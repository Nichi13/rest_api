package sqlstore

import (
	"gismart-rest-api/internal/app/model"
)

type MenuRepository struct {
	store *Store
}

func (mr *MenuRepository) Create(m *model.Menu) error {
	return mr.store.db.QueryRow(
		"INSERT INTO menu_positions (number, name) VALUES ($1, $2) RETURNING id",
		m.Number,
		m.Name,
	).Scan(&m.ID)
}
