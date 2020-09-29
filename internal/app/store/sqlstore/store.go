package sqlstore

import (
	"database/sql"
	"gismart-rest-api/internal/app/store"

	_ "github.com/lib/pq"
)

type Store struct {
	db *sql.DB
	userRepository *UserRepository
	menuRepository *MenuRepository
	orderRepository *OrderRepository
}

func New(db *sql.DB) *Store{
	return &Store{
		db: db,
	}
}

func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
	}
	return s.userRepository
}

func (s *Store) Menu() store.MenuRepository {
	if s.menuRepository != nil {
		return s.menuRepository
	}

	s.menuRepository = &MenuRepository{
		store: s,
	}
	return s.menuRepository
}

func (s *Store) Order() store.OrderRepository {
	if s.orderRepository != nil {
		return s.orderRepository
	}

	s.orderRepository = &OrderRepository{
		store: s,
	}
	return s.orderRepository
}