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
	render "github.com/tulilirockz/fontwriter/internal/render"
	"github.com/tulilirockz/fontwriter/internal/ui"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	defaultScreenWidth  = 640
	defaultScreenHeight = 480
)

type Game struct {
	runes         []rune
	user_text     string
	frame_counter int
}
type Flag bool

var (
	FullConfig            *configuration.ProgramConf
	TitleFont             font.Face
	RenderingCanvas       *ebiten.Image
	RenderingFrameCounter int  = 0
	FlagRender            Flag = false
	RenderingOptions      *render.Options
	OutputPath            string
)

func (g *Game) Update() error {
	g.runes = ebiten.AppendInputChars(g.runes[:0])
	g.user_text += string(g.runes)

	if repeatingKeyPressed(ebiten.KeyBackspace) {
		if len(g.user_text) >= 1 {
			g.user_text = g.user_text[:len(g.user_text)-1]
		}
	}

	if repeatingKeyPressed(ebiten.KeyF1) {
		if len(g.user_text) > 20 {
			OutputPath = g.user_text[:len(g.user_text)-20]
		} else {
			OutputPath = g.user_text
		}
		OutputPath = strings.ReplaceAll(OutputPath, "\n", "_")
		err := os.Mkdir(path.Clean(OutputPath), 0755)
		if !os.IsExist(err) && err != nil {
			log.Fatalf("Failure to create specified folder on output path: %s\n", err)
		}

		rectanglebox := text.BoundString(TitleFont, g.user_text)
		FlagRender = true
		RenderingFrameCounter = 0
		RenderingCanvas = ebiten.NewImage(
			rectanglebox.Dx(),
			rectanglebox.Dy())
	}

	// if repeatingKeyPressed(ebiten.KeyF2) {
	// screenShakeFlag = !screenShakeFlag
	//}

	if repeatingKeyPressed(ebiten.KeyEnter) || repeatingKeyPressed(ebiten.KeyNumpadEnter) {
		g.user_text += "\n"
	}

	if FlagRender {
		var wg sync.WaitGroup

		for i := 0; i < len(g.user_text); i++ {
			wg.Add(1)
			go func(i int) {
				bounding_box := text.BoundString(TitleFont, g.user_text[0:i+1])
				canvas := ebiten.NewImage(
					bounding_box.Dx(),
					bounding_box.Dy())
				text.Draw(canvas, g.user_text[0:i+1], TitleFont, 0, bounding_box.Dy(), color.White)
				if !RenderingOptions.Anti_aliasing {
					render.RemoveAntiAliasing(canvas)
				}
				err := render.WriteImageToFS(canvas, OutputPath, i)
				if err != nil {
					log.Printf("failed to encode: %v", err)
				}
				wg.Done()
			}(i)
		}
		FlagRender = false
		RenderingFrameCounter = 0
	}
	g.frame_counter++
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	t := g.user_text
	if g.frame_counter%60 < 30 {
		t += "_"
	}

	ui.TextBoundingBox(screen, t, bitmapfont.Face, color.RGBA{
		R: 255,
		G: 204,
		B: 0,
		A: 0,
	}, 3, 10, 10, defaultScreenWidth-30, defaultScreenHeight-30)
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

	RenderingOptions = &render.Options{
		Image_opt: &ebiten.DrawImageOptions{
			Filter: ebiten.FilterNearest,
		},
		Anti_aliasing: FullConfig.Output.Anti_aliasing,
	}

	RenderingOptions.Image_opt.GeoM.Scale(float64(FullConfig.Text.Scaling_factor), float64(FullConfig.Text.Scaling_factor))
	RenderingOptions.Image_opt.GeoM.Translate(0, float64(FullConfig.Text.Size))
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
