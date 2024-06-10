package main

import (
	s "paar/internal/server"
)

func main() {
	server := s.NewServer(":8081")
	err := server.Start()
	if err != nil {
		panic(err)
	}
}
