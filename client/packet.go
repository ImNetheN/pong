package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

type PacketType int

const (
	PacketReady PacketType = iota
	PacketInt
	PacketConnected
	PacketEnd
	PacketPaddleMoved
	PacketPaddlePos
	PacketBall
)

type Packet struct {
	Type PacketType
	Data []byte
}

func SerializePacket(pack Packet) []byte {
	b, err := json.Marshal(pack)
	if err != nil {
		panic(err)
	}

	return b
}

func DeserealizePacket(packd []byte) Packet {
	packd = []byte(strings.Trim(string(packd), "\x00"))
	var pack Packet
	err := json.Unmarshal(packd, &pack)
	if err != nil {
		fmt.Println(string(packd))
		panic(err)
	}

	return pack
}

func NewReadyPacket() Packet {
	return Packet{Type: PacketReady, Data: []byte{}}
}

func NewConnectedPacket() Packet {
	return Packet{Type: PacketConnected, Data: []byte{}}
}

func NewIntPacket(x int) Packet {
	b, err := json.Marshal(x)
	if err != nil {
		panic(err)
	}

	return Packet{Type: PacketInt, Data: b}
}

func NewEndPacket() Packet {
	return Packet{Type: PacketEnd, Data: []byte{}}
}

type PaddleMovedData struct {
	IsRight  bool
	MoveDist float32
}

func NewPaddleMovedPacket(isRight bool, moveDist float32) Packet {
	b, err := json.Marshal(PaddleMovedData{
		IsRight:  isRight,
		MoveDist: moveDist,
	})
	// TODO: return error instead of panicking
	if err != nil {
		panic(err)
	}
	return Packet{
		Type: PacketPaddleMoved,
		Data: b,
	}
}

// func ReadPaddleMovedPacket(pack Packet) PaddleMovedData {
//     var data PaddleMovedData
//     err := json.Unmarshal(pack.Data, &data)
//     // TODO: ou
//     if err != nil {
//         panic(err)
//     }

//     return data
// }

type PaddlePosData struct {
	IsRight bool
	Pos     float32
}

func NewPaddlePosPacket(isRight bool, pos float32) Packet {
	b, err := json.Marshal(PaddlePosData{
		IsRight: isRight,
		Pos:     pos,
	})
	// TODO: return error instead of panicking
	if err != nil {
		panic(err)
	}
	return Packet{
		Type: PacketPaddlePos,
		Data: b,
	}
}

// func ReadPaddlePosPacket(pack Packet) PaddlePosData {
//     var data PaddlePosData
//     err := json.Unmarshal(pack.Data, &data)
//     // TODO: ou
//     if err != nil {
//         panic(err)
//     }

//     return data
// }

func NewBallPacket(ball Ball) Packet {
	b, err := json.Marshal(ball)
	// TODO YES !!
	if err != nil {
		panic(err)
	}

	return Packet{
		Type: PacketBall,
		Data: b,
	}
}