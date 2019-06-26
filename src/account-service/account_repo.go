package main

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"os"
)

type accountRepository struct {
	mongoClient *mongo.Client
}

func newAccountRepository(client *mongo.Client) *accountRepository {
	return &accountRepository{mongoClient:client}
}

func (repo *accountRepository) DB() *mongo.Database {
	return repo.mongoClient.Database("accounts_database")
}

func (repo *accountRepository) Col() *mongo.Collection {
	return repo.DB().Collection("accounts")
}

func (repo *accountRepository) GetAccountByAttr(attr string, value interface{}) (*Account, error) {
	if attr == "_id" {
		return repo.getAccountByObjectId(value . (string))
	}

	filter := bson.M{attr : bson.M {"$eq" : value}}
	col := repo.Col()
	var account Account
	if err := col.FindOne(context.Background(), filter).Decode(&account); err != nil {
		return nil, err
	}

	return &account, nil
}

func (repo *accountRepository) getAccountByObjectId(id string) (*Account, error) {
	var account Account
	filter := bson.M {"_id" : id}
	if err := repo.Col().FindOne(context.Background(), filter).Decode(&account); err != nil {
		return nil, err
	}

	return &account, nil
}

func (repo *accountRepository) CreateAccount(name, email, password string) (*Account, error) {
	account := newAccount(name, email, password)
	if exists := repo.accountExists(account.Email); exists {
		return nil, errors.New("email is already in use")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	account.Password = string(hashedPassword)
	col := repo.Col()
	result, err := col.InsertOne(context.Background(), account)
	if err != nil {
		return nil, err
	}

	id := result.InsertedID . (primitive.ObjectID).Hex()
	newAccount, err := repo.getAccountByObjectId(id)
	if err != nil {
		return nil, err
	}

	newAccount.Token = repo.generateToken(newAccount.AccountNumber, newAccount.Id.Hex())
	return newAccount, nil
}

func (repo *accountRepository) generateToken(accountNumber, accountId string) string {
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), newToken(accountNumber, accountId))
	tkString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return ""
	}

	return tkString
}

func (repo *accountRepository) verifyToken(tokenString string) (*Token, error) {
	tk := &Token{}
	token, err := jwt.ParseWithClaims(tokenString, tk, func(token *jwt.Token) (i interface{}, e error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("token is invalid")
	}

	return tk, nil
}

func (repo *accountRepository) accountExists(email string) bool {
	_, err := repo.GetAccountByAttr("email", email)
	if err != nil {
		return false
	}

	return true
}