package models

type CartItem struct {
	Id       int    `db:"id" json:"id"`
	Product  string `db:"product" json:"product"`
	Quantity int    `db:"quantity" json:"quantity"`
	Cart_id  int    `db:"cart_id" json:"cart_id"`
}
