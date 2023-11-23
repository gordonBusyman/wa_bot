package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-chi/chi"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gordonBusyman/wa_bot/api"
	"github.com/gordonBusyman/wa_bot/flow"
	"github.com/gordonBusyman/wa_bot/internal/userFlows"
	"github.com/gordonBusyman/wa_bot/internal/users"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"time"
)

var (
	db  *sql.DB
	bot *tgbotapi.BotAPI
	u   tgbotapi.UpdateConfig
)

func init() {
	var err error

	bot, err = tgbotapi.NewBotAPI("X")
	if err != nil {
		log.Panic(err)
	}

	// Set up an update configuration
	u = tgbotapi.NewUpdate(0)
	u.Timeout = 60

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		"localhost", 5432, "postgres", "postgres", "postgres")

	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	flowDriver := &flow.Driver{
		Bot:            bot,
		UserFlowsStore: userFlows.NewStore(db),
	}
	go startTelegramBot(flowDriver)
	go startHTTPServer(flowDriver)

	select {}

	defer func() {
		db.Close()
	}()

}

func startHTTPServer(flowDriver *flow.Driver) {
	r := chi.NewRouter()
	fmt.Println("http server listening on 8080")

	m := api.Mux{
		DB:     db,
		Bot:    bot,
		Driver: flowDriver,
	}
	// Define a route with a URL parameter
	r.Post("/orders/{user_id}", m.CreateOrder)

	// Start the HTTP server on port 8080
	http.ListenAndServe(":8080", r)
}

func startTelegramBot(flowDriver *flow.Driver) {
	// Retrieve updates from the bot
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("telegram bot is up and ready")

	// Process each update received
	for update := range updates {
		fmt.Println("update: ", update)
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

		//if update.CallbackQuery != nil {
		//	handleCallbackQuery(bot, update.CallbackQuery)
		//	continue
		//}

		if update.Message == nil && update.CallbackQuery == nil {
			continue
		}

		var userID int
		if update.Message != nil {
			userID = update.Message.From.ID
		} else if update.CallbackQuery != nil {
			userID = update.CallbackQuery.From.ID
		}

		uStore := users.NewStore(db)
		user, err := uStore.Retrieve(ctx, userID)
		if err != nil {
			log.Printf("error retrieving user: %v", err)

			continue
		}

		if update.Message != nil && update.Message.Text == "/register" {
			if user != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "You are already registered, "+update.Message.From.FirstName+".")

				bot.Send(msg)
			} else {
				uStore.Create(ctx, update.Message)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Thanks, "+update.Message.From.FirstName+", for registering!")

				bot.Send(msg)
			}

			continue
		}

		// Unregistered users are not allowed to proceed
		if user == nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hello, "+update.Message.From.FirstName+"! Please type /register to complete your registration.")
			msg.ParseMode = "Markdown"
			bot.Send(msg)

			continue
		}

		// Get the current step for the user
		userFlowStep, err := userFlows.NewStore(db).RetrieveCurrent(ctx, userID)
		if err != nil {
			log.Printf("error retrieving user flow: %v", err)

			continue
		}

		if userFlowStep != nil {
			if update.Message != nil {
				flowDriver.HandleResponse(user, update.Message.Text, flow.MsgTypeMessage)
			} else if update.CallbackQuery != nil {
				flowDriver.HandleResponse(user, update.CallbackQuery.Data, flow.MsgTypeCQ)
			}

			continue
		}

		// Reply to the user when there is no further action required.
		reply := tgbotapi.NewMessage(update.Message.Chat.ID, "No further action is required from you at this moment. Thank you for your time!")
		bot.Send(reply)
	}
}

//func newFlow(userID int) any {
//	return userID
//}
//func scoreKeyboard() tgbotapi.InlineKeyboardMarkup {
//	var keyboard tgbotapi.InlineKeyboardMarkup
//	row := make([]tgbotapi.InlineKeyboardButton, 5)
//	for i := 1; i <= 5; i++ {
//		row[i-1] = tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%d", i), fmt.Sprintf("score_%d", i))
//	}
//	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)
//	return keyboard
//}

func handleCallbackQuery(bot *tgbotapi.BotAPI, query *tgbotapi.CallbackQuery) {
	callbackData := query.Data
	chatID := query.Message.Chat.ID
	messageID := query.Message.MessageID

	// Respond to the callback query, acknowledging it
	callbackConfig := tgbotapi.NewCallback(query.ID, "")
	bot.AnswerCallbackQuery(callbackConfig)

	// Handle different scores
	var responseText string
	switch callbackData {
	case "score_1", "score_2", "score_3", "score_4", "score_5":
		responseText = "Thanks for rating! Could you please share a brief opinion or description of your experience with our product? Your feedback is valuable to us!"

		// responseText = fmt.Sprintf(": %s", callbackData[len("score_"):])
	default:
		responseText = "Unknown option selected"
	}

	editMsg := tgbotapi.NewEditMessageText(chatID, messageID, responseText)
	bot.Send(editMsg)
}

// docker run --name connectly -e POSTGRES_PASSWORD=postgres -d postgres
// https://reintech.io/blog/creating-a-sentiment-analysis-tool-with-go
