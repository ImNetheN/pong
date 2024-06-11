package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	LeftPaddle  Paddle
	RightPaddle Paddle
	Ball        Ball

	IsRight        bool
	OtherPaddlePos chan float32
	ServerConn     net.Conn
}

func InitGame(isRight bool, conn net.Conn) Game {
	return Game{
		LeftPaddle:     InitLeftPaddle(),
		RightPaddle:    InitRightPaddle(),
		Ball:           InitBall(),
		IsRight:        isRight,
		OtherPaddlePos: make(chan float32),
		ServerConn:     conn,
	}
}

func (g *Game) Update() error {
	var controlledPaddle *Paddle
	var otherPaddle *Paddle
	if g.IsRight {
		controlledPaddle = &g.RightPaddle
		otherPaddle = &g.LeftPaddle
	} else {
		controlledPaddle = &g.LeftPaddle
		otherPaddle = &g.RightPaddle
	}

	for end := false; !end; {
		select {
		case posY := <-g.OtherPaddlePos:
			otherPaddle.posY = posY
		default:
			end = true
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		controlledPaddle.posY += PADDLE_SPEED
	}

	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		controlledPaddle.posY -= PADDLE_SPEED
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.LeftPaddle.Draw(screen)
	g.RightPaddle.Draw(screen)
	g.Ball.Draw(screen)
}

func (g *Game) Layout(w, h int) (int, int) {
	return 800, 600
}

func (g *Game) SendPaddleInfo() {
	var paddle *Paddle
	if g.IsRight {
		paddle = &g.RightPaddle
	} else {
		paddle = &g.LeftPaddle
	}

	for {
		paddlePosPacket := NewPaddlePosPacket(g.IsRight, paddle.posY)
		_, err := g.ServerConn.Write(SerializePacket(paddlePosPacket))

		if err != nil {
			panic(err)
		}
		time.Sleep(time.Millisecond * 66)
	}
}

func (g *Game) ReceivePaddleInfo() {
	for {
		buf := make([]byte, 256)

		_, err := g.ServerConn.Read(buf)
		if err != nil {
			panic(err)
		}

		paddlePosPacket := DeserealizePacket(buf)
		var paddleData PaddlePosData
		err = json.Unmarshal(paddlePosPacket.Data, &paddleData)
		if err != nil {
			panic(err)
		}

		g.OtherPaddlePos <- paddleData.Pos
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
	go game.ReceivePaddleInfo()
	if err := ebiten.RunGame(&game); err != nil {
		panic(err)
	}
}
