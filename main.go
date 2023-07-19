package main

import (
	"image/color"
	"log"
	"path"

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
	fullconfig            *configuration.ProgramConf
	titleFont             font.Face
	renderingCanvas       *ebiten.Image
	renderingFrameCounter int  = 0
	flagRender            Flag = false
	renderingOptions      *render.Options
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
		err := os.Mkdir(path.Clean(fullconfig.Output.Path), 0755)
		if !os.IsExist(err) && err != nil {
			log.Fatalf("Failure to create specified folder on output path: %s\n", err)
		}

		rectanglebox := text.BoundString(titleFont, g.user_text)
		flagRender = true
		renderingFrameCounter = 0
		renderingCanvas = ebiten.NewImage(
			rectanglebox.Dx(),
			rectanglebox.Dy())
	}

	// if repeatingKeyPressed(ebiten.KeyF2) {
	// screenShakeFlag = !screenShakeFlag
	//}

	if repeatingKeyPressed(ebiten.KeyEnter) || repeatingKeyPressed(ebiten.KeyNumpadEnter) {
		g.user_text += "\n"
	}

	if flagRender {
		if renderingFrameCounter == len(g.user_text) {
			flagRender = false
			renderingFrameCounter = 0
			return nil
		}

		rectanglebox := text.BoundString(titleFont, "h")
		text.Draw(renderingCanvas, g.user_text[0:renderingFrameCounter+1], titleFont, 0, rectanglebox.Dy(), color.White)

		if !renderingOptions.Anti_aliasing {
			render.RemoveAntiAliasing(renderingCanvas)
		}

		err := render.WriteImageToFS(renderingCanvas, fullconfig.Output.Path, renderingFrameCounter)
		if err != nil {
			log.Printf("failed to encode: %v", err)
		}

		renderingFrameCounter++
	}
	g.frame_counter++
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	t := g.user_text
	if g.frame_counter%60 < 30 {
		t += "_"
	}

	ui.BoundingBox(screen, t, bitmapfont.Face, 3, 10, 10, defaultScreenWidth-30, defaultScreenHeight-30)

	if flagRender {	
		boundingbox := text.BoundString(bitmapfont.Face, "Rendering!")
		ui.BoundingBox(screen, "Rendering!", bitmapfont.Face, 3, 100, 50, float32(100+boundingbox.Dx()), float32(50+boundingbox.Dy()))
	}
}

func init() {
	configbytes, err := os.ReadFile("config.toml")
	if err != nil {
		log.Fatalf("Failed reading configuration file: %s", err)
	}
	err = toml.Unmarshal(configbytes, &fullconfig)
	if err != nil {
		log.Panic(err)
	}

	font_bytes, err := os.ReadFile(fullconfig.Text.Font_path)
	if err != nil {
		log.Fatal(err)
	}

	font_type_parsed, err := opentype.Parse(font_bytes)
	if err != nil {
		log.Fatal(err)
	}

	var font_hinting font.Hinting = font.HintingNone

	switch fullconfig.Text.Hinting {
	case "vertical":
		font_hinting = font.HintingVertical
	case "full":
		font_hinting = font.HintingFull
	default:
		font_hinting = font.HintingNone
	}

	titleFont, _ = opentype.NewFace(font_type_parsed, &opentype.FaceOptions{
		Size:    float64(fullconfig.Text.Size),
		DPI:     float64(fullconfig.Text.Dpi),
		Hinting: font_hinting,
	})

	renderingOptions = &render.Options{
		Image_opt: &ebiten.DrawImageOptions{
			Filter: ebiten.FilterNearest,
		},
		Anti_aliasing: fullconfig.Output.Anti_aliasing,
	}

	renderingOptions.Image_opt.GeoM.Scale(float64(fullconfig.Text.Scaling_factor), float64(fullconfig.Text.Scaling_factor))
	renderingOptions.Image_opt.GeoM.Translate(0, float64(fullconfig.Text.Size))
}

func main() {
	g := &Game{
		user_text:     "",
		frame_counter: 0,
	}

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
