package ui

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"log"
)

type room struct {
}

func (r *room) draw(screen *ebiten.Image) {
	msg := fmt.Sprint(`welcome pupucar! 
enter 'c' creat game.`)

	_ = ebitenutil.DebugPrint(screen, msg)
}

func (r *room) update() error {
	if ebiten.IsKeyPressed(ebiten.KeyC) {
		log.Println("enter c !")
	}
	return nil
}

func newRoom() *room {
	r := new(room)
	return r
}
