package bot

// ReplyMessage
type LineMessage struct {
	Destination string   `json:"destination"`
	Events      []Events `json:"events"`
}

type Events struct {
	ReplyToken string      `json:"replyToken"`
	Type       string      `json:"type"`
	Timestamp  int64       `json:"timestamp"`
	Source     Source      `json:"source"`
	Message    TestMessage `json:"message"`
}

type Source struct {
	Type   string `json:"type"`
	UserID string `json:"userId"`
}

type TestMessage struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Text string `json:"text"`
}

// Profile
type Profile struct {
	UserID        string `json:"userId"`
	DisplayName   string `json:"displayName"`
	PictureURL    string `json:"pictureUrl"`
	StatusMessage string `json:"statusMessage"`
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

// Test
type Test struct {
	Name   string `json:"name"`
	Number int    `json:"number"`
}

// PushText
type PushText struct {
	To      string        `json:"to"`
	Message []TextMessage `json:"messages"`
}

type TextMessage struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// PushBubble
type PushBubble struct {
	To       string    `json:"to"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Type     string      `json:"type"`
	AltText  string      `json:"altText"`
	Contents interface{} `json:"contents"`
}

// BubbleVoting
type Bubble struct {
	Type      string `json:"type"`
	Direction string `json:"direction"`
	Header    Header `json:"header"`
	Hero      Hero   `json:"hero"`
	Body      Body   `json:"body"`
}

type Body struct {
	Type            string        `json:"type"`
	Layout          string        `json:"layout"`
	BackgroundColor string        `json:"backgroundColor"`
	Contents        []BodyContent `json:"contents"`
}

type BodyContent struct {
	Type     string           `json:"type"`
	Layout   string           `json:"layout,omitempty"`
	Contents []ContentContent `json:"contents,omitempty"`
	Color    string           `json:"color,omitempty"`
}

type ContentContent struct {
	Type     string      `json:"type"`
	Text     string      `json:"text,omitempty"`
	Weight   string      `json:"weight,omitempty"`
	Color    string      `json:"color"`
	Align    string      `json:"align,omitempty"`
	Gravity  string      `json:"gravity,omitempty"`
	Action   *Action     `json:"action,omitempty"`
	Contents interface{} `json:"contents,omitempty"`
}

type Action struct {
	Type  string `json:"type"`
	Label string `json:"label"`
	Text  string `json:"text"`
}

type Header struct {
	Type            string          `json:"type"`
	Layout          string          `json:"layout"`
	BackgroundColor string          `json:"backgroundColor"`
	BorderColor     string          `json:"borderColor"`
	Contents        []HeaderContent `json:"contents"`
}

type HeaderContent struct {
	Type     string      `json:"type"`
	Text     string      `json:"text"`
	Weight   string      `json:"weight"`
	Size     string      `json:"size"`
	Color    string      `json:"color"`
	Align    string      `json:"align"`
	Gravity  string      `json:"gravity"`
	Margin   string      `json:"margin"`
	Wrap     bool        `json:"wrap"`
	Style    string      `json:"style"`
	Contents interface{} `json:"contents"`
}

type Hero struct {
	Type            string `json:"type"`
	URL             string `json:"url"`
	Align           string `json:"align"`
	Gravity         string `json:"gravity"`
	Size            string `json:"size"`
	AspectRatio     string `json:"aspectRatio"`
	AspectMode      string `json:"aspectMode"`
	BackgroundColor string `json:"backgroundColor"`
}
