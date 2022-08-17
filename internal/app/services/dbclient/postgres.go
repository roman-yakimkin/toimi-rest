package dbclient

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"toimi/internal/app/services/configmanager"
)

type PostgresDBClient struct {
	cfg  *configmanager.Config
	conn *pgx.Conn
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

func (c *PostgresDBClient) Connect(ctx context.Context) (*pgx.Conn, error) {
	var err error
	st := c.getDatabaseURL()
	c.conn, err = pgx.Connect(ctx, st)
	if err != nil {
		return nil, err
	}
	return c.conn, nil
}

func (c *PostgresDBClient) Disconnect(ctx context.Context) {
	_ = c.conn.Close(ctx)
}
