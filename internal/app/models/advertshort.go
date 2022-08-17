package models

type AdvertShort struct {
	Title string `json:"title"`
	Photo string `json:"photo,omitempty"`
	Price uint   `json:"price"`
}

func (as *AdvertShort) GetFromFull(a *Advert) {
	as.Title = a.Title
	if len(a.Photos) > 0 {
		as.Photo = a.Photos[0]
	}
	as.Price = a.Price
}
