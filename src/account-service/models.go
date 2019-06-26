package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Account struct {
	Id primitive.ObjectID `bson:"id" json:"id"`
	Name string `bson:"name" json:"name"`
	Email string `bson:"email" json:"email"`
	Password string `bson:"password" json:"password"`
	AccountNumber string `bson:"account_number" json:"account_number"`
	Token string `json:"token"`
}

type Token struct {
	jwt.StandardClaims
	AccountNumber string
	AccountId string
}

func newToken(number, id string) *Token {
	return &Token{AccountNumber: number, AccountId: id}
}

func newAccount(name, email, password string) *Account {
	return &Account{Name: name, Email: email, Password: password, AccountNumber: newAccountNumber()}
}

// newAccountNumber returns a random 8-length long number string
func newAccountNumber() string {
	return fmt.Sprintf("%08d", rand.Intn(99999999))
}