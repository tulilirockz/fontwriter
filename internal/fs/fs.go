package fs

import (
	"errors"
	"image/png"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

func WriteImageToFS(user_text string, frame *ebiten.Image, base_path string, frame_counter int) error {
	var final_path string = path.Clean(base_path + "/img_frame_" + strconv.Itoa(frame_counter) + ".png")
	f, err := os.Create(final_path)
	if err != nil {
		return err
	}
	defer f.Close()
	if err = png.Encode(f, frame); err != nil {
		return err
	}
	return nil
}

func ToOutputString(base_path string, user_text *string) (string, error) {
	if base_path == "" {
		return "", errors.New("failed to parse string, string is empty")
	}
	if len(*user_text) > 20 {
		base_path += "/" + (*user_text)[:20]
	} else {
		base_path += "/" + (*user_text)
	}
	base_path = strings.ReplaceAll(base_path, "\n", "_")
	base_path = strings.ReplaceAll(base_path, " ", "_")
	return base_path, nil
}
