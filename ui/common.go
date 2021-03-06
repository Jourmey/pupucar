package ui

import "image/color"

var RainbowPal = []color.RGBA{
	{0xff, 0x00, 0x00, 0xff},
	{0xff, 0x7f, 0x00, 0xff},
	{0xff, 0xff, 0x00, 0xff},
	{0x00, 0xff, 0x00, 0xff},
	{0x00, 0x00, 0xff, 0xff},
	{0x4b, 0x00, 0x82, 0xff},
	{0x8f, 0x00, 0xff, 0xff},
}

type GameSDI int32

const (
	Up GameSDI = 1 << iota
	Down
	Left
	Right

	Pu
)
