package messages

type Message struct {
	UUID    string `json:"uuid"`
	Content string `json:"content"`
	TS      int64  `json:"ts"`
}
