package TcpEchoServer

import (
	"bufio"
	"log"
	"net"
)

func Echo(client net.Conn, reader *bufio.Reader) {
	defer client.Close()
	reder := bufio.NewReader(client)
	for {
		line, err := reder.ReadByte()
		if err != nil {
			log.Println(err)
			break
		}
		_, err = client.Write([]byte{line})
		if err != nil {
			log.Println(err)
			break
		}
	}
}
