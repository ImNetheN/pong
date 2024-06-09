package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	PADDLE_WIDTH   = 10
	PADDLE_HEIGHT  = 100
	LEFT_PADDLE_X  = 10
	RIGHT_PADDLE_X = 780
	PADDLE_SPEED   = 5

	BALL_SIZE = 10
)

type Paddle struct {
	posY    float32
	isRight bool
}

func InitLeftPaddle() Paddle {
	return Paddle{
		posY:    400,
		isRight: false,
	}
}

func InitRightPaddle() Paddle {
	return Paddle{
		posY:    400,
		isRight: true,
	}
}

func (p *Paddle) Draw(dst *ebiten.Image) {
	if !p.isRight {
		vector.DrawFilledRect(dst, LEFT_PADDLE_X, p.posY, PADDLE_WIDTH, PADDLE_HEIGHT, color.White, false)
	} else {
		vector.DrawFilledRect(dst, RIGHT_PADDLE_X, p.posY, PADDLE_WIDTH, PADDLE_HEIGHT, color.White, false)
	}
}

type Ball struct {
	posX float32
	posY float32
	velX float32
	velY float32
}

func InitBall() Ball {
	return Ball{
		posX: 400,
		posY: 400,
		velX: 1,
		velY: 1,
	}
}

func (b *Ball) Draw(dst *ebiten.Image) {
	vector.DrawFilledRect(dst, b.posX, b.posY, BALL_SIZE, BALL_SIZE, color.White, false)
}
