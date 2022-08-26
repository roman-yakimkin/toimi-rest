package dbclient

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"toimi/internal/app/services/configmanager"
)

type PostgresDBClient struct {
	cfg  *configmanager.Config
	pool *pgxpool.Pool
}

func NewPostgresDBClient(cfg *configmanager.Config) *PostgresDBClient {
	return &PostgresDBClient{
		cfg: cfg,
	}
}

func (c *PostgresDBClient) getDatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.cfg.DB.Username,
		c.cfg.DB.Password,
		c.cfg.DB.Host,
		c.cfg.DB.Port,
		c.cfg.DB.DBName,
		c.cfg.DB.SSLMode)
}

func (c *PostgresDBClient) Connect(ctx context.Context) (*pgxpool.Pool, error) {
	var err error
	st := c.getDatabaseURL()
	c.pool, err = pgxpool.Connect(ctx, st)
	if err != nil {
		return nil, err
	}
	return c.pool, nil
}

func (c *PostgresDBClient) Disconnect(ctx context.Context) {
	c.pool.Close()
}
