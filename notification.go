package main

import (
	"github.com/gregdel/pushover"
)

func sendNotification(msg, title string, config *Config) error {
	if config.PushoverToken == "" || config.PushoverUserKey == "" {
		return nil
	}
	// Create a new pushover app with a token
	app := pushover.New(config.PushoverToken)

	// Create a new recipient
	recipient := pushover.NewRecipient(config.PushoverUserKey)

	// Create the message to send
	message := pushover.NewMessageWithTitle(msg, title)

	// Send the message to the recipient
	_, err := app.SendMessage(message, recipient)
	if err != nil {
		return err
	}
	return nil
}
