package model

type Product struct {
	Id            uint     `json:"id"`
	Name          string   `json:"name"`
	Price         uint     `json:"price"`
	ImageUrl      string   `json:"imageUrl"`
	Stock         uint     `json:"stock"`
	Condition     string   `json:"condition"`
	Tags          []string `json:"tags"`
	IsPurchasable bool     `json:"isPurchasable"`
	PurchaseCount uint     `json:"purchaseCount"`
	UserId        uint     `json:"userId"`
}
