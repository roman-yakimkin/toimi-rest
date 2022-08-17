package models

import (
	"errors"
	"fmt"
	"time"
	errors2 "toimi/internal/app/errors"
	"toimi/internal/app/services/configmanager"
	"unicode/utf8"
)

type Advert struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Created     time.Time `json:"created"`
	Price       uint      `json:"price"`
	Photos      []string  `json:"photos"`
}

func (a *Advert) BeforeSave() {
	if a.Created.IsZero() {
		a.Created = time.Now()
	}
}

func (a *Advert) Validate(cfg *configmanager.Config) error {
	var err error
	if utf8.RuneCountInString(a.Title) > cfg.Validate.TitleMax {
		err = errors.New(fmt.Sprintf("Title must not be more that %d symbols", cfg.Validate.TitleMax))
	}
	if utf8.RuneCountInString(a.Description) > cfg.Validate.DescriptionMax {
		err = errors.New(fmt.Sprintf("Description must not be more that %d symbols", cfg.Validate.DescriptionMax))
	}
	if len(a.Photos) > cfg.Validate.PhotosMax {
		err = errors.New(fmt.Sprintf("Amount of photos must not be more that %d", cfg.Validate.PhotosMax))
	}
	if err != nil {
		return &errors2.ValidationError{Err: err}
	}
	return nil
}
