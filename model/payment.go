package model

type Payment struct {
	Id                   uint
	AccountId            uint
	ProductId            uint
	PaymentProofImageUrl string
	Quantity             uint
}
