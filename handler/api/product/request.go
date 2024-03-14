package product

type createProductRequest struct {
	Name           string   `json:"name" validate:"required,min=5,max=60"`
	Price          int      `json:"price" validate:"required,min=0"`
	ImageURL       string   `json:"imageUrl" validate:"required,url"`
	Stock          int      `json:"stock" validate:"min=0"`
	Condition      string   `json:"condition" validate:"required,oneof=new second"`
	Tags           []string `json:"tags" validate:"required,min=0,dive,required"`
	IsPurchaseable bool     `json:"isPurchaseable" validate:"required"`
}
