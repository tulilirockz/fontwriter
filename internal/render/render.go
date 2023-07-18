package render

import (
	"image/color"
	"image/png"
	"log"
	"os"
	"path"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
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
	Shakeit       bool
}

func ByteToInteger(b bool) int {
	if b {
		return 1
	}
	return 0
}

// Modifies the canvas image!
func RenderTextToCanvas(input_text string, frame_counter *int, canvas *ebiten.Image, target_font font.Face, options Options) error {
	text.DrawWithOptions(canvas, input_text[0:*frame_counter+1], target_font, options.Image_opt)
	log.Printf("%s", input_text[0:*frame_counter+1])
	log.Printf("%v", options.Image_opt)

	if options.Anti_aliasing {
		return nil
	}

	for x := 0; x < canvas.Bounds().Max.X; x++ {
		for y := 0; y < canvas.Bounds().Max.Y; y++ {
			if canvas.At(x, y) != nullPixel {
				canvas.Set(x, y, color.White)
			}
		}
	}
	return nil
}

func WriteImageToFS(frame *ebiten.Image, base_path string, frame_counter int) error {
	f, err := os.Create(path.Clean(base_path + "/img_frame_" + strconv.Itoa(frame_counter) + ".png"))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if err = png.Encode(f, frame); err != nil {
		log.Printf("failed to encode: %v", err)
		return err
	}
	return nil
}
