package store

import "gismart-rest-api/internal/app/model"

type UserRepository interface {
	Create(*model.User) error
	Find(int) (*model.User, error)
	FindByEmail(string)(*model.User, error)
}

type MenuRepository interface {
	Create(*model.Menu) error
}

type OrderRepository interface {
	Create([]string, []string) (int, error)
	Update(*model.Order) error
	Get(string)([]model.Order, error)
}