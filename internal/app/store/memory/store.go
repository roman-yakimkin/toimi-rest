package memory

import (
	"toimi/internal/app/interfaces"
)

type Store struct {
	advertRepo interfaces.AdvertRepo
}

func NewStore(advertRepo interfaces.AdvertRepo) interfaces.Store {
	return &Store{
		advertRepo: advertRepo,
	}
}

func (s *Store) Advert() interfaces.AdvertRepo {
	return s.advertRepo
}
