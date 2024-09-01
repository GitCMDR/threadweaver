package handlers

import (
	"strings"
	"os"

	"github.com/GitCMDR/threadweaver/internal/controllers"
	"gopkg.in/telebot.v3"
)

type Handler struct {
	Controller *controllers.Controller
	Codeword string
}

func NewHandler(controller *controllers.Controller) *Handler {
	return &Handler{
		Controller: controller,
		Codeword: os.Getenv("CODEWORD"),
	}
}

func (h *Handler) HandleMessage(tCtx telebot.Context) error {
	message := tCtx.Text()

	// check if the message starts with the codeword
	if strings.HasPrefix(strings.ToLower(message), strings.ToLower(h.Codeword)) {
		// remove the codeword from the message
		msg := strings.TrimSpace(strings.TrimPrefix(message, h.Codeword))

		// call the controller to process the message without the codeword
		err := h.Controller.ProcessMessageWithText(tCtx, msg)
		if err != nil {
			return tCtx.Send("Failed to process message: " + err.Error())
		} 
	} else {
		// if the codeword is not present, return a response with attitude
		return tCtx.Send("I don't even know you.")
	}
		
	return nil
}

func (h *Handler) HandleStart(tCtx telebot.Context) error {
	return h.Controller.StartCommand(tCtx)
}

func (h *Handler) HandleHelp(tCtx telebot.Context) error {
	return h.Controller.HelpCommand(tCtx)
}
