package configuration

type Text struct {
	Dpi            int     `json:"dpi"`
	Size           int     `json:"size"`
	Hinting        string  `json:"hinting"`
	Scaling_factor float32 `json:"scaling_factor"`
	Font_path      string  `json:"font_path"`
	Anti_aliasing     bool   `json:"anti_aliasing"`
}

type Output struct {
	Letters_per_frame int    `json:"letters_per_frame"`
	Gif               bool   `json:"gif"`
	Path              string `json:"path"`
}

type ProgramConf struct {
	Text   Text   `json:"text"`
	Output Output `json:"output"`
}
