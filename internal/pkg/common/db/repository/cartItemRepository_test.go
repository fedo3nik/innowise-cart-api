package repository

import (
	"cart-api/internal/pkg/common/models"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestPostgresItemRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewPostgresItemRepository(sqlxDB)

	t.Run("Delete_Success", func(t *testing.T) {
		cartId := 1
		itemId := 2

		mock.ExpectBegin()
		mock.ExpectExec("DELETE FROM cart_item WHERE id = \\$1 AND cart_id = \\$2").
			WithArgs(itemId, cartId).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err = repo.Delete(cartId, itemId)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Delete_Error", func(t *testing.T) {
		cartId := 1
		itemId := 2

		mock.ExpectBegin()
		mock.ExpectExec("DELETE FROM cart_item WHERE id = \\$1 AND cart_id = \\$2").
			WithArgs(itemId, cartId).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit().WillReturnError(sqlmock.ErrCancelled)

		err = repo.Delete(cartId, itemId)

		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestPostgresItemRepository_Create(t *testing.T) {
	// Create a new sqlmock database connection
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")

	repo := NewPostgresItemRepository(sqlxDB)

	itemDto := &models.ItemDto{
		Product:  "Test Product",
		Quantity: 2,
	}

	t.Run("Cart not found", func(t *testing.T) {
		// Mock the count query
		mock.ExpectQuery("SELECT count\\(\\*\\) FROM cart WHERE id = \\$1").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0)) // No cart found

		_, err := repo.Create(*itemDto, 1)

		assert.Error(t, err)
		assert.Equal(t, "cart id now found", err.Error())
	})

	t.Run("Successful creation", func(t *testing.T) {
		// Mock the count query
		mock.ExpectQuery("SELECT count\\(\\*\\) FROM cart WHERE id = \\$1").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1)) // Cart found

		// Mock the insert statement
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO cart_item").
			WithArgs(itemDto.Product, itemDto.Quantity, 1).
			WillReturnResult(sqlmock.NewResult(1, 1)) // Simulate successful insert
		mock.ExpectCommit()

		// Mock the select for the last inserted item
		mock.ExpectQuery("SELECT \\* FROM cart_item ORDER BY id DESC LIMIT 1").
			WillReturnRows(sqlmock.NewRows([]string{"id", "product", "quantity", "cart_id"}).
				AddRow(1, itemDto.Product, itemDto.Quantity, 1))

		// Call the Create method
		createdItem, err := repo.Create(*itemDto, 1)

		assert.NoError(t, err)
		assert.Equal(t, models.CartItem{Id: 1, Product: itemDto.Product, Quantity: itemDto.Quantity, Cart_id: 1}, createdItem)
	})

	// Ensure all expectations were met
	assert.NoError(t, mock.ExpectationsWereMet())
}
