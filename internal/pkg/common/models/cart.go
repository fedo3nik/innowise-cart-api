package models

type Cart struct {
	Id    int        `db:"id"`
	Items []CartItem `json:"items"`
}
