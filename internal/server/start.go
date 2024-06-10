package server

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"paar/internal/process"
	"paar/internal/store"
	"sync"
	"syscall"
)

type Server struct {
	listeningPort string
	ln            net.Listener
	quitCh        chan struct{}
	messageCh     chan []byte
	wg            sync.WaitGroup
	process  			*process.Process
}

func NewServer(listeningPort string) *Server {
	return &Server{
		listeningPort: listeningPort,
		quitCh:        make(chan struct{}),
		messageCh:     make(chan []byte),
		process:       process.New(),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.listeningPort)
	if err != nil {
		return err
	}
	s.ln = ln

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	disk := store.NewDisk(s.process.Store.GetMap())
	m, err := disk.Load("data.json")
	if err != nil {
		return err
	}
	s.process.Store.Initialize(m)

	go s.handleSignals(cancel)
	go s.Accept(ctx)

	<-s.quitCh
	s.ln.Close()
	s.wg.Wait()
	disk.Save("data.json")

	return nil
}

func (s *Server) handleSignals(cancel context.CancelFunc) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	cancel()
	close(s.quitCh)
}

func (s *Server) Accept(ctx context.Context) {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				return 
			default:
				fmt.Println("Error accepting connection:", err)
				continue
			}
		}
		fmt.Println("Accepted connection from", conn.RemoteAddr())
		s.wg.Add(1)
		go s.handle(ctx, conn)
		
	}
}

func (s *Server) handle(ctx context.Context, conn net.Conn) {
	defer s.wg.Done()
	defer conn.Close()
	buf := make([]byte, 1024)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			n, err := conn.Read(buf)
			if err != nil {
				if err == io.EOF {
					fmt.Println("Client closed the connection")
					return
				}
				fmt.Println("Error reading from connection:", err)
				return
			}
			fmt.Println("Received data:", string(buf[:n]))
			s.process.Handle(conn, buf[:n])
		}
	}
}