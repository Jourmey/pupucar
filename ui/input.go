package ui

import (
	"github.com/hajimehoshi/ebiten"
	"lockstepuiclient/client"
)

type input struct {
}

func newInput() *input {
	return new(input)
}

var frameID uint32

func (i *input) update() error {
	frameID++
	var sid GameSDI
	// 移动类
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		sid |= Up
	} else if ebiten.IsKeyPressed(ebiten.KeyDown) {
		sid |= Down
	} else if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		sid |= Left
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) {
		sid |= Right
	}

	if sid == 0 {
		return nil
	}

	return client.SendAction(frameID, int32(sid))
}
