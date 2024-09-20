package repository

import (
	"cart-api/internal/pkg/common/models"
	"errors"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type ItemRepository interface {
	Create(item models.ItemDto, cartId int) (models.CartItem, error)
	Delete(cartId, id int) error
}

type PostgresItemRepository struct {
	pool *sqlx.DB
}

func NewPostgresItemRepository(dbPool *sqlx.DB) *PostgresItemRepository {
	return &PostgresItemRepository{
		pool: dbPool,
	}
}

func (c *PostgresItemRepository) Create(item models.ItemDto, cartId int) (models.CartItem, error) {
	var count int
	if err := c.pool.Get(&count, "SELECT count(*) FROM cart WHERE id = $1", cartId); err != nil {
		return models.CartItem{}, err
	}

	if count == 0 {
		return models.CartItem{}, errors.New("cart id now found")
	}

	tx := c.pool.MustBegin()
	tx.MustExec("INSERT INTO cart_item (id, product, quantity, cart_id) VALUES (DEFAULT, $1, $2, $3);", item.Product, item.Quantity, cartId)
	err := tx.Commit()

	if err != nil {
		return models.CartItem{}, nil
	}

	row := c.pool.QueryRowx("SELECT * FROM cart_item ORDER BY id DESC LIMIT 1")

	var rowItem models.CartItem
	if err := row.StructScan(&rowItem); err != nil {
		return models.CartItem{}, err
	}

	return rowItem, nil
}

func (c *PostgresItemRepository) Delete(cartId, id int) error {
	tx := c.pool.MustBegin()
	tx.MustExec("DELETE FROM cart_item WHERE id = $1 AND cart_id = $2", id, cartId)
	err := tx.Commit()

	if err != nil {
		return err
	}

	return nil
}
