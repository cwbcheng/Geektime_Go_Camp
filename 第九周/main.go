package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"time"
)

func echo(conn *net.TCPConn) {
	tick := time.Tick(5 * time.Second) // 五秒的心跳间隔
	for now := range tick {
		n, err := conn.Write([]byte(now.String()))
		if err != nil {
			log.Println(err)
			conn.Close()
			return
		}
		fmt.Printf("send %d bytes to %s\n", n, conn.RemoteAddr())
	}
}

func main() {
	address := net.TCPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 7300,
	}

	listener, err := net.ListenTCP("tcp4", &address)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Fatal(err)
		}
		var model = decodeGoimModel(conn)
		handle(model)
	}
}

func handle(model goimModel) {
	//
}

func decodeGoimModel(conn *net.TCPConn) goimModel {
	buffer1 := [16]byte{}
	_, err := conn.Read(buffer1[0:16])
	if err != nil {
		log.Fatal(err)
	}
	var model = new(goimModel)
	model.packageLength = binary.BigEndian.Uint32(buffer1[0:4])
	model.headerLength = binary.BigEndian.Uint16(buffer1[4:6])
	model.protocalVersion = binary.BigEndian.Uint16(buffer1[6:8])
	model.operation = binary.BigEndian.Uint32(buffer1[8:12])
	model.sequenceId = binary.BigEndian.Uint32(buffer1[12:16])
	bodyLength := model.packageLength - uint32(model.headerLength)
	buffer := make([]byte, bodyLength)
	count, err1 := conn.Read(buffer[0:bodyLength])
	if err1 != nil {
		log.Fatal(err1)
	}
	model.body = buffer[0:count]
	return *model
}

type goimModel struct {
	packageLength   uint32
	headerLength    uint16
	protocalVersion uint16
	operation       uint32
	sequenceId      uint32
	body            []byte
}
