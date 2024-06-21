package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"pong/pkg/packet"
	"pong/pkg/pongs"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type State struct {
	LeftPaddle  pongs.Paddle
	RightPaddle pongs.Paddle
	Ball        pongs.Ball
}

type Game struct {
	State State

	IsRight bool
	// OtherPaddlePos chan float32
	StateChannel chan State
	ServerConn   net.Conn
}

func InitGame(isRight bool, conn net.Conn) Game {
	return Game{
		State: State{
			LeftPaddle:  pongs.InitLeftPaddle(),
			RightPaddle: pongs.InitRightPaddle(),
			Ball:        pongs.InitBall(),
		},
		IsRight: isRight,
		// OtherPaddlePos: make(chan float32),
		// TODO: check if it needs to be buffered
		StateChannel: make(chan State, 1),
		ServerConn:   conn,
	}
}

func (g *Game) ControlledPaddle() *pongs.Paddle {
	if g.IsRight {
		return &g.State.RightPaddle
	} else {
		return &g.State.LeftPaddle
	}
}

func (g *Game) OtherPaddle() *pongs.Paddle {
	if !g.IsRight {
		return &g.State.RightPaddle
	} else {
		return &g.State.LeftPaddle
	}
}

func (g *Game) Update() error {
	for end := false; !end; {
		select {
		case newState := <-g.StateChannel:
			// g.State = newState
			g.State.Ball = newState.Ball
			if g.IsRight {
				g.State.LeftPaddle = newState.LeftPaddle
			} else {
				g.State.RightPaddle = newState.RightPaddle
			}
		default:
			end = true
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		g.ControlledPaddle().PosY += pongs.PADDLE_SPEED
	}

	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		g.ControlledPaddle().PosY -= pongs.PADDLE_SPEED
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButton0) {
		_, mouseY := ebiten.CursorPosition()
		g.ControlledPaddle().PosY = float32(mouseY)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.State.LeftPaddle.Draw(screen)
	g.State.RightPaddle.Draw(screen)
	g.State.Ball.Draw(screen)
}

func (g *Game) Layout(w, h int) (int, int) {
	return 800, 600
}

func (g *Game) SendPaddleInfo() {
	for {
		paddlePosPacket := packet.NewPaddlePosPacket(g.IsRight, g.ControlledPaddle().PosY)
		_, err := g.ServerConn.Write(packet.SerializePacket(paddlePosPacket))

		if err != nil {
			panic(err)
		}
		time.Sleep(time.Millisecond * 17)
	}
}

func (g *Game) ReceiveStateInfo() {
	for {
		buf := make([]byte, 256)

		_, err := g.ServerConn.Read(buf)
		if err != nil {
			panic(err)
		}

		statePacket := packet.DeserealizePacket(buf)
		var newState State
		err = json.Unmarshal(statePacket.Data, &newState)
		if err != nil {
			panic(err)
		}

		g.StateChannel <- newState
	}
}

func main() {
	var serverAddressString string
	flag.StringVar(&serverAddressString, "a", ":4321", "server address")
	flag.Parse()

	serverConnection, isRight := Connect(serverAddressString)
	WaitUntilReady(serverConnection)

	ebiten.SetWindowSize(800, 600)

	fmt.Println(isRight)
	game := InitGame(isRight, serverConnection)

	go game.SendPaddleInfo()
	go game.ReceiveStateInfo()
	if err := ebiten.RunGame(&game); err != nil {
		panic(err)
	}
}
