package bot

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"line-notification/handler"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type botHandler struct {
}

func NewBotHandler() *botHandler {
	return &botHandler{}
}

const ChannalToken string = "cdsyUJnmmWbfU8zHbWb4pZVjIw8jrMIXOceX0zBP8e/keH6KAnt4TyG6ZCXGstYP03m68Q1BZOU/7/DvmTyKSDR4EnxLxyAuVq94zQjT6HrjI0bMf9XW5spws2hwmD2ebD+OGHKrVVR3u9i2VlrXCAdB04t89/1O/w1cDnyilFU="

func (s *botHandler) ReplyBot(c *handler.Ctx) error {
	c.Log().Info(string(c.Body()))

	var req LineMessage

	if err := c.BodyParser(&req); err != nil {
		c.Log().Error(err.Error())
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"Error": err.Error(),
		})
	}

	if len(req.Events) != 0 {
		text := Text{
			Type: "text",
			Text: "ข้อความเข้ามา : " + req.Events[0].Message.Text + " ยินดีต้อนรับ : ",
		}

		message := ReplyMessage{
			ReplyToken: req.Events[0].ReplyToken,
			Messages: []Text{
				text,
			},
		}

		go replyMessageLine(message)
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{})
}

func replyMessageLine(Message ReplyMessage) error {
	value, _ := json.Marshal(Message)

	url := "https://api.line.me/v2/bot/message/reply"

	var jsonStr = []byte(value)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+ChannalToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	log.Println("response Status:", resp.Status)
	log.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	log.Println("response Body:", string(body))

	return err
}
