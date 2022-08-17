package postgres

import (
	"context"
	"github.com/jackc/pgx/v4"
	"strconv"
	"toimi/internal/app/interfaces"
	"toimi/internal/app/models"
	"toimi/internal/app/services/configmanager"
	"toimi/internal/app/services/dbclient"
)

type AdvertRepo struct {
	cfg *configmanager.Config
	db  *dbclient.PostgresDBClient
}

func NewAdvertRepo(cfg *configmanager.Config, db *dbclient.PostgresDBClient) interfaces.AdvertRepo {
	return &AdvertRepo{
		cfg: cfg,
		db:  db,
	}
}

func (r *AdvertRepo) Save(a *models.Advert) (string, error) {
	if err := a.Validate(r.cfg); err != nil {
		return "", err
	}
	a.BeforeSave()
	ctx := context.Background()
	conn, err := r.db.Connect(ctx)
	defer r.db.Disconnect(ctx)
	if err != nil {
		return "", err
	}
	tx, err := conn.Begin(ctx)
	id := 0
	if a.ID == "" {
		err = tx.QueryRow(ctx, "INSERT INTO adverts (title, description, created, price) values($1, $2, $3, $4) RETURNING id",
			a.Title, a.Description, a.Created, a.Price).Scan(&id)
		if err != nil {
			tx.Rollback(ctx)
			return "", err
		}
		if err != nil {
			return "", err
		}
		if err = r.insertPhotos(ctx, tx, id, a.Photos); err != nil {
			return "", err
		}
	} else {
		id, err = strconv.Atoi(a.ID)
		if err != nil {
			tx.Rollback(ctx)
			return "", err
		}
		_, err = tx.Exec(ctx, "UPDATE adverts SET title=$1, description=$2, created=$3, price=$4 WHERE id=$5",
			a.Title, a.Description, a.Created, a.Price, id)
		if err != nil {
			tx.Rollback(ctx)
			return "", err
		}
		if err = r.deletePhotos(ctx, tx, id); err != nil {
			return "", err
		}
		if err = r.insertPhotos(ctx, tx, id, a.Photos); err != nil {
			return "", err
		}
	}
	err = tx.Commit(ctx)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(id), nil
}

func (r *AdvertRepo) insertPhotos(ctx context.Context, tx pgx.Tx, id int, photos []string) error {
	for i, photo := range photos {
		_, err := tx.Exec(ctx, "INSERT INTO adverts_photos(advert_id, photo, delta) VALUES ($1, $2, $3)",
			id, photo, i)
		if err != nil {
			tx.Rollback(ctx)
			return err
		}
	}
	return nil
}

func (r *AdvertRepo) deletePhotos(ctx context.Context, tx pgx.Tx, id int) error {
	_, err := tx.Exec(ctx, "delete from adverts_photos where advert_id = $1", id)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}
	return nil
}

func (r *AdvertRepo) Delete(id string) error {
	ctx := context.Background()
	conn, err := r.db.Connect(ctx)
	defer r.db.Disconnect(ctx)
	if err != nil {
		return err
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	if err = r.deletePhotos(ctx, tx, idInt); err != nil {
		return err
	}
	_, err = tx.Exec(ctx, "delete from adverts where id = $1", idInt)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}
	err = tx.Commit(ctx)
	return err
}

func (r *AdvertRepo) GetByID(id string) (*models.Advert, error) {
	ctx := context.Background()
	conn, err := r.db.Connect(ctx)
	defer r.db.Disconnect(ctx)
	if err != nil {
		return nil, err
	}
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	var a models.Advert
	err = conn.QueryRow(ctx,
		"select cast(id as varchar), title, description, created, price 	from adverts where id = $1", idInt).
		Scan(&a.ID, &a.Title, &a.Description, &a.Created, &a.Price)
	if err != nil {
		return nil, err
	}
	rows, err := conn.Query(ctx, "select photo from adverts_photos	where advert_id = $1", idInt)
	defer rows.Close()
	for rows.Next() {
		var photo string
		err = rows.Scan(&photo)
		if err != nil {
			return nil, err
		}
		a.Photos = append(a.Photos, photo)
	}
	return &a, nil
}

func (r *AdvertRepo) GetSortedPage(page int, sortMode int) ([]models.AdvertShort, error) {
	ctx := context.Background()
	conn, err := r.db.Connect(ctx)
	defer r.db.Disconnect(ctx)
	if err != nil {
		return nil, err
	}
	sortParams := "id"
	switch sortMode {
	case interfaces.SortByCreatedAsc:
		sortParams = "created asc"
	case interfaces.SortByCreatedDesc:
		sortParams = "created desc"
	case interfaces.SortByPriceAsc:
		sortParams = "price asc"
	case interfaces.SortByPriceDesc:
		sortParams = "price desc"
	}
	startPos := r.cfg.Paginate.AdvertsPageSize * page

	rows, err := conn.Query(ctx,
		`select a.title, a.price, af.photo from adverts a
		left join adverts_photos af  on a.id = af.advert_id and af.delta = 0		
		order by `+sortParams+` limit $1 offset $2`,
		r.cfg.Paginate.AdvertsPageSize, startPos,
	)
	defer rows.Close()
	result := make([]models.AdvertShort, 0, r.cfg.Paginate.AdvertsPageSize)
	for rows.Next() {
		var as models.AdvertShort
		err = rows.Scan(&as.Title, &as.Price, &as.Photo)
		if err != nil {
			return nil, err
		}
		result = append(result, as)
	}
	return result, nil
}
