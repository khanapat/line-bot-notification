package middleware

import (
	"line-notification/common"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type middleware struct {
	ZapLogger *zap.Logger
}

func NewMiddleware(zapLogger *zap.Logger) *middleware {
	return &middleware{
		ZapLogger: zapLogger,
	}
}

func (m *middleware) JSONMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Accepts(common.ApplicationJSON)

		return c.Next()
	}
}
