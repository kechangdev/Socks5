package main

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
)

func Connect(reader *bufio.Reader, conn net.Conn) (addr string, port string, err error) {
	// +----+-----+-------+------+----------+----------+
	// |VER | CMD |  RSV  | ATYP | DST.ADDR | DST.PORT |
	// +----+-----+-------+------+----------+----------+
	// | 1  |  1  | X'00' |  1   | Variable |    2     |
	// +----+-----+-------+------+----------+----------+
	// VER 版本号，socks5的值为0x05
	// CMD 0x01表示CONNECT请求
	// RSV 保留字段，值为0x00
	// ATYP 目标地址类型，DST.ADDR的数据对应这个字段的类型。
	//   0x01表示IPv4地址，DST.ADDR为4个字节
	//   0x03表示域名，DST.ADDR是一个可变长度的域名
	// DST.ADDR 一个可变长度的值
	// DST.PORT 目标端口，固定2个字节
	buf := make([]byte, 4)
	_, err = io.ReadFull(reader, buf)
	if err != nil {
		return addr, port, errors.New("read header failed")
	}
	ver, cmd, atyp := buf[0], buf[1], buf[3]
	if ver != socks5Ver {
		return addr, port, errors.New("invalid Socks version")
	}
	if cmd != 0x01 {
		return addr, port, errors.New("invalid command")
	}
	switch atyp {
	case 0x01:
		_, err = io.ReadFull(reader, buf)
		if err != nil {
			return addr, port, errors.New("read atyp failed")
		}
		addr = fmt.Sprintf("%d.%d.%d.%d", buf[0], buf[1], buf[2], buf[3])
	case 0x03:
		return addr, port, errors.New("IPv6: no supported yet")
	default:
		return addr, port, errors.New("invalid ATYP")
	}
	_, err = io.ReadFull(reader, buf[:2])
	if err != nil {
		return addr, port, errors.New("read port failed")
	}
	port = strconv.Itoa(int(binary.BigEndian.Uint16(buf[:2])))
	log.Println("addr", addr, "port", port)
	// 按照协议，收到请求后需要回复
	// +----+-----+-------+------+----------+----------+
	// |VER | REP |  RSV  | ATYP | BND.ADDR | BND.PORT |
	// +----+-----+-------+------+----------+----------+
	// | 1  |  1  | X'00' |  1   | Variable |    2     |
	// +----+-----+-------+------+----------+----------+
	// VER socks版本，这里为0x05
	// REP Relay field,内容取值如下 X’00’ succeeded
	// RSV 保留字段
	// ATYPE 地址类型
	// BND.ADDR 服务绑定的地址
	// BND.PORT 服务绑定的端口DST.PORT
	_, err = conn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0, 0, 0, 0, 0, 0})
	if err != nil {
		return addr, port, fmt.Errorf("write failed: %w", err)
	}
	return addr, port, nil
}
