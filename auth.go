package main

import (
	"bufio"
	"errors"
	"io"
	"log"
	"net"
)

func Auth(reader *bufio.Reader, client net.Conn) (err error) {
	// +----+----------+----------+
	// |VER | NMETHODS | METHODS  |
	// +----+----------+----------+
	// | 1  |    1     | 1 to 255 |
	// +----+----------+----------+
	// VER: 协议版本，socks5为0x05
	// NMETHODS: 支持认证的方法数量
	// METHODS: 对应NMETHODS，NMETHODS的值为多少，METHODS就有多少个字节。RFC预定义了一些值的含义，内容如下:
	// X’00’ NO AUTHENTICATION REQUIRED
	// X’02’ USERNAME/PASSWORD
	ver, err := reader.ReadByte()
	if err != nil {
		return errors.New("read var failed")
	}
	if ver != socks5Ver {
		return errors.New("socks5 ver not supported")
	}
	methodSize, err := reader.ReadByte()
	if err != nil {
		return errors.New("read methodSize failed")
	}
	method := make([]byte, methodSize)
	_, err = io.ReadFull(reader, method)
	if err != nil {
		return errors.New("read method failed")
	}
	log.Println("ver", ver, "method", method)
	// +----+--------+
	// |VER | METHOD |
	// +----+--------+
	// | 1  |   1    |
	// +----+--------+
	_, err = client.Write([]byte{socks5Ver, 0x00})
	if err != nil {
		return errors.New("write socks5 ver failed")
	}
	return nil
}
