package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
)

func TextBoundingBox(screen *ebiten.Image, label string, font font.Face, usercolor color.Color, padding, x0, y0, x1, y1 float32) {
	vector.DrawFilledRect(screen, x0, y0, x1, y1, usercolor, false)                                       
	vector.DrawFilledRect(screen, x0+padding, y0+padding, x1-padding*2, y1-padding*2, color.White, false)

	text.Draw(screen, label, font, int(x0+padding+10), int(y0+padding+10), color.Black)
}
