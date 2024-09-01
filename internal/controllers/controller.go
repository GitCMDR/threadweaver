package controllers

import (
"gopkg.in/telebot.v3"
)

type Controller struct {
	// lets add some gateways here later
}

func NewController() *Controller {
	return &Controller{}
}

func (c *Controller) ProcessMessage (tCtx telebot.Context) error {
	return tCtx.Send("You said " + tCtx.Text())
}

func (c *Controller) StartCommand (tCtx telebot.Context) error {
	return tCtx.Send("Ready to post.")
}

func (c *Controller) HelpCommand (tCtx telebot.Context) error {
	return tCtx.Send("You are on your own, kiddo *smirks*")
}