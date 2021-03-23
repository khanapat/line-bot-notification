package bot

import (
	"encoding/json"
	"fmt"
	"line-notification/client"
	"line-notification/common"
	"net/http"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type GetProfileFn func(logger *zap.Logger, xRequestID, userID string) (*Profile, error)

func NewGetProfileFn(cli *client.Client) GetProfileFn {
	return func(logger *zap.Logger, xRequestID, userID string) (*Profile, error) {
		m := map[string]string{
			common.AuthorizationHeader: fmt.Sprintf("%s %s", common.Bearer, viper.GetString("client.line-notification.channel-token")),
		}
		clientRequest := client.Request{
			URL:                 strings.Replace(viper.GetString("client.line-notification.get-profile.url"), "{userID}", userID, 1),
			Method:              http.MethodGet,
			XRequestID:          xRequestID,
			Header:              m,
			HideLogRequestBody:  viper.GetBool("client.hidebody"),
			HideLogResponseBody: viper.GetBool("client.hidebody"),
			Body:                nil,
			Logger:              logger,
		}
		clientResponse, err := cli.Do(&clientRequest)
		if err != nil {
			return nil, err
		}
		var profile Profile
		if err = json.Unmarshal(clientResponse.Body, &profile); err != nil {
			return nil, err
		}
		return &profile, nil
	}
}

type ReplyMessageClientFn func(logger *zap.Logger, xRequestID string, request *ReplyMessage) error

func NewReplyMessageClientFn(cli *client.Client) ReplyMessageClientFn {
	return func(logger *zap.Logger, xRequestID string, request *ReplyMessage) error {
		byteRequest, err := json.Marshal(&request)
		if err != nil {
			return err
		}
		m := map[string]string{
			common.AuthorizationHeader: fmt.Sprintf("%s %s", common.Bearer, viper.GetString("client.line-notification.channel-token")),
		}
		clientRequest := client.Request{
			URL:                 viper.GetString("client.line-notification.reply-message.url"),
			Method:              http.MethodPost,
			XRequestID:          xRequestID,
			Header:              m,
			HideLogRequestBody:  viper.GetBool("client.hidebody"),
			HideLogResponseBody: viper.GetBool("client.hidebody"),
			Body:                byteRequest,
			Logger:              logger,
		}
		_, err = cli.Do(&clientRequest)
		if err != nil {
			return err
		}

		return nil
	}
}

type PushMessageClientFn func(logger *zap.Logger, xRequestID string, request interface{}) error

func NewPushMessageClientFn(cli *client.Client) PushMessageClientFn {
	return func(logger *zap.Logger, xRequestID string, request interface{}) error {
		byteRequest, err := json.Marshal(&request)
		if err != nil {
			return err
		}
		m := map[string]string{
			common.AuthorizationHeader: fmt.Sprintf("%s %s", common.Bearer, viper.GetString("client.line-notification.channel-token")),
		}
		clientRequest := client.Request{
			URL:                 viper.GetString("client.line-notification.push-message.url"),
			Method:              http.MethodPost,
			XRequestID:          xRequestID,
			Header:              m,
			HideLogRequestBody:  viper.GetBool("client.hidebody"),
			HideLogResponseBody: viper.GetBool("client.hidebody"),
			Body:                byteRequest,
			Logger:              logger,
		}
		_, err = cli.Do(&clientRequest)
		if err != nil {
			return err
		}

		return nil
	}
}
