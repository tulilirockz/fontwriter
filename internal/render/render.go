package render

import (
	"image/color"
	"image/png"
	"log"
	"os"
	"path"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
)

var nullPixel color.RGBA = color.RGBA{
	0,
	0,
	0,
	0,
}

type Options struct {
	Image_opt     *ebiten.DrawImageOptions
	Anti_aliasing bool
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

func WriteImageToFS(frame *ebiten.Image, base_path string, frame_counter int) error {
	f, err := os.Create(path.Clean(base_path + "/img_frame_" + strconv.Itoa(frame_counter) + ".png"))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if err = png.Encode(f, frame); err != nil {
		return err
	}
	return nil
}

func ByteToInteger(b bool) int {
	if b {
		return 1
	}
	return 0
}
