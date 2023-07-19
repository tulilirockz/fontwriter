package configuration

type Text struct {
	Dpi            int    `json:"dpi"`
	Size           int    `json:"size"`
	Hinting        string `json:"hinting"`
	Scaling_factor int    `json:"scaling_factor"`
	Font_path      string `json:"font_path"`
}

type Output struct {
	Letters_per_frame int    `json:"letters_per_frame"`
	Anti_aliasing     bool   `json:"anti_aliasing"`
	Gif               bool   `json:"gif"`
	Path              string `json:"path"`
}

type ProgramConf struct {
	Text   Text   `json:"text"`
	Output Output `json:"output"`
}
