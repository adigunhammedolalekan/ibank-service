package main

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// accountHandler creates functions to handle incoming http requests
type accountHandler struct {
	repo *accountRepository
}

// newAccountHandler creates a usable accountHandler
func newAccountHandler(repo *accountRepository) *accountHandler {
	return &accountHandler{repo:repo}
}

// createAccountHandler handles create_account http request
func (handler *accountHandler) createAccountHandler(ctx *gin.Context) {
	payload := &Account{}
	if err := ctx.ShouldBindJSON(payload); err != nil {
		ctx.JSON(400, &Response{Success: false, Message: "bad request: malformed request body"})
		return
	}

	newAccount, err := handler.repo.CreateAccount(payload.Name, payload.Email, payload.Password)
	if err != nil {
		ctx.JSON(200, &Response{Success: false, Message: err.Error()})
		return
	}

	ctx.JSON(200, &Response{Success: true, Message: "account created", Data: newAccount})
}

// authenticateAccountHandler handles authenticate_account http request
func (handler *accountHandler) authenticateAccountHandler(ctx *gin.Context) {
	payload := &Account{}
	if err := ctx.ShouldBindJSON(payload); err != nil {
		ctx.JSON(400, &Response{Success: false, Message: "bad request: malformed request body"})
		return
	}

	account, err := handler.repo.GetAccountByAttr("account_number", payload.AccountNumber)
	if err != nil {
		ctx.JSON(404, &Response{Success: false, Message: "404: account not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(payload.Password)); err != nil {
		ctx.JSON(200, &Response{Success: false, Message: "Invalid login credentials"})
		return
	}

	account.Token = handler.repo.generateToken(account.AccountNumber, account.Id.Hex())
	ctx.JSON(200, &Response{Success: true, Message: "authenticated", Data: account})
}

type Response struct {
	Success bool `json:"success"`
	Message string `json:"message"`
	Data interface{} `json:"data"`
}
