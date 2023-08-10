package render

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

var nullPixel color.RGBA = color.RGBA{
	0,
	0,
	0,
	0,
}

func RemoveAntiAliasing(canvas *ebiten.Image) {
	for x := 0; x < canvas.Bounds().Max.X; x++ {
		for y := 0; y < canvas.Bounds().Max.Y; y++ {
			if canvas.At(x, y) != nullPixel {
				canvas.Set(x, y, color.White)
			}
		}
	}
}
