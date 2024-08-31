package main

import (
    "log"
    "time"

    "gopkg.in/telebot.v3"
)

func main() {
    // Replace with your Bot API token
    botToken := "TOKEN_GOES_HERE"

    bot, err := telebot.NewBot(telebot.Settings{
        Token:  botToken,
        Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
    })
    if err != nil {
        log.Fatal(err)
    }

    // Handle /start command
    bot.Handle("/start", func(c telebot.Context) error {
        return c.Send("Hello! I’m your new bot.")
    })

    // Handle any text message
    bot.Handle(telebot.OnText, func(c telebot.Context) error {
        return c.Send("You said: " + c.Text())
    })

    // Start the bot
    bot.Start()
}
