package main

import (
	"context"
	"echo/internal/client"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := client.Config{
		Addr: "0.0.0.0:9854",
	}
	c, err := client.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		if err := c.Run(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(ch)
	<-ch
	cancel()
}
