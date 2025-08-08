package http

import (
	"golang-clean-architecture/internal/delivery/http/middleware"

	"github.com/gofiber/fiber/v2"
)

type HelloController struct {
}

func NewHelloController() *HelloController {
	return &HelloController{}
}

func (h *HelloController) SayHello(ctx *fiber.Ctx) error {
	auth := middleware.GetUser(ctx)
	return ctx.Send([]byte("Hello " + auth.ID))
}
