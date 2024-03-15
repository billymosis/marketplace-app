package account

type GetAccountResponse struct {
	Message string        `json:"message"`
	Data    []BankAccount `json:"data"`
}

type BankAccount struct {
	BankAccountId     string `json:"bankAccountId"`
	BankName          string `json:"bankName"`
	BankAccountName   string `json:"bankAccountName"`
	BankAccountNumber string `json:"bankAccountNumber"`
}
