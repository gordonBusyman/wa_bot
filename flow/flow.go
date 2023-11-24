package flow

import (
	"context"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/gordonBusyman/wa_bot/internal/userFlows"
	"github.com/gordonBusyman/wa_bot/internal/users"
)

const (
	// MsgTypeMessage is a message type simple message.
	MsgTypeMessage = "message"
	// MsgTypeCQ is a message type callback query.
	MsgTypeCQ = "cq"
)

// Driver is a driver for a flow.
type Driver struct {
	Bot *tgbotapi.BotAPI

	UserFlowsStore *userFlows.Store
}

// StartConversation starts a new conversation with a user.
func (d *Driver) StartConversation(user *users.Resource) {
	step, err := d.UserFlowsStore.RetrieveCurrent(context.Background(), user.ID)
	if err != nil {
		log.Printf("error retrieving first step for user %d: %v", user.ID, err)
	}

	d.SendStepMessage(user.ChatID, step)
}

// HandleResponse handles a response from a user.
func (d *Driver) HandleResponse(user *users.Resource, response string, _ string) {
	currentStep, err := d.UserFlowsStore.RetrieveCurrent(context.Background(), user.ID)
	if err != nil {
		log.Printf("error retrieving flow step for user %d: %v", user.ID, err)
	}
	// Update the user's response
	d.UserFlowsStore.Update(context.Background(), currentStep, response)

	nextStep, err := d.UserFlowsStore.RetrieveCurrent(context.Background(), user.ID)
	if err != nil {
		log.Printf("error retrieving next step for user %d: %v", user.ID, err)
	}

	if nextStep != nil {
		d.SendStepMessage(user.ChatID, nextStep)
	} else {
		d.EndConversation(user.ChatID)
		d.UserFlowsStore.Update(context.Background(), currentStep, response)
	}
}

// SendStepMessage sends a message to a user.
func (d *Driver) SendStepMessage(chatID int64, step *userFlows.Resource) {
	var msg tgbotapi.MessageConfig

	if step.Step.Options != nil {
		msg = tgbotapi.NewMessage(chatID, step.Step.Details+": "+step.Product.Name)

		var keyboard tgbotapi.InlineKeyboardMarkup
		options := step.Step.Options
		row := make([]tgbotapi.InlineKeyboardButton, len(options))
		for i := 0; i < len(options); i++ {
			row[i] = tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf("%v", options[i]),
				fmt.Sprintf("%v", options[i]),
			)
		}

		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
		msg.ReplyMarkup = keyboard
	} else {
		msg = tgbotapi.NewMessage(chatID, step.Step.Details)
	}

	if _, err := d.Bot.Send(msg); err != nil {
		log.Printf("error sending message to user %d: %v", chatID, err)
	}
}

// EndConversation ends a conversation with a user.
func (d *Driver) EndConversation(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Thank you for your responses!")
	d.Bot.Send(msg)
}
