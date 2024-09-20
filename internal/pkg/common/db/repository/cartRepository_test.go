package repository

import (
	"cart-api/internal/pkg/common/models"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestPostgresCartRepository_GetById(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")

	repo := NewPostgresCartRepository(sqlxDB)

	cartId := 1

	cartItem := &models.CartItem{
		Id:       1,
		Cart_id:  cartId,
		Product:  "test product",
		Quantity: 10,
	}

	expectedCart := models.Cart{
		Id:    cartId,
		Items: []models.CartItem{*cartItem},
	}

	t.Run("Success", func(t *testing.T) {
		// Mock if exist cart query
		mock.ExpectQuery("SELECT id FROM cart WHERE id = \\$1 LIMIT 1").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		// mock select cart items
		mock.ExpectQuery("SELECT \\*\\ FROM cart_item where cart_id = \\$1").
			WithArgs(cartId).
			WillReturnRows(sqlmock.NewRows([]string{"id", "product", "quantity", "cart_id"}).
				AddRow(cartItem.Id, cartItem.Product, cartItem.Quantity, cartItem.Cart_id))

		returnedCart, err := repo.GetById(1)

		assert.NoError(t, err)
		assert.Equal(t, &expectedCart, returnedCart)
	})

	t.Run("StructScan error", func(t *testing.T) {
		mock.ExpectQuery("SELECT id FROM cart WHERE id = \\$1 LIMIT 1").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(0))

		_, err := repo.GetById(cartId)

		if assert.Error(t, err) {
			assert.Equal(t, "cannot find cart by given id", err.Error())
		}

	})

	assert.NoError(t, mock.ExpectationsWereMet())
}
