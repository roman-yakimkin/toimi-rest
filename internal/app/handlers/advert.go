package handlers

import (
	"encoding/json"
	errors2 "errors"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strconv"
	"strings"
	"toimi/internal/app/errors"
	"toimi/internal/app/interfaces"
	"toimi/internal/app/models"
)

type AdvertController struct {
	store interfaces.Store
}

func NewAdvertController(store interfaces.Store) *AdvertController {
	return &AdvertController{
		store: store,
	}
}

func (c *AdvertController) getPageFromRequest(r *http.Request) int {
	query := r.URL.Query()
	page, err := strconv.Atoi(query.Get("page"))
	if err != nil {
		page = 0
	}
	return page
}

func (c *AdvertController) getSortFromRequest(r *http.Request) int {
	query := r.URL.Query()
	sort := strings.ToLower(query.Get("sort"))
	switch sort {
	case "created-asc":
		return interfaces.SortByCreatedAsc
	case "created-desc":
		return interfaces.SortByCreatedDesc
	case "price-asc":
		return interfaces.SortByPriceAsc
	case "price-desc":
		return interfaces.SortByPriceDesc
	default:
		return interfaces.SortNone
	}
}

func (c *AdvertController) Create(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if returnErrorResponse(err != nil, w, r, http.StatusInternalServerError, err, "") {
		return
	}
	var advert models.Advert
	err = json.Unmarshal(body, &advert)
	if returnErrorResponse(err != nil, w, r, http.StatusInternalServerError, err, "") {
		return
	}

	advertId, err := c.store.Advert().Save(&advert)
	var ve *errors.ValidationError
	if returnErrorResponse(errors2.As(err, &ve), w, r, http.StatusBadRequest, err, "") {
		return
	}
	if returnErrorResponse(err != nil, w, r, http.StatusInternalServerError, err, "") {
		return
	}
	returnSuccessResponse(w, r, "A new advert has been created", struct {
		Id string `json:"id"`
	}{Id: advertId})
}

func (c *AdvertController) Update(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if returnErrorResponse(err != nil, w, r, http.StatusInternalServerError, err, "") {
		return
	}
	advert := new(models.Advert)
	err = json.Unmarshal(body, advert)
	if returnErrorResponse(err != nil, w, r, http.StatusInternalServerError, err, "") {
		return
	}
	advert, err = c.store.Advert().GetByID(advert.ID)
	if returnErrorResponse(err != nil, w, r, http.StatusNotFound, err, "") {
		return
	}
	_, err = c.store.Advert().Save(advert)
	var ve *errors.ValidationError
	if returnErrorResponse(errors2.As(err, &ve), w, r, http.StatusBadRequest, err, "") {
		return
	}
	if returnErrorResponse(err != nil, w, r, http.StatusInternalServerError, err, "") {
		return
	}
	returnSuccessResponse(w, r, "advert has been updated", struct {
		Id string `json:"id"`
	}{Id: advert.ID})
}

func (c *AdvertController) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	advertId, ok := vars["id"]
	if returnErrorResponse(!ok, w, r, http.StatusNotFound, errors.ErrAdvertNotFound, "") {
		return
	}
	_, err := c.store.Advert().GetByID(advertId)
	if returnErrorResponse(err != nil, w, r, http.StatusNoContent, err, "") {
		return
	}
	err = c.store.Advert().Delete(advertId)
	if returnErrorResponse(err != nil, w, r, http.StatusNoContent, err, "") {
		return
	}
	returnSuccessResponse(w, r, "advert has been deleted", struct {
		Id string `json:"id"`
	}{Id: advertId})
}

func (c *AdvertController) GetOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	advertId, ok := vars["id"]
	if returnErrorResponse(!ok, w, r, http.StatusNotFound, errors.ErrAdvertNotFound, "") {
		return
	}
	advert, err := c.store.Advert().GetByID(advertId)
	if returnErrorResponse(err != nil, w, r, http.StatusNotFound, err, "") {
		return
	}
	returnSuccessResponse(w, r, "task has been updated", advert)
}

func (c *AdvertController) GetPage(w http.ResponseWriter, r *http.Request) {
	page, sort := c.getPageFromRequest(r), c.getSortFromRequest(r)
	adverts, err := c.store.Advert().GetSortedPage(page, sort)
	if returnErrorResponse(err != nil, w, r, http.StatusInternalServerError, err, "") {
		return
	}
	returnSuccessResponse(w, r, "", adverts)
}
