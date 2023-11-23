package flow

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gordonBusyman/wa_bot/internal/userFlows"
	"github.com/gordonBusyman/wa_bot/internal/users"
	"log"
)

const (
	MsgTypeMessage = "message"
	MsgTypeCQ      = "cq"
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

	fmt.Println("HandleResponse NEXT: ", nextStep)

	if nextStep != nil {
		d.SendStepMessage(user.ChatID, nextStep)
	} else {
		d.EndConversation(user.ChatID)
		d.UserFlowsStore.Update(context.Background(), currentStep, response)
	}
}

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
				fmt.Sprintf("score_%v", options[i]),
			)
		}

		keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
		msg.ReplyMarkup = keyboard
	} else {
		msg = tgbotapi.NewMessage(chatID, step.Step.Details)
	}

	_, err := d.Bot.Send(msg)
	if err != nil {
		log.Printf("error sending message to user %d: %v", chatID, err)
	}
}

func (d *Driver) EndConversation(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Thank you for your responses!")
	d.Bot.Send(msg)
}
