package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

func SerializeFloat(x float32) []byte {
	return []byte(fmt.Sprintf("%v", x))
}

func DeserializeFloat(b []byte) float32 {
	s := string(b)
	s = strings.Trim(s, string(0))
	x, err := strconv.ParseFloat(s, 32)
	if err != nil {
		panic(err)
	}

	return float32(x)
}

func Connect(serverAddressString string) (*net.UDPConn, bool) {
	serverAddress, err := net.ResolveUDPAddr("udp", serverAddressString)
	if err != nil {
		panic(err)
	}

	conn, err := net.DialUDP("udp", nil, serverAddress)
	if err != nil {
		panic(err)
	}
	_, err = conn.Write([]byte{})
	if err != nil {
		panic(err)
	}

	buf := make([]byte, 8)
	n, err := conn.Read(buf)
	if err != nil {
		panic(err)
	}
	if n != 1 {
		fmt.Println("n is not 1, buf is", buf)
		panic("mrrow")
	}

	var isRight bool
	if buf[0] == '1' {
		isRight = false
	} else if buf[0] == '2' {
		isRight = true
	} else {
		panic("aioshentoeiatsoeht connect bad")
	}

	return conn, isRight
}

func WaitUntilReady(conn *net.UDPConn) {
	for {
		buf := make([]byte, 5)
		_, err := conn.Read(buf)
		if err != nil {
			panic(err)
		}

		if string(buf) == "ready" {
			return
		}
	}
}
