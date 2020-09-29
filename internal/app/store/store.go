package store

type Store interface {
	User() UserRepository
	Menu() MenuRepository
	Order() OrderRepository
}
