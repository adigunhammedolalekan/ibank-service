package main

import (
	"context"
	"errors"
	pb "github.com/adigunhammedolalekan/ibank-service/src/account-service/proto"
	"github.com/jinzhu/gorm"
)

type transactionRepository struct {
	db *gorm.DB
	client pb.AccountServiceClient
}

func newTransactionRepository(db *gorm.DB, client pb.AccountServiceClient) *transactionRepository {
	return &transactionRepository{db: db, client:client}
}

func (repo *transactionRepository) doTransfer(fromAccount, toAccount string, amount int64) (*Transaction, error) {
	fromWallet, err := repo.wallet(fromAccount)
	if err != nil {
		return nil, err
	}
	toWallet, err := repo.wallet(toAccount)
	if err != nil {
		return nil, err
	}

	if fromWallet.Balance < amount {
		return nil, errors.New("insufficient balance")
	}

	tx := repo.db.Begin()
	if err := tx.Error; err != nil {
		return nil, err
	}

	credit := toWallet.Balance + amount
	debit := fromWallet.Balance - amount

	err = tx.Table("wallets").Where("account_number = ?", fromAccount).UpdateColumn("balance", debit).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Table("wallets").Where("account_number = ?", toAccount).UpdateColumn("balance", credit).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	txn := newTransaction(fromAccount, toAccount, amount)
	if err := tx.Create(txn).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return txn, nil
}

func (repo *transactionRepository) history(account string) ([]*Transaction, error) {
	data := make([]*Transaction, 0)
	err := repo.db.Table("transactions").Where("account_number = ?", account).Find(&data).Error
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (repo *transactionRepository) fundAccount(account string, amount int64) (*Wallet, error) {
	w, err := repo.wallet(account)
	if err != nil {
		return nil, err
	}

	newBalance := w.Balance + amount
	err = repo.db.Table("wallets").Where("account_number = ?", account).UpdateColumn("balance", newBalance).Error
	if err != nil {
		return nil, err
	}

	return repo.wallet(account)
}

func (repo *transactionRepository) createWallet(account string) (*Wallet, error) {
	w := newWallet(account)
	if err := repo.db.Create(w).Error; err != nil {
		return w, nil
	}

	return w, nil
}

// wallet fetches wallet where accountNumber == account params
// a wallet will be created for this account if doesn't exists before
func (repo *transactionRepository) wallet(account string) (*Wallet, error) {
	w := &Wallet{}
	err := repo.db.Table("wallets").Where("account_number = ?", account).First(w).Error
	if err != nil || err == gorm.ErrRecordNotFound {
		err = repo.db.Create(newWallet(account)).Error
		if err != nil {
			return nil, err
		}
	}

	return w, nil
}

func (repo *transactionRepository) verifyToken(token string) (*pb.Account, error) {
	res, err := repo.client.VerifyAccountToken(context.Background(), &pb.VerifyTokenRequest{Token:token})
	if err != nil {
		return nil, err
	}

	return res.Account, nil
}