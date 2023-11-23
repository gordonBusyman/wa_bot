package flowSteps

import "encoding/json"

// Resource represents a flow step.
type Resource struct {
	ID      int         `json:"id"`
	Name    string      `json:"name"`
	Order   string      `json:"description"`
	Details string      `json:"details"`
	Options StringArray `json:"options"`
}

// StringArray represents a string array.
type StringArray []string

// Scan implements the sql.Scanner interface.
func (sa *StringArray) Scan(value interface{}) error {
	if value == nil {
		*sa = StringArray{}
		return nil
	}
	return json.Unmarshal(value.([]byte), sa)
}
