package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net"
)

func Relay(conn net.Conn, src *bufio.Reader, addr string, port string) error {
	dest, err := net.Dial("tcp", fmt.Sprintf("%s:%s", addr, port))
	log.Println("dial", addr, port)
	if err != nil {
		return err
	}
	defer dest.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		_, _ = io.Copy(dest, src)
	}()
	go func() {
		_, _ = io.Copy(conn, dest)
	}()
	<-ctx.Done()
	return nil
}
