package process

import (
	"fmt"
	"paar/internal/store"
	"paar/internal/utils"
	"strings"
	"time"
)

type Process struct {
	Store *store.Storage
}

func New() *Process{
	return &Process{
		Store: store.NewStorage(),
	}
}



type Conn interface {
	Write([]byte) (int, error)
	Close() error
}

func (p *Process) Handle(conn Conn, buf []byte) {
	command := strings.Fields(string(buf))

	switch command[0] {
		case "GET":
			p.handleGet(conn, command)
			return
		case "SET":
			p.handleSet(conn, command)
			return
		case "DEL":
			p.handleDel(conn, command)
			return
		case "KEYS":
			p.handleKeys(conn, command)
			return
		case "QUIT":
			p.handleQuit(conn, command)
			return
		case "PING":
			p.handlePing(conn, command)
			return
		case "EXPIRE":
			p.handleExpire(conn, command)
			return
		default:
			conn.Write([]byte(fmt.Sprintf("UNKNOWN COMMAND %s\n", command[0])))
			return
	}
}

func (p *Process) handlePing(conn Conn, command []string) {
	conn.Write([]byte("PONG\n"))
}

func (p *Process) handleGet(conn Conn, command []string) {
	if len(command) < 2 {
		conn.Write([]byte("KEY MISSING\n"))
		return
	}
	key := command[1]
	value, ok := p.Store.Load(key)
	if !ok {
		conn.Write([]byte("NOT FOUND\n"))
		return
	}
	if value.ExpireTo.Before(time.Now()) {
		p.Store.Delete(key)
		conn.Write([]byte("NOT FOUND\n"))
		return
	}
	conn.Write([]byte(value.Value + "\n"))
}

func (p *Process) handleSet(conn Conn, command []string) {
	if len(command) < 2 {
		conn.Write([]byte("KEY MISSING\n"))
		return
	}
	key := command[1]
	if len(command) < 3 {
		conn.Write([]byte("VALUE MISSING\n"))
		return
	}
	value := command[2]
	p.Store.Store(key, store.Values{
		Value: value,
		ExpireTo: time.Now().Add(time.Duration(1<<63 - 1)),
	})
	conn.Write([]byte("STORED\n"))
}

func (p *Process) handleDel(conn Conn, command []string) {
	if len(command) < 2 {
		conn.Write([]byte("KEY MISSING\n"))
		return
	}
	key := command[1]
	p.Store.Delete(key)
	conn.Write([]byte("DELETED\n"))
}

func (p *Process) handleQuit(conn Conn, command []string) {
	conn.Write([]byte("OK\n"))
	conn.Close()
}

func (p *Process) handleKeys(conn Conn, command []string) {
	key := ""
	if len(command) > 1 {
		key = command[1]
	}
	keys := []string{}
	p.Store.Range(func(k string, v store.Values) bool {
		if strings.HasPrefix(k, key) {
			if v.ExpireTo.Before(time.Now()) {
				p.Store.Delete(k)
			} else {
				keys = append(keys, k)
			}
		}
		return true
	})
	conn.Write([]byte(strings.Join(keys, "\n")))
	conn.Write([]byte("\n"))
}

func (p *Process) handleExpire(conn Conn, command []string) {
	if len(command) < 3 {
		conn.Write([]byte("EXPIRE DATA MISSING\n"))
		return
	}
	key := command[1]
	value , ok := p.Store.Load(key)
	if !ok {
		conn.Write([]byte("KEY NOT FOUND\n"))
		return
	}
	expireToSrt := command[2]
	now := time.Now()
	duration, err := utils.ParseDuration(expireToSrt)
	if err!= nil {
		conn.Write([]byte("EXPIRE DATA INVALID\n"))
		return
	}
	expireTo := now.Add(duration)
	value.ExpireTo = expireTo
	p.Store.Store(key, value)
	conn.Write([]byte(fmt.Sprintf("EXPIRE TO %v\n", expireTo.Format("2006-01-02 15:04:05"))))
}







