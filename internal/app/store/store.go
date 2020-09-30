package store

type Store interface {
	Menu() MenuRepository
	Order() OrderRepository
}
