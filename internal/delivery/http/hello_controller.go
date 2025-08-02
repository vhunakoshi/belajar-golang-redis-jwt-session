package http

import (
	"github.com/gofiber/fiber/v2"
	"golang-clean-architecture/internal/delivery/http/middleware"
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
