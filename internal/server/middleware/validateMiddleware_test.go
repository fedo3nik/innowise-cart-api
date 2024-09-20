package middleware

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateProductEmptyProduct(t *testing.T) {
	val := validateProduct("")
	assert.False(t, val)
}

func TestValidateProductLenLowerOne(t *testing.T) {
	val := validateProduct("r")
	assert.True(t, val)
}

func TestValidateProductNormal(t *testing.T) {
	val := validateProduct("Shoes")
	assert.True(t, val)
}

func TestValidateProductRegOnlyNumbers(t *testing.T) {
	val, err := validateProductReg("12345")
	assert.False(t, val)

	if assert.Error(t, err) {
		assert.EqualError(t, err, "product need to be valid(word)")
	}
}

func TestValidateProductRegOnlyBlank(t *testing.T) {
	val, err := validateProductReg("     ")
	assert.False(t, val)

	if assert.Error(t, err) {
		assert.EqualError(t, err, "product need to be valid(word)")
	}
}

func TestValidateProductRegNormalText(t *testing.T) {
	val, err := validateProductReg("Basketball ring")
	assert.True(t, val)
	assert.Nil(t, err)
}

func TestValidateQuantityZero(t *testing.T) {
	val := validateQuantity(0)
	assert.False(t, val)
}

func TestValidateQuantityPositive(t *testing.T) {
	val := validateQuantity(1)
	assert.True(t, val)
}
