package product

import "github.com/billymosis/marketplace-app/model"

type createProductResponse struct {
	Name           string   `json:"name"`
	Price          uint     `json:"price"`
	ImageUrl       string   `json:"imageUrl"`
	Stock          uint     `json:"stock"`
	Condition      string   `json:"condition"`
	Tags           []string `json:"tags"`
	IsPurchaseable bool     `json:"isPurchaseable"`
}

type GetProductsResponse struct {
	Message string     `json:"message"`
	Data    []*model.Product `json:"data"`
	Meta    struct {
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
		Total  int `json:"total"`
	} `json:"meta"`
}
