package handlers

import (
	"github.com/GitCMDR/microblogreposter-bot/internal/controllers"
	"gopkg.in/telebot.v3"
)

type Handler struct {
	Controller *controllers.Controller
}

func NewHandler(controller *controllers.Controller) *Handler {
	return &Handler{Controller: controller}
}

func (h *Handler) HandleMessage(tCtx telebot.Context) error {
	return h.Controller.ProcessMessage(tCtx)
}

func (h *Handler) HandleStart(tCtx telebot.Context) error {
	return h.Controller.StartCommand(tCtx)
}

func (h *Handler) HandleHelp(tCtx telebot.Context) error {
	return h.Controller.HelpCommand(tCtx)
}
