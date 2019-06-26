package main

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
)

type Transaction struct {
	gorm.Model
	AccountNumber string `json:"account_number"`
	RecipientAccountNumber string `json:"recipient_account_number"`
	Amount int64 `json:"amount"`
}

type Wallet struct {
	gorm.Model
	AccountNumber string `json:"account_number"`
	Balance int64 `json:"balance"`
}

func newWallet(account string) *Wallet {
	return &Wallet{AccountNumber: account, Balance: 0}
}

func newTransaction(fromAccount, toAccount string, amount int64) *Transaction {
	return &Transaction{
		AccountNumber: fromAccount, RecipientAccountNumber: toAccount, Amount: amount,
	}
}

type TransferRequestPayload struct {
	ToAccount string `json:"to_account"`
	Amount json.Number `json:"amount"`
}

type FundAccountPayload struct {
	Amount json.Number `json:"amount"`
}

func (t *TransferRequestPayload) amount() int64 {
	i, err := t.Amount.Int64()
	if err != nil {
		return 0
	}

	return i
}


func (t *FundAccountPayload) amount() int64 {
	i, err := t.Amount.Int64()
	if err != nil {
		return 0
	}

	return i
}