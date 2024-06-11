package main

import (
	"flag"
	"fmt"
	"net"
)

type Ball struct {
	posX float32
	posY float32
	velX float32
	velY float32
}

type Manager struct {
	Conn   *net.UDPConn
	P1Addr *net.UDPAddr
	P2Addr *net.UDPAddr

	P1Chan chan []byte
	P2Chan chan []byte
}

func (m *Manager) SendReady() {
	_, err := m.Conn.WriteToUDP(SerializePacket(NewReadyPacket()), m.P1Addr)
	if err != nil {
		panic(err)
	}

	_, err = m.Conn.WriteToUDP(SerializePacket(NewReadyPacket()), m.P2Addr)
	if err != nil {
		panic(err)
	}
}

func (m *Manager) Receive() {
	for {
		buf := make([]byte, 256)
		_, addr, err := m.Conn.ReadFromUDP(buf)
		if err != nil {
			panic(err)
		}

		if addr.IP.Equal(m.P1Addr.IP) && addr.Port == m.P1Addr.Port {
			// fmt.Println("received from 1:", string(buf))
			m.P2Chan <- buf
		} else if addr.IP.Equal(m.P2Addr.IP) && addr.Port == m.P2Addr.Port {
			// fmt.Println("received from 2:", string(buf))
			m.P1Chan <- buf
		} else {
			panic(fmt.Sprint("bad client address", addr))
		}
	}
}

func (m *Manager) Send() {
	for {
		select {
		case p := <-m.P1Chan:
			_, err := m.Conn.WriteToUDP(p, m.P1Addr)
			if err != nil {
				panic(err)
			}
			// fmt.Println("sent to 1:", string(p))
		case p := <-m.P2Chan:
			_, err := m.Conn.WriteToUDP(p, m.P2Addr)
			if err != nil {
				panic(err)
			}
			// fmt.Println("sent to 2:", string(p))
		}
	}
}

func main() {
	var serverAddressString string
	flag.StringVar(&serverAddressString, "a", ":4321", "server address")
	flag.Parse()

	serverAddress, err := net.ResolveUDPAddr("udp", serverAddressString)
	if err != nil {
		panic(err)
	}

	conn, err := net.ListenUDP("udp", serverAddress)
	if err != nil {
		panic(err)
	}

	fmt.Println("listening")

	_, p1Addr, err := conn.ReadFromUDP([]byte{})
	if err != nil {
		panic(err)
	}
	// p := (NewIntPacket(1))
	// fmt.Println(string(p.Data))
	// pp := SerializePacket(p)
	// fmt.Println((pp))
	_, err = conn.WriteToUDP(SerializePacket(NewIntPacket(1)), p1Addr)
	if err != nil {
		panic(err)
	}

	fmt.Println("Player 1 connected", p1Addr)

	var p2Addr *net.UDPAddr
	for {
		_, p2Addr, err = conn.ReadFromUDP([]byte{})
		if err != nil {
			panic(err)
		}

		if p2Addr != p1Addr {
			break
		}
	}

	_, err = conn.WriteToUDP(SerializePacket(NewIntPacket(2)), p2Addr)
	if err != nil {
		panic(err)
	}

	fmt.Println("Player 2 connected", p2Addr)

	manager := Manager{
		Conn:   conn,
		P1Addr: p1Addr,
		P2Addr: p2Addr,

		P1Chan: make(chan []byte),
		P2Chan: make(chan []byte),
	}

	manager.SendReady()

	go manager.Receive()
	manager.Send()
}
