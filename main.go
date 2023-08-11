package main

import (
	"image/color"
	"log"
	"path"
	"runtime"
	"strings"
	"sync"

	"os"

	"github.com/hajimehoshi/bitmapfont/v2"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/pelletier/go-toml/v2"
	"github.com/tulilirockz/fontwriter/internal/configuration"
	"github.com/tulilirockz/fontwriter/internal/fs"
	render "github.com/tulilirockz/fontwriter/internal/render"
	"github.com/tulilirockz/fontwriter/internal/ui"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	defaultScreenWidth  = 640
	defaultScreenHeight = 480
)

var (
	FullConfig *configuration.ProgramConf
	TitleFont  font.Face
)

type Game struct {
	runes         []rune
	user_text     string
	frame_counter int
}
type Flag bool

func (g *Game) Update() error {
	g.runes = ebiten.AppendInputChars(g.runes[:0])
	g.user_text += string(g.runes)

	if repeatingKeyPressed(ebiten.KeyBackspace) {
		if len(g.user_text) >= 1 {
			g.user_text = g.user_text[:len(g.user_text)-1]
		}
	}

	if repeatingKeyPressed(ebiten.KeyF1) {
		out_path, err := fs.ToPathString(FullConfig.Output.Path, &g.user_text)
		if err != nil {
			log.Printf("Failed translating output string")
			return err
		}

		err = os.MkdirAll(path.Clean(out_path), 0755)
		if !os.IsExist(err) && err != nil {
			return err
		}

		filtered_text := strings.ReplaceAll(g.user_text, "\n", "_")
		filtered_text = strings.ReplaceAll(filtered_text, " ", "_")

		var waitgrp sync.WaitGroup

		for i := 0; i < len(g.user_text); i++ {
			waitgrp.Add(1)
			go func(i int) {
				bounding_box := text.BoundString(TitleFont, filtered_text[0:i+1])
				canvas := ebiten.NewImage(
					bounding_box.Dx(),
					bounding_box.Dy())
				text.Draw(canvas, (filtered_text)[0:i+1], TitleFont, 0, bounding_box.Dy(), color.White)
				if !FullConfig.Text.Anti_aliasing {
					render.RemoveAntiAliasing(canvas)
				}
				err := fs.WriteImageToFS(g.user_text, canvas, out_path, i)
				if err != nil {
					log.Printf("Failed writing image to disk, %v", err)
				}
				waitgrp.Done()
			}(i)
		}
	}

	if repeatingKeyPressed(ebiten.KeyEnter) || repeatingKeyPressed(ebiten.KeyNumpadEnter) {
		g.user_text += "\n"
	}

	g.frame_counter++
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	t := g.user_text
	if g.frame_counter%60 < 30 {
		t += "_"
	}

	const (
		PADDING = 3
		OFFSET  = 10
	)

	var default_color color.Color = color.RGBA{
		R: 72,
		G: 12,
		B: 122,
		A: 1,
	}

	ui.TextBoundingBox(screen, t, bitmapfont.Face, default_color, 3, 10, 10, defaultScreenWidth-30, defaultScreenHeight-30)
	text.Draw(screen, "Instructions:\nF1: Render separate png files\nF2: Render a .gif file", bitmapfont.Face, int(10+PADDING+OFFSET), int(400+PADDING+OFFSET), color.Black)
}

func init() {
	configbytes, err := os.ReadFile("config.toml")
	if err != nil {
		log.Fatalf("Failed reading configuration file: %s", err)
	}
	err = toml.Unmarshal(configbytes, &FullConfig)
	if err != nil {
		log.Panic(err)
	}

	font_bytes, err := os.ReadFile(FullConfig.Text.Font_path)
	if err != nil {
		log.Fatal(err)
	}

	font_type_parsed, err := opentype.Parse(font_bytes)
	if err != nil {
		log.Fatal(err)
	}

	var font_hinting font.Hinting = font.HintingNone

	switch FullConfig.Text.Hinting {
	case "vertical":
		font_hinting = font.HintingVertical
	case "full":
		font_hinting = font.HintingFull
	default:
		font_hinting = font.HintingNone
	}

	TitleFont, _ = opentype.NewFace(font_type_parsed, &opentype.FaceOptions{
		Size:    float64(FullConfig.Text.Size),
		DPI:     float64(FullConfig.Text.Dpi),
		Hinting: font_hinting,
	})
}

func main() {
	g := &Game{
		user_text:     "",
		frame_counter: 0,
	}
	runtime.GOMAXPROCS(runtime.NumCPU())
	ebiten.SetWindowTitle("Type your thing here")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

func repeatingKeyPressed(key ebiten.Key) bool {
	const (
		delay    = 30
		interval = 3
	)
	d := inpututil.KeyPressDuration(key)
	if d == 1 {
		return true
	}
	if d >= delay && (d-delay)%interval == 0 {
		return true
	}
	return false
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return defaultScreenWidth, defaultScreenHeight
}
