package orders

// Resource represents an order.
type Resource struct {
	ID     int `json:"id"`
	UserID int `json:"user_id"`

	Items []Item `json:"items"`
}

// Item represents an order item.
type Item struct {
	ID        int `json:"id"`
	OrderID   int `json:"order_id"`
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}
