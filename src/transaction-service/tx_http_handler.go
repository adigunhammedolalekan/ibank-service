package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const headerKey = "Authorization"
type transactionHandler struct {
	repo *transactionRepository
}

func newTransactionHandler(repo *transactionRepository) *transactionHandler {
	return &transactionHandler{repo:repo}
}

func (handler *transactionHandler) transferHandler(ctx *gin.Context) {
	token := handler.token(ctx)
	if token == "" {
		ctx.JSON(http.StatusUnauthorized, &Response{Success: false, Message: "unauthorized request: token is missing"})
		return
	}

	account, err := handler.repo.verifyToken(token)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, &Response{Success: false, Message: "unauthorized request: " + err.Error()})
		return
	}

	payload := &TransferRequestPayload{}
	if err := ctx.ShouldBindJSON(payload); err != nil {
		ctx.JSON(http.StatusBadRequest, &Response{Success: false, Message: "bad request: malformed request body"})
		return
	}

	txn, err := handler.repo.doTransfer(account.AccountNumber, payload.ToAccount, payload.amount())
	if err != nil {
		ctx.JSON(200, &Response{Success: false, Message: err.Error()})
		return
	}

	ctx.JSON(200, &Response{Success: true, Message: "transaction successful", Data: txn})
}

func (handler *transactionHandler) transactionHistoryHandler(ctx *gin.Context) {
	token := handler.token(ctx)
	if token == "" {
		ctx.JSON(http.StatusUnauthorized, &Response{Success: false, Message: "unauthorized request: token is missing"})
		return
	}

	account, err := handler.repo.verifyToken(token)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, &Response{Success: false, Message: "unauthorized request: " + err.Error()})
		return
	}

	txns, err := handler.repo.history(account.AccountNumber)
	if err != nil {
		ctx.JSON(200, &Response{Success: false, Message: err.Error()})
		return
	}

	ctx.JSON(200, &Response{Success: true, Message: "history", Data: txns})
}

func (handler *transactionHandler) fundAccountHandler(ctx *gin.Context)  {
	token := handler.token(ctx)
	if token == "" {
		ctx.JSON(http.StatusUnauthorized, &Response{Success: false, Message: "unauthorized request: token is missing"})
		return
	}

	account, err := handler.repo.verifyToken(token)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, &Response{Success: false, Message: "unauthorized request: " + err.Error()})
		return
	}
	payload := &FundAccountPayload{}
	if err := ctx.ShouldBindJSON(payload); err != nil {
		ctx.JSON(http.StatusBadRequest, &Response{Success: false, Message: "bad request: malformed request body"})
		return
	}
	w, err := handler.repo.fundAccount(account.AccountNumber, payload.amount())
	if err != nil {
		ctx.JSON(200, &Response{Success: false, Message: err.Error()})
		return
	}

	ctx.JSON(200, &Response{Success: true, Message: "account funded", Data: w})
}

func (handler *transactionHandler) token(ctx *gin.Context) string {
	header := ctx.GetHeader(headerKey)
	values := strings.Split(header, " ")
	if len(values) != 2 {
		return ""
	}

	return values[1]
}

type Response struct {
	Success bool `json:"success"`
	Message string `json:"message"`
	Data interface{} `json:"data"`
}