package bot

// ReplyMessage
type LineMessage struct {
	Destination string   `json:"destination"`
	Events      []Events `json:"events"`
}

type Events struct {
	ReplyToken string  `json:"replyToken"`
	Type       string  `json:"type"`
	Timestamp  int64   `json:"timestamp"`
	Source     Source  `json:"source"`
	Message    Message `json:"message"`
}

type Source struct {
	Type   string `json:"type"`
	UserID string `json:"userId"`
}

type Message struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Text string `json:"text"`
}

// PushMessage
type ReplyMessage struct {
	ReplyToken string `json:"replyToken"`
	Messages   []Text `json:"messages"`
}

type Text struct {
	Type string `json:"type"`
	Text string `json:"text"`
}
