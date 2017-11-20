package utils

type ImageMagicPreparer interface {
	Prepare()
}

type Object interface {
	getType() string
}

type PngPage struct {
	B64            string         `json:"b64" binding:"required"`
	Page           int            `json:"page" binding:"required"`
	CanvasElements CanvasElements `json:"canvasElements" binding:"required"`
}

type PngToPdf struct {
	Data  []DataType            `json:"data" binding:"required"`
	Pages []PngPageWithElements `json:"pages" binding:"required"`
}

type PngPageWithElements struct {
	Page           int            `json:"page" binding:"required"`
	CanvasElements CanvasElements `json:"canvasElements" binding:"required"`
}

type DataType struct {
	Placeholder string `json:"placeholder" binding:"required"`
	Value       string `json:"value" binding:"required"`
}

type CanvasElements struct {
	//BackgroundImage Image  `json:"backgroundImage" `
	Objects []Text `json:"objects"`
}

type BaseObject struct {
	Angle                    float64 `json:"angle"`
	CharSpacing              int     `json:"charSpacing"`
	ClipTo                   string  `json:"clipTo"`
	Fill                     string  `json:"fill"`
	FillRule                 string  `json:"fillRule"`
	FlipX                    bool    `json:"flipX"`
	FlipY                    bool    `json:"flipY"`
	Opacity                  float64 `json:"opacity"`
	OriginX                  string  `json:"originX"`
	OriginY                  string  `json:"originY"`
	ScaleX                   float64 `json:"scaleX"`
	ScaleY                   float64 `json:"scaleY"`
	Shadow                   string  `json:"shadow"`
	SkewX                    int     `json:"skewX"`
	SkewY                    int     `json:"skewY"`
	Stroke                   string  `json:"stroke"`
	StrokeDashArray          string  `json:"strokeDashArray"`
	StrokeLineCap            string  `json:"strokeLineCap"`
	StrokeLineJoin           string  `json:"strokeLineJoin"`
	StrokeMiterLimit         int     `json:"strokeMiterLimit"`
	StrokeWidth              int     `json:"strokeWidth"`
	Top                      float64 `json:"top"`
	TransformMatrix          string  `json:"transformMatrix"`
	TypeCanvas               string  `json:"type" binding:"required"`
	Visible                  bool    `json:"visible"`
	Width                    float64 `json:"width"`
	GlobalCompositeOperation string  `json:"globalCompositeOperation"`
	Height                   float64 `json:"height"`
	Left                     float64 `json:"left"`
}

type Image struct {
	BaseObject
	AlignX      string `json:"alignX"`
	AlignY      string `json:"alignY"`
	MeetOrSlice string `json:"meetOrSlice"`
}

func (i Image) getType() string {
	return "image"
}

type Text struct {
	BaseObject
	BackgroundColor     string      `json:"backgroundColor"`
	FontFamily          string      `json:"fontFamily"`
	FontSize            string      `json:"fontSize"`
	FontStyle           string      `json:"fontStyle"`
	FontWeight          int         `json:"fontWeight"`
	LineHeight          float64     `json:"lineHeight"`
	Styles              interface{} `json:"styles"`
	Text                string      `json:"text"`
	TextAlign           string      `json:"textAlign"`
	TextBackgroundColor string      `json:"textBackgroundColor"`
	TextDecoration      string      `json:"textDecoration"`
}

func (t Text) getType() string {
	return "text"
}

func (t *Text) Prepare() {
	t.Top += 50
	//font-zie must be converted from px to dpi
	//we convert pdf to 300dpi solution for high quality printing
	//ratio for 300dpi is 1 pixel/cm = 2.54 dpi
	//t.FontSize = t.FontSize * 2.54
	if t.Fill == "" {
		t.Fill = "#000"
	}
}

type PngWithProps struct {
	Page           int            `json:"page" binding:"required"`
	CanvasElements CanvasElements `json:"canvasElements" binding:"required"`
	Generated      string         `json:"generated" binding:"required"`
	Original       string         `json:"original" binding:"required"`
}
