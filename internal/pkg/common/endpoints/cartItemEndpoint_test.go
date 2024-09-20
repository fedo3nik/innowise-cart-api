package endpoints

import (
	"bytes"
	"cart-api/internal/pkg/common/models"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPostgresItemRepo struct {
	mock.Mock
}

func (m *MockPostgresItemRepo) Create(item models.ItemDto, cartId int) (models.CartItem, error) {
	args := m.Called()

	itemDto := args.Get(0).(models.ItemDto)
	tmp := models.CartItem{
		Id:       args.Int(1),
		Product:  itemDto.Product,
		Quantity: itemDto.Quantity,
		Cart_id:  1,
	}
	return tmp, m.Called().Error(2)
}
func (m *MockPostgresItemRepo) Delete(cartId, id int) error {
	return nil
}

func TestAddToCart(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")

	repo := &MockPostgresItemRepo{}

	t.Run("Success_create", func(t *testing.T) {
		dto := models.ItemDto{
			Product:  "test product",
			Quantity: 10,
		}

		repo.On("Create").Return(dto, 1, nil)

		handler := http.HandlerFunc(NewCartItemHandler(sqlxDB).AddToCart(repo))

		bytesDto, _ := json.Marshal(dto)

		req, _ := http.NewRequest("http.MethodPost", "/carts/{cartId}/items", bytes.NewBuffer(bytesDto))
		req.SetPathValue("cartId", "1")

		w := httptest.NewRecorder()

		handler(w, req)
		res := w.Result()
		defer res.Body.Close()

		var cartItem models.CartItem
		if err := json.NewDecoder(res.Body).Decode(&cartItem); err != nil {
			t.Fatalf("Error: %v", err)
		}

		expEq := models.CartItem{
			Id:       1,
			Cart_id:  1,
			Product:  "test product",
			Quantity: 10,
		}

		assert.Equal(t, expEq, cartItem)
	})

	t.Run("Error_create", func(t *testing.T) {
		dto := models.Cart{
			Id: 1,
		}

		repo.On("Create").Return(dto, 1, nil)

		handler := http.HandlerFunc(NewCartItemHandler(sqlxDB).AddToCart(repo))

		bytesDto, _ := json.Marshal(dto)

		req, _ := http.NewRequest("http.MethodPost", "/carts/{cartId}/items", bytes.NewBuffer(bytesDto))
		req.SetPathValue("cartId", "1")

		w := httptest.NewRecorder()

		handler(w, req)
		res := w.Result()
		defer res.Body.Close()

		bodyStr, _ := io.ReadAll(res.Body)

		assert.Equal(t, "Cannot decode body\n", string(bodyStr))
		assert.Equal(t, res.StatusCode, http.StatusBadRequest)
	})
}
