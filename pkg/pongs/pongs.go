package pongs

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

	SCREEN_HEIGHT = 600
	SCREEN_WIDTH  = 800

	BALL_SIZE = 10
)

type Paddle struct {
	PosY    float32
	IsRight bool
}

func InitLeftPaddle() Paddle {
	return Paddle{
		PosY:    400,
		IsRight: false,
	}
}

func InitRightPaddle() Paddle {
	return Paddle{
		PosY:    400,
		IsRight: true,
	}
}

func (p *Paddle) Draw(dst *ebiten.Image) {
	if !p.IsRight {
		vector.DrawFilledRect(dst, LEFT_PADDLE_X, p.PosY, PADDLE_WIDTH, PADDLE_HEIGHT, color.White, false)
	} else {
		vector.DrawFilledRect(dst, RIGHT_PADDLE_X, p.PosY, PADDLE_WIDTH, PADDLE_HEIGHT, color.White, false)
	}
}

type Ball struct {
	PosX float32
	PosY float32
	VelX float32
	VelY float32
}

func InitBall() Ball {
	return Ball{
		PosX: 400,
		PosY: 400,
		VelX: 1,
		VelY: 1,
	}
}

func (b *Ball) Draw(dst *ebiten.Image) {
	vector.DrawFilledRect(dst, b.PosX, b.PosY, BALL_SIZE, BALL_SIZE, color.White, false)
}

type State struct {
	LeftPaddle  Paddle
	RightPaddle Paddle
	Ball        Ball
}

func BallIntersectsLeftPaddle(ball Ball, leftPadl Paddle) bool {
	if ball.PosX > LEFT_PADDLE_X+PADDLE_WIDTH || ball.PosX+BALL_SIZE < LEFT_PADDLE_X {
		return false
	}
	if ball.PosY > leftPadl.PosY+PADDLE_HEIGHT || ball.PosY+BALL_SIZE < leftPadl.PosY {
		return false
	}

	return true
}

func BallIntersectsRightPaddle(ball Ball, rightPadl Paddle) bool {
	if ball.PosX > RIGHT_PADDLE_X+PADDLE_WIDTH || ball.PosX+BALL_SIZE < RIGHT_PADDLE_X {
		return false
	}
	if ball.PosY > rightPadl.PosY+PADDLE_HEIGHT || ball.PosY+BALL_SIZE < rightPadl.PosY {
		return false
	}

	return true
}
