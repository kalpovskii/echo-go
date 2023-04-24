package main

import (
	"context"
	"echo/internal/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := server.Config{
		Port: 9854,
	}

	srv, err := server.New(cfg)
	if err != nil {
		log.Fatal()
	}

	ctx, cancel := context.WithCancel(context.Background())

	go srv.Run(ctx)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(ch)
	<-ch
	cancel()
}
