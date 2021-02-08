package ui

import (
	"github.com/hajimehoshi/ebiten"
	"lockstepuiclient/client/pb"
	"log"
)

type ci struct {
	c  *car
	id uint64
}

type cars struct {
	c []ci
}

func newCars(mid uint64) *cars {
	c := new(cars)
	c.c = append(c.c, ci{
		c:  newCar(mid),
		id: mid,
	})
	return c
}

func (c cars) receiveData(input *pb.InputData) {
	if input == nil {
		return
	}

	id := *input.Id
	for i := 0; i < len(c.c); i++ {
		if c.c[i].id == id {
			c.c[i].c.receiveData(input)
		}
	}
}

func (c cars) joinRoom(rec *pb.S2C_JoinRoomMsg) {
	log.Print("JoinRoom success. rec = ", rec)
	//TODO
}

func (c cars) draw(screen *ebiten.Image) {
	for i := 0; i < len(c.c); i++ {
		c.c[i].c.draw(screen)
	}
}
