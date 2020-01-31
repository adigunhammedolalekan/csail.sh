package services

import (
	"bufio"
	"errors"
	"log"
	"net"
	"strings"
	"sync"
)

type TcpServer struct {
	conns map[string]net.Conn
	mtx sync.Mutex
}

func NewTcpServer() *TcpServer {
	return &TcpServer{conns: make(map[string]net.Conn)}
}

func (s *TcpServer) Run(addr string) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			// could be any type of error: Log error and continue to listen
			log.Println("failed to accept a new conn: ", err)
			continue
		}
		go s.handle(conn)
	}
}

func (s *TcpServer) handle(conn net.Conn) {
	for {
		buf := bufio.NewReader(conn)
		m, err := buf.ReadString(byte('\n'))
		if err != nil {
			log.Println("[TCP]: failed to read data from TCP client: ", err)
			continue
		}
		key := strings.TrimSuffix(m, "\n")
		if key == "" || len(key) != 10 {
			log.Println("invalid key: ", key)
			continue
		}
		s.register(conn, key)
		break
	}
}

func (s *TcpServer) register(conn net.Conn, k string) {
	s.mtx.Lock()
	s.mtx.Unlock()
	s.conns[k] = conn
}

func (s *TcpServer) Send(key string, data []byte) error {
	s.mtx.Lock()
	conn, ok := s.conns[key]
	s.mtx.Unlock()
	if !ok {
		log.Println("client not found: Key=", key)
		return errors.New("client not found")
	}
	if _, err := conn.Write(data); err != nil {
		return err
	}
	return nil
}
