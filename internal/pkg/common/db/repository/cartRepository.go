package repository

import (
	"cart-api/internal/pkg/common/models"
	"errors"

	"github.com/jmoiron/sqlx"
)

const IdAfterInsertMock int = 10

type ICartRepository interface {
	GetById(id int) (*models.Cart, error)
	GetAll() ([]models.Cart, error)
	Create() (models.Cart, error)
	Delete(id int) error
}

type PostgresCartRepository struct {
	pool *sqlx.DB
}

func NewPostgresCartRepository(dbPool *sqlx.DB) *PostgresCartRepository {
	return &PostgresCartRepository{
		pool: dbPool,
	}
}

func (c *PostgresCartRepository) GetById(id int) (*models.Cart, error) {
	existsCartIdRow := c.pool.QueryRowx("SELECT id FROM cart WHERE id = $1 LIMIT 1", id)
	var cartId int

	if err := existsCartIdRow.Scan(&cartId); err != nil {
		return &models.Cart{}, err
	}

	if cartId == 0 {
		return &models.Cart{}, errors.New("cannot find cart by given id")
	}

	cart := models.Cart{}

	itemRows, err := c.pool.Queryx("SELECT * FROM cart_item where cart_id = $1", cartId)

	if err != nil {
		return &models.Cart{}, err
	}

	defer itemRows.Close()

	var items []models.CartItem
	for itemRows.Next() {
		var item models.CartItem

		if err := itemRows.StructScan(&item); err != nil {
			return &models.Cart{}, err
		}

		items = append(items, item)
	}

	cart.Id = id

	if len(items) == 0 {
		cart.Items = make([]models.CartItem, 0)
	} else {
		cart.Items = items
	}

	return &cart, nil
}

func (c *PostgresCartRepository) GetAll() ([]models.Cart, error) {
	var count int
	if err := c.pool.Get(&count, "SELECT count(*) FROM cart;"); err != nil {
		return []models.Cart{}, err
	}

	rows, err := c.pool.Queryx("Select * from cart")

	if err != nil {
		return []models.Cart{}, err
	}

	defer rows.Close()

	carts := make([]models.Cart, 0, count)

	//remove append
	for rows.Next() {
		cart := &models.Cart{}

		if err := rows.StructScan(cart); err != nil {
			return []models.Cart{}, err
		}

		itemRows, err := c.pool.Queryx("SELECT * FROM cart_item where cart_id = $1", cart.Id)
		if err != nil {
			break
		}

		defer itemRows.Close()

		var items []models.CartItem

		for itemRows.Next() {

			var item models.CartItem
			if err := itemRows.StructScan(&item); err != nil {
				break
			}
			items = append(items, item)
		}

		if len(items) == 0 {
			cart.Items = make([]models.CartItem, 0)
		} else {
			cart.Items = items

		}

		carts = append(carts, *cart)
	}

	return carts, nil

}

func (c *PostgresCartRepository) Create() (models.Cart, error) {
	tx := c.pool.MustBegin()
	tx.MustExec("INSERT INTO cart DEFAULT VALUES")
	if err := tx.Commit(); err != nil {
		return models.Cart{}, err
	}

	var createdId int
	if err := c.pool.QueryRowx("SELECT id FROM cart ORDER BY id DESC LIMIT 1").Scan(&createdId); err != nil {
		return models.Cart{}, err
	}
	items := make([]models.CartItem, 0)

	return models.Cart{Id: createdId, Items: items}, nil
}

func (c *PostgresCartRepository) Delete(id int) error {
	tx := c.pool.MustBegin()
	tx.MustExec("DELETE FROM cart WHERE id = $1", id)
	err := tx.Commit()

	if err != nil {
		return err
	}

	return nil
}

func (c *PostgresCartRepository) Update(models.Cart) error {
	return nil
}
