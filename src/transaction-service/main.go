package main

import (
	pb "github.com/adigunhammedolalekan/ibank-service/src/account-service/proto"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"log"
	"os"
)

const accountServiceRpcUrl = "account-service:9003"
func main()  {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("failed to load env data ", err)
	}

	db, err := createDbConnection(os.Getenv("DB_URL"))
	conn, err := grpc.Dial(accountServiceRpcUrl, grpc.WithInsecure())
	if err != nil {
		log.Fatal("failed to connect to account service ", err)
	}

	client := pb.NewAccountServiceClient(conn)
	repo := newTransactionRepository(db, client)
	handler := newTransactionHandler(repo)

	router := gin.Default()
	group := router.Group("/api/txn")
	group.POST("/transfer", handler.transferHandler)
	group.POST("/fund", handler.fundAccountHandler)
	group.GET("/me/history", handler.transactionHistoryHandler)

	port := os.Getenv("PORT")
	if err := router.Run(":" + port); err != nil {
		log.Fatal("transaction http service error ", err)
	}
}

func createDbConnection(uri string) (*gorm.DB, error) {
	db, err := gorm.Open("mysql", uri)
	if err != nil {
		return nil, err
	}

	return db, nil
}