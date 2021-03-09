package main

import (
	"fmt"
	"line-notification/bot"
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

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func init() {
	runtime.GOMAXPROCS(1)
	initViper()
	initTimezone()
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

	logger := logz.NewLogConfig()

	middle := middleware.NewMiddleware(logger)

	line := app.Group(viper.GetString("app.context"))

	// api.Use(recover.New())
	// api.Use(requestid.New())
	// api.Use(logger.New())
	line.Use(middle.JSONMiddleware())

	line.Post("/webhook", handler.Helper(bot.NewBotHandler().ReplyBot, logger))

	line.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	log.Println(fmt.Sprintf("⇨ http server started on [::]:%s", viper.GetString("app.port")))

	go func() {
		if err := app.Listen(fmt.Sprintf(":%s", viper.GetString("app.port"))); err != nil {
			log.Panic(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	select {
	case <-c:
		log.Println("terminating: by signal")
	}

	app.Shutdown()

	log.Println("shutting down")
	os.Exit(0)
}

func initViper() {
	viper.SetDefault("app.name", "line-notification-api")
	viper.SetDefault("app.port", "9090")
	viper.SetDefault("app.timeout", "60s")
	viper.SetDefault("app.context", "/line")

	viper.SetDefault("log.level", "debug")
	viper.SetDefault("log.env", "dev")

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
