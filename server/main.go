package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"time"
)

const (
	PADDLE_WIDTH   = 10
	PADDLE_HEIGHT  = 100
	LEFT_PADDLE_X  = 10
	RIGHT_PADDLE_X = 780
	PADDLE_SPEED   = 5

	BALL_SIZE = 10
)

type Ball struct {
	PosX float32
	PosY float32
	VelX float32
	VelY float32
}

type Paddle struct {
	IsRight bool
	PosY    float32
}

type State struct {
	LeftPaddle  Paddle
	RightPaddle Paddle
	Ball        Ball
}

type Manager struct {
	Conn   *net.UDPConn
	P1Addr *net.UDPAddr
	P2Addr *net.UDPAddr

	State      State
	LastUpdate time.Time
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

// TODO: mutex
func (m *Manager) UpdateBall() {
	updateTicker := time.NewTicker(time.Millisecond * 17)

	for {
		select {
		case <-updateTicker.C:
			m.State.Ball.PosX += m.State.Ball.VelX
			m.State.Ball.PosY += m.State.Ball.VelY
		}
	}
}

func (m *Manager) Receive() {
	for {
		buf := make([]byte, 256)
		_, addr, err := m.Conn.ReadFromUDP(buf)
		if err != nil {
			panic(err)
		}

		paddlePacket := DeserealizePacket(buf)

		var paddle Paddle
		err = json.Unmarshal(paddlePacket.Data, &paddle)
		if err != nil {
			panic(err)
		}

		// TODO: put player data in the packet instead of cheaking ip
		// TODO: mutex
		if addr.IP.Equal(m.P1Addr.IP) && addr.Port == m.P1Addr.Port {
			m.State.LeftPaddle = paddle
		} else if addr.IP.Equal(m.P2Addr.IP) && addr.Port == m.P2Addr.Port {
			m.State.RightPaddle = paddle
		} else {
			panic(fmt.Sprint("bad client address", addr))
		}

	}
}

func (m *Manager) Send() {
	for {
		statePacketSerialized := SerializePacket(NewStatePacket(m.State))

		_, err := m.Conn.WriteToUDP(statePacketSerialized, m.P1Addr)
		if err != nil {
			panic(err)
		}

		_, err = m.Conn.WriteToUDP(statePacketSerialized, m.P2Addr)
		if err != nil {
			panic(err)
		}

		time.Sleep(time.Millisecond * 33)
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

		State: State{
			Ball: Ball{
				PosX: 400,
				PosY: 400,
				VelX: 1,
				VelY: 1,
			},
		},
	}

	manager.SendReady()

	manager.LastUpdate = time.Now()

	go manager.Receive()
	go manager.UpdateBall()
	manager.Send()
}
