package model

// Модель для заказа
type Order struct {
	ID int `json:"id"`
	Number int `json:"number"`
	Status string `json:"status"`
}
