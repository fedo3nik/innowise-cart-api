package middleware

import (
	"bytes"
	"cart-api/internal/pkg/common/models"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"regexp"
	"strings"
)

type ValidateItemMiddleware struct {
	nextHandleFunc http.HandlerFunc
}

func NewValiDateMiddleWare(handlerFunc http.HandlerFunc) *ValidateItemMiddleware {
	return &ValidateItemMiddleware{
		nextHandleFunc: handlerFunc,
	}
}

func (v *ValidateItemMiddleware) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
	}

	req.Body = io.NopCloser(bytes.NewReader(body))

	var dto models.ItemDto
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(&dto); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if !validateProduct(dto.Product) {
		http.Error(res, "Product cannot be blank ", http.StatusBadRequest)
		return
	}

	if !validateQuantity(dto.Quantity) {
		http.Error(res, "Quantity need to be positive", http.StatusBadRequest)
		return
	}

	valid, err := validateProductReg(dto.Product)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	if !valid {
		http.Error(res, "Product must have letters", http.StatusBadRequest)
		return
	}
	v.nextHandleFunc(res, req)
}

func validateProduct(product string) bool {
	if len([]rune(product)) < 1 {
		return false
	}

	if strings.EqualFold(product, " ") {
		return false
	}
	return true
}

func validateProductReg(product string) (bool, error) {
	pattern := `.*[a-zA-Z].*`
	re, err := regexp.Compile(pattern)
	if err != nil {
		return false, err
	}
	match := re.MatchString(product)

	if !match {
		return false, errors.New("product need to be valid(word)")
	}

	return match, nil
}

func validateQuantity(quantity int) bool {
	return quantity >= 1
}
