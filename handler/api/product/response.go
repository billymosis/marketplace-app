package product

type createProductResponse struct {
	Name          string   `json:"name"`
	Price         uint     `json:"price"`
	ImageUrl      string   `json:"imageUrl"`
	Stock         uint     `json:"stock"`
	Condition     string   `json:"condition"`
	Tags          []string `json:"tags"`
	IsPurchasable bool     `json:"isPurchasable"`
}

type ProductResponse struct {
	ProductId     string   `json:"productId"`
	Name          string   `json:"name"`
	Price         uint     `json:"price"`
	ImageUrl      string   `json:"imageUrl"`
	Stock         uint     `json:"stock"`
	Condition     string   `json:"condition"`
	Tags          []string `json:"tags"`
	IsPurchasable bool     `json:"isPurchasable"`
	PurchaseCount uint     `json:"purchaseCount"`
}

type GetProductsResponse struct {
	Message string            `json:"message"`
	Data    []ProductResponse `json:"data"`
	Meta    struct {
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
		Total  int `json:"total"`
	} `json:"meta"`
}

type BankAccount struct {
	BankAccountID     string `json:"bankAccountId"`
	BankName          string `json:"bankName"`
	BankAccountName   string `json:"bankAccountName"`
	BankAccountNumber string `json:"bankAccountNumber"`
}

type Seller struct {
	Name             string        `json:"name"`
	ProductSoldTotal int           `json:"productSoldTotal"`
	BankAccounts     []BankAccount `json:"bankAccounts"`
}

type GetProductDetailResponse struct {
	Message string `json:"message"`
	Data    struct {
		Product ProductResponse `json:"product"`
		Seller  Seller          `json:"seller"`
	} `json:"data"`
}
