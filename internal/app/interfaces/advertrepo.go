package interfaces

import "toimi/internal/app/models"

const (
	SortNone = iota
	SortByCreatedAsc
	SortByCreatedDesc
	SortByPriceAsc
	SortByPriceDesc
)

type AdvertRepo interface {
	Save(a *models.Advert) (string, error)
	Delete(string) error
	GetByID(id string) (*models.Advert, error)
	GetSortedPage(int, int) ([]models.AdvertShort, error)
}
