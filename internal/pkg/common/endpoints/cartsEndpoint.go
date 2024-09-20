package endpoints

import (
	"cart-api/internal/pkg/common/db/repository"
	"cart-api/internal/pkg/common/models"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jmoiron/sqlx"
)

type CartHandler struct {
	pool *sqlx.DB
}

func NewCarHandler(db *sqlx.DB) *CartHandler {
	return &CartHandler{
		pool: db,
	}
}

// CreateCart creates a new cart
//
//	@Summary		Creates a new cart and id generated
//	@Description	create cart
//	@Tags			carts
//	@Produce		json
//	@Success		200	{object}	models.Cart
//	@Failure		500	{object}	models.ResponseError
//	@Router			/carts [post]
func (c *CartHandler) CreateNew(repo repository.ICartRepository) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		cart, err := repo.Create()

		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)

		json.NewEncoder(res).Encode(cart)
	}
}

// ListCarts lists all existing carts with items
//
//	@Summary		List carts
//	@Description	get all carts with composition item slice inside
//	@Tags			carts
//	@Produce		json
//	@Success		200	{array}		models.Cart
//	@Failure		500	{object}	models.ResponseError
//	@Router			/carts [get]
func (c *CartHandler) GetAll(repo repository.ICartRepository) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		carts, err := repo.GetAll()

		if err != nil {
			models.NewResponseError(http.StatusInternalServerError, err.Error()).ShowError(res)
			return
		}

		if len(carts) == 0 {
			res.Write([]byte("{}"))
		}

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)

		json.NewEncoder(res).Encode(carts)
	}
}

// ViewCart shows cart
//
//	@Summary		Shows cart by id
//	@Description	get cart by id
//	@Tags			carts
//	@Produce		json
//	@Param			id	path		int	true	"id to find cart"
//	@Success		200	{object}	models.Cart
//	@Failure		400	{object}	models.ResponseError
//	@Failure		404	{object}	models.ResponseError
//	@Router			/carts/{id} [get]
func (c *CartHandler) ViewCart(repo repository.ICartRepository) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {

		idPath := req.PathValue("id")

		id, err := strconv.Atoi(idPath)

		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		cart, err := repo.GetById(id)

		if err != nil {
			http.Error(res, err.Error(), http.StatusNotFound)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)

		json.NewEncoder(res).Encode(cart)
	}
}

// DeleteCart  delete cart
//
//	@Summary		Delete a cart
//	@Description	delete a cart with its items recursively
//	@Tags			carts
//	@Produce		json
//	@Param			id	path		int	true	"Cart id"
//	@Success		200	{array}		byte
//	@Failure		400	{object}	models.ResponseError
//	@Failure		500	{object}	models.ResponseError
//	@Router			/carts/{id} [delete]
func (c *CartHandler) DeleteCart(repository repository.ICartRepository) func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		pathId := req.PathValue("id")
		id, err := strconv.Atoi(pathId)

		if err != nil {
			models.NewResponseError(http.StatusBadRequest, err.Error())
			return
		}

		if err := repository.Delete(id); err != nil {
			models.NewResponseError(http.StatusInternalServerError, err.Error())
			return
		}

		res.WriteHeader(http.StatusOK)
		res.Header().Set("Content-Type", "application/json")
		res.Write([]byte("{}"))
	}
}
