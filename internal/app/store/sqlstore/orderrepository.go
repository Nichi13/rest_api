package sqlstore

import (
	"fmt"
	"gismart-rest-api/internal/app/model"
)

type OrderRepository struct {
	store *Store
}

func (order *OrderRepository) Create(d []string, c []string) (int, error) {
	var (
		order_id int
		order_number int
		err error = nil
	)
	if err :=order.store.db.QueryRow(
		"INSERT INTO orders (status) VALUES ('new') RETURNING id, number",
	).Scan(&order_id, &order_number); err != nil{
		return order_number, err
	}
	for index, value := range d{
		var (
			menu_number int
			order_line int
		)
		if err := order.store.db.QueryRow(
			"SELECT id FROM menu_positions WHERE number = $1",
			value,
		).Scan(&menu_number); err != nil {
			return order_number, err
		}
		if err :=order.store.db.QueryRow(
			"INSERT INTO orders_line (order_id, menu_position_id, count) VALUES ($1, $2, $3)  RETURNING id",
			order_id,
			menu_number,
			c[index],
		).Scan(&order_line); err != nil{
			return order_number, err
		}

	}
	return order_number, err
}

func (order *OrderRepository) Update(m *model.Order) error {
	return order.store.db.QueryRow(
		"UPDATE orders SET status = $1 WHERE number = $2 RETURNING id",
		m.Status,
		m.Number,
	).Scan(&m.ID)
}

func (o *OrderRepository) Get(status string) ([]model.Order, error) {
	data, err := o.store.db.Query(
		"SELECT * FROM orders WHERE status = $1",
		status,
	)
	if err != nil {
		return nil, err
	}
	orders := []model.Order{}
	for data.Next(){
		o := model.Order{}
		err := data.Scan(&o.ID,  &o.Status, &o.Number)
		if err != nil {
			fmt.Println(err)
			continue
		}
		orders = append(orders, o)
	}
	return orders, nil
}