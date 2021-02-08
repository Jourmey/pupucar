package ui

import (
	"github.com/hajimehoshi/ebiten"
	"lockstepuiclient/client/pb"
)

type PupuCar struct {
	//m *car
	o *cars
	i *input
	s *status
}

func NewPupuCar(ss *status) *PupuCar {
	g := new(PupuCar)
	//g.m = newCar(1)
	g.o = newCars(ss.id)
	g.i = newInput()
	g.s = ss
	return g
}

func (g *PupuCar) Update(_ *ebiten.Image) error {
	_ = g.i.Update()
	return nil
}

func (g *PupuCar) Draw(screen *ebiten.Image) {
	g.o.draw(screen)
}

func (g *PupuCar) Layout(_, _ int) (screenWidth, screenHeight int) {
	return 320, 240
}

func (g *PupuCar) ReceiveData(input []*pb.InputData) {
	if len(input) == 0 {
		return
	}

	for _, i := range input {
		g.o.receiveData(i)
	}
}

func (g *PupuCar) JoinRoom(rec *pb.S2C_JoinRoomMsg) {
	g.o.joinRoom(rec)
}
