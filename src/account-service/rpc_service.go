package main

import (
	"context"
	pb "github.com/adigunhammedolalekan/ibank-service/src/account-service/proto"
)

type accountService struct {
	repo *accountRepository
}

func newAccountService(repo *accountRepository) *accountService {
	return &accountService{
		repo: repo,
	}
}

func (a *accountService) GetAccount(ctx context.Context, req *pb.GetAccountRequest) (*pb.GetAccountResponse, error) {
	account, err := a.repo.GetAccountByAttr("account_number", req.AccountNumber)
	if err != nil {
		return nil, err
	}

	return &pb.GetAccountResponse{
		Account: &pb.Account{
			Name: account.Name, Email: account.Email, AccountNumber: account.AccountNumber,
		},
	}, nil
}

func (a *accountService) VerifyAccountToken(ctx context.Context, req *pb.VerifyTokenRequest) (*pb.VerifyTokenResponse, error) {
	token, err := a.repo.verifyToken(req.Token)
	if err != nil {
		return nil, err
	}

	account, err := a.repo.GetAccountByAttr("account_number", token.AccountNumber)
	if err != nil {
		return nil, err
	}

	return &pb.VerifyTokenResponse{
		Account: &pb.Account{
			Name: account.Name, Email: account.Email, AccountNumber: account.AccountNumber,
		},
	}, nil
}

