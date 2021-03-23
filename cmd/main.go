package main

import (
	"fmt"
	"line-notification/bot"
	"line-notification/client"
	"line-notification/handler"
	"line-notification/logz"
	"line-notification/middleware"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"
	_ "time/tzdata"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func init() {
	runtime.GOMAXPROCS(1)
	initTimezone()
	initViper()
}

func main() {
	app := fiber.New(fiber.Config{
		StrictRouting: true,
		CaseSensitive: true,
		Immutable:     true,
		ReadTimeout:   viper.GetDuration("app.timeout"),
		WriteTimeout:  viper.GetDuration("app.timeout"),
		IdleTimeout:   viper.GetDuration("app.timeout"),
	})

	logger, err := logz.NewLogConfig()
	if err != nil {
		log.Fatal(err)
	}

	httpClient := client.NewClient()

	middle := middleware.NewMiddleware(logger)

	line := app.Group(viper.GetString("app.context"))

	line.Use(middle.JSONMiddleware())
	line.Use(middle.LoggingMiddleware())

	line.Post("/webhook", handler.Helper(bot.NewBotHandler(
		bot.NewGetProfileFn(httpClient),
		bot.NewReplyMessageClientFn(httpClient),
		bot.NewPushMessageClientFn(httpClient),
	).ReplyBot, logger))

	line.Get("/push", handler.Helper(bot.NewBotHandler(
		bot.NewGetProfileFn(httpClient),
		bot.NewReplyMessageClientFn(httpClient),
		bot.NewPushMessageClientFn(httpClient),
	).Voting, logger))

	line.Post("/test", handler.Helper(bot.NewBotHandler(
		bot.NewGetProfileFn(httpClient),
		bot.NewReplyMessageClientFn(httpClient),
		bot.NewPushMessageClientFn(httpClient),
	).Test, logger))

	line.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	logger.Info(fmt.Sprintf("⇨ http server started on [::]:%s", viper.GetString("app.port")))

	go func() {
		if err := app.Listen(fmt.Sprintf(":%s", viper.GetString("app.port"))); err != nil {
			logger.Info(err.Error())
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	select {
	case <-c:
		logger.Info("terminating: by signal")
	}

	app.Shutdown()

	logger.Info("shutting down")
	os.Exit(0)
}

func initViper() {
	viper.SetDefault("app.name", "line-notification-api")
	viper.SetDefault("app.port", "9090")
	viper.SetDefault("app.timeout", "60s")
	viper.SetDefault("app.context", "/line")

	viper.SetDefault("log.level", "debug")
	viper.SetDefault("log.env", "dev")

	viper.SetDefault("client.timeout", "60s")
	viper.SetDefault("client.hidebody", true)
	viper.SetDefault("client.line-notification.channel-token", "cdsyUJnmmWbfU8zHbWb4pZVjIw8jrMIXOceX0zBP8e/keH6KAnt4TyG6ZCXGstYP03m68Q1BZOU/7/DvmTyKSDR4EnxLxyAuVq94zQjT6HrjI0bMf9XW5spws2hwmD2ebD+OGHKrVVR3u9i2VlrXCAdB04t89/1O/w1cDnyilFU=")
	viper.SetDefault("client.line-notification.get-profile.url", "https://api.line.me/v2/bot/profile/{userID}")
	viper.SetDefault("client.line-notification.reply-message.url", "https://api.line.me/v2/bot/message/reply")
	viper.SetDefault("client.line-notification.push-message.url", "https://api.line.me/v2/bot/message/push")

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

func initTimezone() {
	ict, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		log.Printf("error loading location 'Asia/Bangkok': %v\n", err)
	}
	time.Local = ict
}
