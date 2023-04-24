package client

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
)

type Config struct {
	Addr string
}

type Client struct {
	conn net.Conn
}

func New(cfg Config) (*Client, error) {
	conn, err := net.Dial("tcp", cfg.Addr)
	if err != nil {
		return nil, err
	}
	return &Client{
		conn: conn,
	}, nil
}

func (c *Client) Run(ctx context.Context) error {
	defer func() {
		err := c.conn.Close()
		if err != nil {
			log.Println("error closing connect")
		}
	}()
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("text to send:")
			text, err := reader.ReadString('\n')
			if err != nil {
				return err
			}
			if _, err := c.conn.Write([]byte(text + "\r")); err != nil {
				log.Println("error writing message: ", err.Error())
				return err
			}
			message, _ := bufio.NewReader(c.conn).ReadString('\r')
			fmt.Print("Message from server: " + message)
		}
	}
}
