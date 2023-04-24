package server

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net"
)

type Config struct {
	Port uint16
}

type Server struct {
	listener net.Listener
}

func New(cfg Config) (*Server, error) {
	fmt.Println("starting server on port ", fmt.Sprintf("0.0.0.0:%d", cfg.Port))
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", cfg.Port))
	if err != nil {
		return nil, err
	}
	return &Server{
		listener: listener,
	}, nil
}

func (s *Server) Run(ctx context.Context) {
	defer func() {
		err := s.listener.Close()
		if err != nil {
			log.Println("error closing connect")
		}
	}()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			conn, err := s.listener.Accept()
			if err != nil {
				log.Println("failed to accept connection, err:", err)
				continue
			}
			go s.handleConnection(ctx, conn)
		}
	}
}

func (s *Server) handleConnection(ctx context.Context, conn net.Conn) {
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Println("error closing connect")
		}
	}()
	reader := bufio.NewReader(conn)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			bytes, err := reader.ReadBytes(byte('\r'))
			if err != nil {
				if err != io.EOF {
					log.Println("failed to read data, err:", err)
				}
				return
			}
			log.Printf("request: %s", bytes)

			line := fmt.Sprintf("%s", bytes)
			fmt.Printf("response: %s", line)
			if _, err := conn.Write([]byte(line)); err != nil {
				log.Println(err)
				return
			}
		}
	}
}
