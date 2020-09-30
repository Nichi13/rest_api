package store

import "gismart-rest-api/internal/app/model"


type MenuRepository interface {
	Create(*model.Menu) error
}

type OrderRepository interface {
	Create([]string, []string) (int, error)
	Update(*model.Order) error
	Get(string)([]model.Order, error)
}