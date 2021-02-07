package main

import (
	"flag"
	"fmt"
	"image/color"
	"lockstepuiclient/client"
	"lockstepuiclient/game"
	"lockstepuiclient/pb"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten"
)

var mid = flag.Int("mid", 1, "help message for mid")
var roomid = flag.Int("roomid", 1, "help message for roomid")
var ip = flag.String("ip", "192.168.16.152", "lockstep server ip")

var g *Game
var frameID uint32

func getRainbowColor() color.Color {
	return game.RainbowPal[game.Id]
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
	var sid game.GameSDI
	// 移动类
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		sid |= game.Up
	} else if ebiten.IsKeyPressed(ebiten.KeyDown) {
		sid |= game.Down
	} else if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		sid |= game.Left
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		sid |= game.Right
	}

	//// 技能类
	//if ebiten.IsKeyPressed(ebiten.KeySpace) {
	//	sid |= game.Pu
	//}

	return client.SendAction(frameID, int32(sid))
}

func (g *Game) Draw(screen *ebiten.Image) {
	geom := ebiten.GeoM{}
	geom.Translate(g.m.x, g.m.y)
	screen.DrawImage(g.m.img, &ebiten.DrawImageOptions{
		GeoM: geom,
	})

	geomo := ebiten.GeoM{}
	geomo.Translate(g.o.x, g.o.y)
	screen.DrawImage(g.m.img, &ebiten.DrawImageOptions{
		GeoM: geomo,
	})
}

func (g *Game) Layout(_, _ int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	flag.Parse()

	rand.Seed(time.Now().UnixNano())
	game.RoomID = uint64(*roomid)
	game.Id = uint64(*mid)
	log.Print("welcome to pupucar ,mid = ", game.Id, ", roomid = ", game.RoomID)

	client.MHandler = mHandler
	client.Run(game.RoomID, game.Id, *ip)

	s := <-client.StartChan
	if s {
		ebiten.SetWindowSize(640, 480)
		ebiten.SetWindowTitle(fmt.Sprintf("pupucar id:%d room:%d", game.Id, game.RoomID))
		g = &Game{}
		if err := g.Init(); err != nil {
			log.Fatal(err)
		}
		if err := ebiten.RunGame(g); err != nil {
			log.Fatal(err)
		}
	}
}

func mHandler(input []*pb.InputData) {
	if input == nil || len(input) == 0 {
		return
	}

	for _, i := range input {
		handleInput(i)
	}

}

func handleInput(input *pb.InputData) {
	if input == nil {
		return
	}
	sid := game.GameSDI(*input.Sid)

	if *input.Id == game.Id {
		if (sid & game.Up) != 0 {
			g.m.y -= 1
		} else if (sid & game.Down) != 0 {
			g.m.y += 1
		} else if (sid & game.Left) != 0 {
			g.m.x -= 1
		} else if (sid & game.Right) != 0 {
			g.m.x += 1
		}
	} else {
		if (sid & game.Up) != 0 {
			g.o.y -= 1
		} else if (sid & game.Down) != 0 {
			g.o.y += 1
		} else if (sid & game.Left) != 0 {
			g.o.x -= 1
		} else if (sid & game.Right) != 0 {
			g.o.x += 1
		}
	}
}
