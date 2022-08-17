package memory

import (
	"sort"
	"strconv"
	"toimi/internal/app/errors"
	"toimi/internal/app/interfaces"
	"toimi/internal/app/models"
	"toimi/internal/app/services/configmanager"
)

type AdvertRepo struct {
	adverts map[string]models.Advert
	lastID  int
	cfg     *configmanager.Config
}

func NewAdvertRepo(cfg *configmanager.Config) interfaces.AdvertRepo {
	return &AdvertRepo{
		adverts: make(map[string]models.Advert),
		cfg:     cfg,
	}
}

func (r *AdvertRepo) Save(a *models.Advert) (string, error) {
	if err := a.Validate(r.cfg); err != nil {
		return "", err
	}
	a.BeforeSave()
	if a.ID == "" {
		r.lastID++
		a.ID = strconv.Itoa(r.lastID)
	}
	r.adverts[a.ID] = *a
	return a.ID, nil
}

func (r *AdvertRepo) Delete(id string) error {
	_, ok := r.adverts[id]
	if !ok {
		return errors.ErrAdvertNotFound
	}
	delete(r.adverts, id)
	return nil
}

func (r *AdvertRepo) GetByID(id string) (*models.Advert, error) {
	advert, ok := r.adverts[id]
	if !ok {
		return nil, errors.ErrAdvertNotFound
	}
	return &advert, nil
}

func (r *AdvertRepo) GetSortedPage(page int, sortMode int) ([]models.AdvertShort, error) {
	sorted, err := r.getAllSorted(sortMode)
	if err != nil {
		return nil, err
	}
	startPos := r.cfg.Paginate.AdvertsPageSize * page
	endPos := startPos + r.cfg.Paginate.AdvertsPageSize

	if startPos > len(sorted)-1 {
		return nil, nil
	}
	if endPos > len(sorted) {
		endPos = len(sorted)
	}

	paged := sorted[startPos:endPos]
	result := make([]models.AdvertShort, 0, len(paged))
	for _, advert := range paged {
		var advertShort models.AdvertShort
		advertShort.GetFromFull(&advert)
		result = append(result, advertShort)
	}

	return result, nil
}

func (r *AdvertRepo) getAll() ([]models.Advert, error) {
	result := make([]models.Advert, 0, len(r.adverts))
	for _, task := range r.adverts {
		result = append(result, task)
	}
	return result, nil
}

func (r *AdvertRepo) getAllSorted(sortMode int) ([]models.Advert, error) {
	adverts, err := r.getAll()
	if err != nil {
		return nil, err
	}
	switch sortMode {
	case interfaces.SortByCreatedAsc:
		sort.Slice(adverts, func(i, j int) bool {
			return adverts[i].Created.Before(adverts[j].Created)
		})
	case interfaces.SortByCreatedDesc:
		sort.Slice(adverts, func(i, j int) bool {
			return adverts[i].Created.After(adverts[j].Created)
		})
	case interfaces.SortByPriceAsc:
		sort.Slice(adverts, func(i, j int) bool {
			return adverts[i].Price < adverts[j].Price
		})
	case interfaces.SortByPriceDesc:
		sort.Slice(adverts, func(i, j int) bool {
			return adverts[i].Price > adverts[j].Price
		})
	}
	return adverts, nil
}
