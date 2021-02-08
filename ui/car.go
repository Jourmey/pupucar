package ui

import (
	"github.com/hajimehoshi/ebiten"
	"lockstepuiclient/client/pb"
)

type car struct {
	img    *ebiten.Image
	x      float64
	y      float64
	width  int
	height int
}

func (c *car) draw(screen *ebiten.Image) {
	geom := ebiten.GeoM{}
	geom.Translate(c.x, c.y)
	_ = screen.DrawImage(c.img, &ebiten.DrawImageOptions{
		GeoM: geom,
	})
}

func newCar(id uint64) *car {
	b := new(car)
	b.width = 3
	b.height = 3
	b.img, _ = ebiten.NewImage(b.width, b.height, ebiten.FilterDefault)
	b.img.Fill(RainbowPal[id])
	return b
}

func (c *car) receiveData(i *pb.InputData) {
	if i == nil {
		return
	}
	sid := GameSDI(*i.Sid)

	if (sid & Up) != 0 {
		c.y -= 1
	} else if (sid & Down) != 0 {
		c.y += 1
	} else if (sid & Left) != 0 {
		c.x -= 1
	} else if (sid & Right) != 0 {
		c.x += 1
	}
}
