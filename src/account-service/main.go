package main

import (
	"context"
	pb "github.com/adigunhammedolalekan/ibank-service/src/account-service/proto"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

func main()  {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("failed to load env var ", err)
	}
	db, err := createDbConnection(os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal("failed to connect to db ", err)
	}

	repo := newAccountRepository(db)
	go func() {
		err := startRpcServer(repo)
		if err != nil {
			log.Fatal("gRPC server error ", err)
		}
	}()

	handler := newAccountHandler(repo)
	router := gin.Default()
	group := router.Group("/api")
	group.POST("/account", handler.createAccountHandler)
	group.POST("/account/authenticate", handler.authenticateAccountHandler)

	port := os.Getenv("PORT")
	if err := router.Run(":" + port); err != nil {
		log.Fatal("account http service error ", err)
	}
}

func startRpcServer(repo *accountRepository) error {
	s, err := net.Listen("tcp", ":9003")
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	pb.RegisterAccountServiceServer(server, newAccountService(repo))
	return server.Serve(s)
}

func createDbConnection(uri string) (*mongo.Client, error) {
	log.Println("Connecting to ", uri)
	ctx := context.Background()
	client, err := mongo.NewClient(
		options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	if err := client.Connect(ctx); err != nil {
		return nil, err
	}

	return client, nil
}