package api

type PrintItem struct {
	Type    string `json:"type"`
	Content string `json:"content,omitempty"`
	Align   string `json:"align,omitempty"`
	Count   int    `json:"count,omitempty"`
}

type PrintRequest struct {
	Items []PrintItem `json:"items"`
}
