package ui

import (
	"github.com/hajimehoshi/ebiten"
	"lockstepuiclient/client/pb"
)

type PupuCar struct {
	c *cars
	i *input
	s *status
	r *room
}

func NewPupuCar(ss *status) *PupuCar {
	g := new(PupuCar)
	g.c = newCars()
	g.i = newInput()
	g.r = newRoom()
	g.s = ss
	return g
}

func (g *PupuCar) Update(_ *ebiten.Image) error {
	switch g.s.level {
	case 0:
		return g.r.update()
	case 1:
		return g.i.update()
	}
	return nil
}

func (g *PupuCar) Draw(screen *ebiten.Image) {
	switch g.s.level {
	case 0:
		g.r.draw(screen)
	case 1:
		g.c.draw(screen)
	}
}

func (g *PupuCar) Layout(_, _ int) (screenWidth, screenHeight int) {
	return 320, 240
}

func (g *PupuCar) ReceiveData(input []*pb.InputData) {
	if len(input) == 0 {
		return
	}

	for _, i := range input {
		g.c.receiveData(i)
	}
}

func (g *PupuCar) JoinRoom(rec *pb.S2C_JoinRoomMsg) {
	g.c.joinRoom(rec)
}

func (g *PupuCar) GameStart() {
	g.s.level = 1
}
