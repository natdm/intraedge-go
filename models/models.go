package models

// Coffee ...
type Coffee struct {
	Name   string `json:"name"`
	Sugars int    `json:"sugars"`
	Cream  bool   `json:"cream"`
}

type Storer interface {
	State() map[string]interface{}
	Val(key string) (interface{}, error)
	Add(key string, val interface{}) error
}
