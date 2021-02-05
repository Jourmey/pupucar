package main

import (
	"image/color"
	"lockstepuiclient/client"
	"lockstepuiclient/pb"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten"
)

var rainbowPal = []color.RGBA{
	{0xff, 0x00, 0x00, 0xff},
	{0xff, 0x7f, 0x00, 0xff},
	{0xff, 0xff, 0x00, 0xff},
	{0x00, 0xff, 0x00, 0xff},
	{0x00, 0x00, 0xff, 0xff},
	{0x4b, 0x00, 0x82, 0xff},
	{0x8f, 0x00, 0xff, 0xff},
}

var sameColorCounter = 0
var rainbowColorIndex = 0

var game *Game
var frameID uint32

func getRainbowColor() color.Color {
	if sameColorCounter == 0 {
		rainbowColorIndex = rand.Intn(len(rainbowPal))
	}
	sameColorCounter = (sameColorCounter + 1) % 10
	return rainbowPal[rainbowColorIndex]
}

type Game struct {
	m *Block
	o *Block
}

type Block struct {
	img    *ebiten.Image
	x      float64
	y      float64
	width  int
	height int
}

func NewBlock() *Block {
	b := new(Block)
	b.width = 3
	b.height = 3
	b.img, _ = ebiten.NewImage(b.width, b.height, ebiten.FilterDefault)
	b.img.Fill(getRainbowColor())
	return b
}

func (g *Game) Init() error {
	g.m = NewBlock()
	g.o = NewBlock()
	//g.direction = Right
	return nil
}

func (g *Game) Update(_ *ebiten.Image) error {
	frameID++
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		client.SendAction(frameID, 1)
	} else if ebiten.IsKeyPressed(ebiten.KeyDown) {
		client.SendAction(frameID, 2)
	}

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		client.SendAction(frameID, 3)
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		client.SendAction(frameID, 4)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	geom := ebiten.GeoM{}
	geom.Translate(g.m.x, g.m.y)
	screen.DrawImage(g.m.img, &ebiten.DrawImageOptions{
		GeoM: geom,
	})

	geomo := ebiten.GeoM{}
	geom.Translate(g.o.x, g.o.y)
	screen.DrawImage(g.m.img, &ebiten.DrawImageOptions{
		GeoM: geomo,
	})
}

func (g *Game) Layout(_, _ int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {

	client.MHandler = mHandler
	client.Run(1, 1)

	rand.Seed(time.Now().UnixNano())

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Snake")
	game = &Game{}
	if err := game.Init(); err != nil {
		log.Fatal(err)
	}
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func mHandler(input *pb.InputData) {
	if input == nil {
		return
	}
	switch *input.Sid {
	case 1:
		game.m.y--
	case 2:
		game.m.y++
	case 3:
		game.m.x--
	case 4:
		game.m.x++
	}
}
