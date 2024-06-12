package main

import (
	"encoding/json"
	"net"
	"pong/pkg/packet"
)

func Connect(serverAddressString string) (*net.UDPConn, bool) {
	serverAddress, err := net.ResolveUDPAddr("udp", serverAddressString)
	if err != nil {
		panic(err)
	}

	conn, err := net.DialUDP("udp", nil, serverAddress)
	if err != nil {
		panic(err)
	}
	_, err = conn.Write(packet.SerializePacket(packet.NewConnectedPacket()))
	if err != nil {
		panic(err)
	}

	buf := make([]byte, 64)
	_, err = conn.Read(buf)
	if err != nil {
		panic(err)
	}

	playerNumPacket := packet.DeserealizePacket(buf)
	var playerNum int
	err = json.Unmarshal(playerNumPacket.Data, &playerNum)
	if err != nil {
		panic(err)
	}

	var isRight bool
	if playerNum == 1 {
		isRight = false
	} else if playerNum == 2 {
		isRight = true
	} else {
		panic("aioshentoeiatsoeht connect bad")
	}

	return conn, isRight
}

func WaitUntilReady(conn *net.UDPConn) {
	for {
		buf := make([]byte, 32)
		_, err := conn.Read(buf)
		if err != nil {
			panic(err)
		}

		readyPacket := packet.DeserealizePacket(buf)

		if readyPacket.Type == packet.PacketReady {
			return
		}
	}
}
