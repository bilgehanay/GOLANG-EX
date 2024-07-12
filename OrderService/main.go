package main

import (
	"context"
	"fmt"
	"log"

	controller "deneme.com/bng-go/Controller"
	rabbitmq "deneme.com/bng-go/RabbitMQ"
	service "deneme.com/bng-go/Service"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	server          *gin.Engine
	orderService    service.OrderService
	orderController controller.OrderController
	ctx             context.Context
	ordercollection *mongo.Collection
	mongoClient     *mongo.Client
	err             error
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

	ordercollection = mongoClient.Database("orderdb").Collection("orders")
	orderService = service.NewOrderService(ordercollection, ctx)
	orderController = controller.New(orderService)
	server = gin.Default()
	rabbitmq.InitRabbitMQ()
}

func main() {
	defer mongoClient.Disconnect(ctx)
	defer rabbitmq.Close()

	basepath := server.Group("/api/v1")
	orderController.RegisterOrderRoutes(basepath)

	log.Fatal(server.Run(":9090"))
}
