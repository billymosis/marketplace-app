package product

type createProductRequest struct {
	Name          string   `json:"name" validate:"required,min=5,max=60"`
	Price         int      `json:"price" validate:"required,min=0"`
	ImageURL      string   `json:"imageUrl" validate:"required,url"`
	Stock         int      `json:"stock" validate:"min=0"`
	Condition     string   `json:"condition" validate:"required,oneof=new second"`
	Tags          []string `json:"tags" validate:"required,min=0,dive,required"`
	IsPurchasable bool     `json:"isPurchasable" validate:"required"`
}

type updateProductStockRequest struct {
	Stock int `json:"stock" validate:"min=0"`
}

type buyProductRequest struct {
	BankAccountId        string `json:"bankAccountId" validate:"required"`
	PaymentProofImageUrl string `json:"paymentProofImageUrl" validate:"required,http_url"`
	Quantity             int    `json:"quantity" validate:"min=1"`
}
