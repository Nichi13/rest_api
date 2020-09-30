package model

// Модель для позиции меню
type Menu struct {
	ID int `json:"id"`
	Number string `json:"number"`
	Name string `json:"name"`
}
