package bot

import (
	"fmt"
	"line-notification/common"
	"line-notification/handler"
	"line-notification/response"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type Data struct {
	Restaurant string   `json:"restaurant"`
	Score      int      `json:"score"`
	VotedBy    []string `json:"voteBy"`
}

var listData []Data

type bothandler struct {
	GetProfileFn         GetProfileFn
	ReplyMessageClientFn ReplyMessageClientFn
	PushMessageClientFn  PushMessageClientFn
}

func NewBotHandler(getProfileFn GetProfileFn, replyMessageClientFn ReplyMessageClientFn, pushMessageClientFn PushMessageClientFn) *bothandler {
	return &bothandler{
		GetProfileFn:         getProfileFn,
		ReplyMessageClientFn: replyMessageClientFn,
		PushMessageClientFn:  pushMessageClientFn,
	}
}

func (s *bothandler) ReplyBot(c *handler.Ctx) error {
	var req LineMessage
	if err := c.BodyParser(&req); err != nil {
		c.Log().Error(err.Error())
		return c.Status(http.StatusBadRequest).JSON(response.NewErrResponse(response.ErrInvalidRequestCode, response.ErrInvalidRequestMessageEn))
	}
	if len(req.Events) != 0 {
		profile, err := s.GetProfileFn(c.Log(), string(c.Request().Header.Peek(common.XRequestID)), req.Events[0].Source.UserID)
		if err != nil {
			c.Log().Error(err.Error())
			return c.Status(http.StatusBadRequest).JSON(response.NewErrResponse(response.ErrThirdPartyCode, response.ErrThirdPartyMessageEn))
		}

		c.Log().Debug(profile.DisplayName)
		c.Log().Debug(profile.UserID)

		var text Text
		var message ReplyMessage
		switch req.Events[0].Type {
		case "message":
			switch req.Events[0].Message.Type {
			case "text":
				var textMessage string
				switch strings.Split(req.Events[0].Message.Text, "")[0] {
				case "!":
					event := strings.Split(req.Events[0].Message.Text, " ")
					switch event[0] {
					case "!list":
						textMessage = fmt.Sprintf("--- ตารางคะแนน ---\n")
						for _, value := range listData {
							textMessage += fmt.Sprintf("%s: %d\n", value.Restaurant, value.Score)
						}
						textMessage += fmt.Sprintf("----------------------")
					case "!add":
						if len(event) != 2 {
							textMessage = fmt.Sprintf("ส่งตัวแปรมาไม่ตรงโคต้าวุย ปวดหัว")
						} else {
							var isDup bool
							for _, value := range listData {
								if value.Restaurant == event[1] {
									isDup = true
								}
							}
							if !isDup {
								data := Data{
									Restaurant: event[1],
									Score:      0,
								}
								listData = append(listData, data)
								textMessage = fmt.Sprintf("คุณ '%s' ได้เพิ่มรายการ '%s' สำหรับการลงคะแนน", profile.DisplayName, event[1])
							} else {
								textMessage = fmt.Sprintf("รายการ '%s' ลงทะเบียนไปแล้วงับ จะลงซ้ำทำเพื่อ (- -)", event[1])
							}
						}
					case "!vote":
						if len(event) != 2 {
							textMessage = fmt.Sprintf("ส่งตัวแปรมาไม่ตรงโคต้าวุย ปวดหัว")
						} else {
							var score int
							for index, value := range listData {
								if value.Restaurant == event[1] {
									score = value.Score + 1
									listData[index].Score = score
								}
							}
							if score == 0 {
								textMessage = fmt.Sprintf("มันยังไม่มีรายการ '%s' นะ กรุณาไป !add ก่อน", event[1])
							} else {
								textMessage = fmt.Sprintf("คุณ '%s' ได้ลงคะแนนให้กับ '%s' เรียบร้อยแล้ว", profile.DisplayName, event[1])
							}
						}
					case "!reset":
						listData = []Data{}
						textMessage = fmt.Sprintf("คุณ '%s' ได้รีเซตการลงคะแนนใหม่", profile.DisplayName)
					case "!help":
						textMessage = fmt.Sprintf("คำสั่งต่างๆ\n" +
							"!list -> สำหรับแสดงรายการที่ลงทะเบียน\n" +
							"!add {v} -> สำหรับการเพิ่มรายการในตาราง\n" +
							"!vote {v} -> สำหรับลงคะแนน\n" +
							"!reset -> สำหรับรีเซตข้อมูลทั้งหมด")
					default:
						textMessage = fmt.Sprintf("ไม่มีคำสั่ง '%s' โว้ย\nอยากได้อะไรไปติดต่อทรัสเอา :D", event[0])
					}
				default:
					textMessage = fmt.Sprintf("ว่ายังไง <3 '%s' \nส่งข้อความ '%s' มาต้องการอะไรหยอ?", profile.DisplayName, req.Events[0].Message.Text)
				}
				text = Text{
					Type: "text",
					Text: textMessage,
				}
				message = ReplyMessage{
					ReplyToken: req.Events[0].ReplyToken,
					Messages: []Text{
						text,
					},
				}
			default:
				text = Text{
					Type: "sticker",
				}
			}
		default:

		}

		if err := s.ReplyMessageClientFn(c.Log(), string(c.Request().Header.Peek(common.XRequestID)), &message); err != nil {
			c.Log().Error(err.Error())
			return c.Status(http.StatusBadRequest).JSON(response.NewErrResponse(response.ErrThirdPartyCode, response.ErrThirdPartyMessageEn))
		}
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{})
}

func (s *bothandler) Voting(c *handler.Ctx) error {
	bubble := Bubble{
		Type:      "bubble",
		Direction: "ltr",
		Header: Header{
			Type:            "box",
			Layout:          "horizontal",
			BackgroundColor: "#000000FF",
			BorderColor:     "#FFFFFFFF",
			Contents: []HeaderContent{
				{
					Type:     "text",
					Text:     "Restaurant",
					Weight:   "bold",
					Size:     "xl",
					Color:    "#EADD01FF",
					Align:    "center",
					Gravity:  "center",
					Margin:   "lg",
					Wrap:     true,
					Style:    "normal",
					Contents: []int{},
				},
			},
		},
		Hero: Hero{
			Type:            "image",
			URL:             "https://i.postimg.cc/Xv44qGRt/best-restaurants-in-bangkok.jpg",
			Align:           "center",
			Gravity:         "center",
			Size:            "full",
			AspectRatio:     "1.51:1",
			AspectMode:      "cover",
			BackgroundColor: "#0C0000FF",
		},
		Body: Body{
			Type:            "box",
			Layout:          "vertical",
			BackgroundColor: "#000000FF",
			Contents: []BodyContent{
				{
					Type:   "box",
					Layout: "horizontal",
					Contents: []ContentContent{
						{
							Type:    "text",
							Text:    "Uncle house",
							Weight:  "bold",
							Color:   "#E9DCDCFF",
							Align:   "center",
							Gravity: "center",
							Action: &Action{
								Type:  "message",
								Label: "Voting",
								Text:  "Uncle house",
							},
							Contents: []int{},
						},
						{
							Type:   "separator",
							Color:  "#6C6363FF",
							Action: nil,
						},
						{
							Type:    "text",
							Text:    "Crispy pork alley",
							Weight:  "bold",
							Color:   "#E9DCDCFF",
							Align:   "center",
							Gravity: "center",
							Action: &Action{
								Type:  "message",
								Label: "Voting",
								Text:  "Crispy pork alley",
							},
							Contents: []int{},
						},
					},
				},
				{
					Type:  "separator",
					Color: "#6C6363FF",
				},
				{
					Type:   "box",
					Layout: "horizontal",
					Contents: []ContentContent{
						{
							Type:    "text",
							Text:    "Curvy suki",
							Weight:  "bold",
							Color:   "#E9DCDCFF",
							Align:   "center",
							Gravity: "center",
							Action: &Action{
								Type:  "message",
								Label: "Voting",
								Text:  "Curvy suki",
							},
							Contents: []int{},
						},
						{
							Type:  "separator",
							Color: "#6C6363FF",
						},
						{
							Type:    "text",
							Text:    "Papaya salad",
							Weight:  "bold",
							Color:   "#E9DCDCFF",
							Align:   "center",
							Gravity: "center",
							Action: &Action{
								Type:  "message",
								Label: "Voting",
								Text:  "Papaya salad",
							},
							Contents: []int{},
						},
					},
				},
			},
		},
	}

	pushBubble := PushBubble{
		// To: "U7f23c5963e6ef29e206e23d7b785660f",
		To: "Cfb1b2b20854a13490e67a4686a96fa61",
		Messages: []Message{
			{
				Type:     "flex",
				AltText:  "This is a Flex",
				Contents: bubble,
			},
		},
	}

	_ = PushText{
		To: "U7f23c5963e6ef29e206e23d7b785660f",
		Message: []TextMessage{
			{
				Type: "text",
				Text: "Hello",
			},
		},
	}

	if err := s.PushMessageClientFn(c.Log(), string(c.Request().Header.Peek(common.XRequestID)), &pushBubble); err != nil {
		c.Log().Error(err.Error())
		return c.Status(http.StatusBadRequest).JSON(response.NewErrResponse(response.ErrThirdPartyCode, response.ErrThirdPartyMessageEn))
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{})
}

func (s *bothandler) Test(c *handler.Ctx) error {
	var req Test
	if err := c.BodyParser(&req); err != nil {
		c.Log().Error(err.Error())
		return c.Status(http.StatusBadRequest).JSON(response.NewErrResponse(response.ErrInvalidRequestCode, response.ErrInvalidRequestMessageEn))
	}

	c.Log().Info(fmt.Sprintf("%s", req.Name))
	c.Log().Info(fmt.Sprintf("%d", req.Number))

	a := strings.Split(req.Name, "")

	c.Log().Info(fmt.Sprintf("%s %s", a[0], a[1]))
	return c.Status(http.StatusOK).JSON(response.NewResponse(response.SuccessCode, response.SuccessMessageEn, nil))
}
