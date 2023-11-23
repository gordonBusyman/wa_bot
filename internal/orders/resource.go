package orders

// Resource represents a user flow.
type Resource struct {
	ID     int `json:"id"`
	UserID int `json:"user_id"`

	Items []Item `json:"items"`
}

type Item struct {
	ID        int `json:"id"`
	OrderID   int `json:"order_id"`
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}
