package model

// Модель для позиции заказа
type Order_line struct {
	ID int `json:"id"`
	OrderId int `json:"order_id"`
	MenuPosition int `json:"menu_position"`
	Count string `json:"count"`
}
