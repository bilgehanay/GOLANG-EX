package main

import (
	"context"
	"deneme.com/bng-go/Init"
	"fmt"
	"log"
	"os"

	controller "deneme.com/bng-go/Controller"
	rabbitmq "deneme.com/bng-go/RabbitMQ"
	service "deneme.com/bng-go/Service"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	server         *gin.Engine
	userService    service.UserService
	userController controller.UserController
	ctx            context.Context
	usercollection *mongo.Collection
	mongoClient    *mongo.Client
	err            error
)

func init() {
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "development"
	}
	err = Init.LoadConfig(env)

	if err != nil {
		log.Fatal("Error loading config: ", err)
	}

	ctx = context.TODO()

	mongoconn := options.Client().ApplyURI(Init.SetConfig.MongoDB.Connection)
	mongoClient, err = mongo.Connect(ctx, mongoconn)
	if err != nil {
		log.Fatal(err)
	}
	err = mongoClient.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB")

	usercollection = mongoClient.Database(Init.SetConfig.MongoDB.Database).Collection(Init.SetConfig.MongoDB.Collection)
	userService = service.NewUserService(usercollection, ctx)
	userController = controller.New(userService)
	server = gin.Default()
	rabbitmq.InitRabbitMQ(Init.SetConfig.RabbitMQ.Connection)
}

func main() {

	defer mongoClient.Disconnect(ctx)
	defer rabbitmq.Close()

	go userService.ListenOrderMessage()

	basepath := server.Group("/api/v1")
	userController.RegisterUserRoutes(basepath)
	log.Fatal(server.Run(Init.SetConfig.App.Url))
}
