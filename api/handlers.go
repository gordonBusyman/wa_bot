package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gordonBusyman/wa_bot/flow"
	"github.com/gordonBusyman/wa_bot/internal/orders"
	"github.com/gordonBusyman/wa_bot/internal/userFlows"
	"github.com/gordonBusyman/wa_bot/internal/users"
	"net/http"
	"strconv"
)

// Mux represents the API.
type Mux struct {
	DB  *sql.DB
	Bot *tgbotapi.BotAPI

	Driver *flow.Driver
}

// CreateOrder handles the POST /orders/{user_id} endpoint.
func (m Mux) CreateOrder(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(chi.URLParam(r, "user_id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	var items []orders.Item

	// Decode the JSON body
	if err = json.NewDecoder(r.Body).Decode(&items); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	// create a new order
	order, err := orders.NewStore(m.DB).Create(r.Context(), userID, items)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
	fmt.Println("1")

	// create a flow for user
	if err := userFlows.NewStore(m.DB).CreateMany(r.Context(), userID, 1, order.ID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
	fmt.Println("2")
	user, err := users.NewStore(m.DB).Retrieve(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
	fmt.Println("3")

	m.Driver.StartConversation(user)

	w.Write([]byte(fmt.Sprintf("Order %d created", order.ID)))
}
