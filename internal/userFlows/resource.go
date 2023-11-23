package userFlows

// Resource represents a user flow.
type Resource struct {
	ID          int `json:"id"`
	UserID      int `json:"user_id"`
	StepID      int `json:"step_id"`
	ProductID   int `json:"product_id"`
	OrderItemID int `json:"order_item_id"`

	Complete bool `json:"complete"`

	Step    *Step
	Product *Product
}

// Step represents a step.
type Step struct {
	ID      int      `json:"id"`
	Name    string   `json:"name"`
	Order   string   `json:"description"`
	Details string   `json:"details"`
	Options []string `json:"options"`
}

//// StringArray represents a string array.
//type StringArray []string
//
//// Scan implements the sql.Scanner interface.
//func (sa *StringArray) Scan(value interface{}) error {
//	if value == nil {
//		*sa = StringArray{}
//		return nil
//	}
//	return json.Unmarshal(value.([]byte), sa)
//}

// Product represents a product.
type Product struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
