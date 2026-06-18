package repositories

type Storage struct {
	ReviewRepository *ReviewRepository
}

func NewStorage() *Storage {
	newStorage := &Storage{
		ReviewRepository: &ReviewRepository{},
	}

	return newStorage
}
