package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/cdipaolo/sentiment"
	"github.com/go-chi/chi"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/lib/pq"

	"github.com/gordonBusyman/wa_bot/api"
	"github.com/gordonBusyman/wa_bot/flow"
	"github.com/gordonBusyman/wa_bot/internal/userFlows"
	"github.com/gordonBusyman/wa_bot/internal/users"
)

var (
	db  *sql.DB
	bot *tgbotapi.BotAPI
	u   tgbotapi.UpdateConfig
)

func init() {
	var err error

	// Retrieve the value of an environment variable
	token := os.Getenv("CONNECTLY_BOT_TOKEN")
	if token == "" {
		panic("CONNECTLY_BOT_TOKEN is not set")
	}

	bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	// Set up an update configuration
	u = tgbotapi.NewUpdate(0)
	u.Timeout = 60

	//export CONNECTLY_DB_HOST=localhost
	//export CONNECTLY_DB_USERNAME=postgres
	//export CONNECTLY_DB_PASSWORD=postgres
	dbHost := os.Getenv("CONNECTLY_DB_HOST")
	if token == "" {
		panic("CONNECTLY_DB_HOST is not set")
	}
	dbUser := os.Getenv("CONNECTLY_DB_USERNAME")
	if token == "" {
		panic("CONNECTLY_DB_USERNAME is not set")
	}
	dbPass := os.Getenv("CONNECTLY_DB_PASSWORD")
	if token == "" {
		panic("CONNECTLY_DB_PASSWORD is not set")
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbHost, 5432, dbUser, dbPass, "postgres")

	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	model, err := sentiment.Restore()
	if err != nil {
		panic(err)
	}

	flowDriver := &flow.Driver{
		Bot:            bot,
		UserFlowsStore: userFlows.NewStore(db, &model),
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
		fmt.Printf("got update message: %+v", update)
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

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
				u := &users.Resource{
					ID:     update.Message.From.ID,
					ChatID: update.Message.Chat.ID,
					Name:   update.Message.From.FirstName,
				}

				uStore.Create(ctx, u)
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
		userFlowStep, err := userFlows.NewStore(db, nil).RetrieveCurrent(ctx, userID)
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
