package jakis

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

type Server struct {
	addr string
	port int

	dict Dict

	closeCh chan struct{}
}

func (s *Server) handle(ctx context.Context, conn net.Conn) {
	b := bufio.NewReader(conn)
	for {
		select {
		case <-ctx.Done():
			conn.Write([]byte("BYE\n"))
			return
		default:
			line, _ := b.ReadBytes('\n')
			cmd := strings.TrimSuffix(string(line), "\n")
			if len(cmd) > 0 {
				s.processCmd(ctx, conn, cmd)
			}
		}
	}
}

func (s *Server) processCmd(ctx context.Context, conn net.Conn, cmd string) {
	log.Printf("[%s] %s\n", conn.RemoteAddr().String(), cmd)
	block := strings.Split(strings.ToLower(cmd), " ")
	switch block[0] {
	case "exit":
		s.closeCh <- struct{}{}
	case "get":
		resp := fmt.Sprintf("%s\n", s.dict.Get(block[1]))
		conn.Write([]byte(resp))
	case "set":
		s.dict.Set(block[1], block[2])
		conn.Write([]byte("OK\n"))
	default:
		resp := fmt.Sprintf("(error) ERR unknown command `%s`, with args beginning with:\n", cmd)
		conn.Write([]byte(resp))
	}
}

func (s *Server) Run() {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.addr, s.port))
	defer l.Close()
	if err != nil {
		log.Printf("Fail to listen %s:%d\n", s.addr, s.port)
		os.Exit(1)
	}
	log.Printf("Listening %s:%d\n", s.addr, s.port)

	for {
		conn, err := l.Accept()
		defer conn.Close()
		if err != nil {
			log.Printf("Fail to accept connection, %v", err)
			return
		}

		ctx, cancel := context.WithCancel(context.Background())
		go s.handle(ctx, conn)

		go func() {
			select {
			case <-s.closeCh:
				conn.Close()
				cancel()
			}
		}()
	}

}

func NewServer(addr string, port int) *Server {

	return &Server{
		addr:    addr,
		port:    port,
		dict:    NewSimpleMap(),
		closeCh: make(chan struct{}),
	}
}
