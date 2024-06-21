package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"pong/pkg/packet"
	"pong/pkg/pongs"
	"time"
)

type Manager struct {
	Conn   *net.UDPConn
	P1Addr *net.UDPAddr
	P2Addr *net.UDPAddr

	State      pongs.State
	LastUpdate time.Time
}

func (m *Manager) SendReady() {
	_, err := m.Conn.WriteToUDP(packet.SerializePacket(packet.NewReadyPacket()), m.P1Addr)
	if err != nil {
		panic(err)
	}

	_, err = m.Conn.WriteToUDP(packet.SerializePacket(packet.NewReadyPacket()), m.P2Addr)
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
			ball := &m.State.Ball
			leftPadl := m.State.LeftPaddle
			rightPadl := m.State.RightPaddle

			ball.PosX += ball.VelX
			ball.PosY += ball.VelY

			if ball.PosY+pongs.BALL_SIZE < 0 || ball.PosY > pongs.SCREEN_HEIGHT {
				ball.VelY *= -1
			}

			if pongs.BallIntersectsLeftPaddle(*ball, leftPadl) || pongs.BallIntersectsRightPaddle(*ball, rightPadl) {
				ball.VelX *= -1
			}

			if ball.PosX < 0 || ball.PosX+pongs.BALL_SIZE > pongs.SCREEN_WIDTH {
				ball.PosX = 400
				ball.PosY = 400
			}
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

		paddlePacket := packet.DeserealizePacket(buf)

		var paddle pongs.Paddle
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
		statePacketSerialized := packet.SerializePacket(packet.NewStatePacket(m.State))

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
	_, err = conn.WriteToUDP(packet.SerializePacket(packet.NewIntPacket(1)), p1Addr)
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

	_, err = conn.WriteToUDP(packet.SerializePacket(packet.NewIntPacket(2)), p2Addr)
	if err != nil {
		panic(err)
	}

	fmt.Println("Player 2 connected", p2Addr)

	manager := Manager{
		Conn:   conn,
		P1Addr: p1Addr,
		P2Addr: p2Addr,

		State: pongs.State{
			Ball: pongs.Ball{
				PosX: 400,
				PosY: 400,
				VelX: 3,
				VelY: 3,
			},
		},
	}

	manager.SendReady()

	manager.LastUpdate = time.Now()

	go manager.Receive()
	go manager.UpdateBall()
	manager.Send()
}
