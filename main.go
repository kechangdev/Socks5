package main

import (
	"bufio"
	"log"
	"net"
)

const socks5Ver = 0x05
const cmdBind = 0x01
const atypeIPV4 = 0x01
const atypeHOST = 0x03
const atypeIPV6 = 0x04

func main() {
	server, err := net.Listen("tcp", ":1080")
	if err != nil {
		panic(err)
	}
	for {
		client, err := server.Accept()
		if err != nil {
			continue
		}
		defer client.Close()
		go func() {
			reader := bufio.NewReader(client)
			err = Auth(reader, client)
			if err == nil {
				log.Println("Auth Success")
				addr := ""
				port := ""
				addr, port, err = Connect(reader, client)
				if err != nil {
					log.Println("Connect Failed")
					return
				}
				err = Relay(client, reader, addr, port)
				if err != nil {
					log.Println("Relay Failed")
					return
				}
			}
		}()
	}
}
