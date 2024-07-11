package main

import (
	"context"
	"fmt"
	"log"

	controller "deneme.com/bng-go/Controller"
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
	ctx = context.TODO()

	mongoconn := options.Client().ApplyURI("mongodb://localhost:27017")
	mongoClient, err = mongo.Connect(ctx, mongoconn)
	if err != nil {
		log.Fatal(err)
	}
	err = mongoClient.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB")

	usercollection = mongoClient.Database("userdb").Collection("users")
	userService = service.NewUserService(usercollection, ctx)
	userController = controller.New(userService)
	server = gin.Default()
}

func main() {
	defer mongoClient.Disconnect(ctx)

	basepath := server.Group("/api/v1")
	userController.RegisterUserRoutes(basepath)

	log.Fatal(server.Run(":8080"))
}
