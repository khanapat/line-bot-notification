package handler

import (
	"line-notification/common"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type Ctx struct {
	*fiber.Ctx
	logger *zap.Logger
}

func (c *Ctx) Log() *zap.Logger {
	return c.logger
}

func (c *Ctx) Status(status int) *Ctx {
	c.Response().Header.Add(common.XRequestID, string(c.Request().Header.Peek(common.XRequestID)))
	ctx := c.Ctx.Status(status)
	return &Ctx{ctx, c.logger}
}

type Handler func(*Ctx) error

func Helper(handler Handler, logger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		xRequestID := c.Request().Header.Peek(common.XRequestID)
		return handler(&Ctx{
			c,
			logger.With(zap.String(common.XRequestID, string(xRequestID))),
		})
	}
}
