package video

type Pixel struct {
	Red   int
	Green int
	Blue  int
	X     int
	Y     int
}

type Frame struct {
	Width    int
	Height   int
	FrameNum int
}

type VideoMeta struct {
	Duration   float64
	Resolution string
	FPS        float64
	Filename   string
}

type TermSize struct {
	Width  int
	Height int
}

type MetaRecord struct {
	Type   string `json:"type"`
	FPS    int    `json:"fps"`
	W      int    `json:"w"`
	H      int    `json:"h"`
	Source string `json:"source"`
}

type FrameRecord struct {
	I int    `json:"i"`
	W int    `json:"w"`
	H int    `json:"h"`
	F string `json:"f"`
}
