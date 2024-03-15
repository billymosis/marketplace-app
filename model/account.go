package model

type Account struct {
	Id            uint   `json:"id"`
	Name          string `json:"name"`
	AccountName   string `json:"accountName"`
	AccountNumber string `json:"accountNumber"`
	UserId        uint   `json:"userId"`
}
