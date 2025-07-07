package printer

type ItemType string

const (
	Text   ItemType = "text"
	QRCode ItemType = "qrcode"
	Blank  ItemType = "blank"
	Line   ItemType = "line"
	Cut    ItemType = "cut"
)

type Alignment string

const (
	AlignLeft   Alignment = "left"
	AlignCenter Alignment = "center"
	AlignRight  Alignment = "right"
)

type PrintItem struct {
	Type    ItemType  `json:"type"`
	Content string    `json:"content,omitempty"`
	Align   Alignment `json:"align,omitempty"`
	Count   int       `json:"count,omitempty"`

	Bold         bool `json:"bold,omitempty"`
	Underline    bool `json:"underline,omitempty"`
	DoubleHeight bool `json:"doubleHeight,omitempty"`
	DoubleWidth  bool `json:"doubleWidth,omitempty"`
}

type PrintRequest struct {
	Items []PrintItem `json:"items"`
}