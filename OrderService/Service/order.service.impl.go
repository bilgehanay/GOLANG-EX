package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"deneme.com/bng-go/Model"
	rabbitmq "deneme.com/bng-go/RabbitMQ"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type OrderServiceImpl struct {
	ordercollection *mongo.Collection
	ctx             context.Context
}

// CreateOrder implements OrderService.
func (o *OrderServiceImpl) CreateOrder(order *Model.Order) error {
	neworder := Model.NewOrder(order.UserId, order.Quantity, order.Price)

	orderJSON, err := json.Marshal(neworder)
	if err != nil {
		return err
	}
	fmt.Println(string(orderJSON))

	result, err := o.ordercollection.InsertOne(o.ctx, neworder)
	if err != nil {
		return err
	}

	fmt.Printf("Inserted document ID: %v\n", result.InsertedID)

	orderMessage := map[string]interface{}{
		"user_id":  neworder.UserId.String(),
		"order_id": neworder.Id.String(),
	}

	orderMessageBytes, err := json.Marshal(orderMessage)
	if err != nil {
		return err
	}

	err = rabbitmq.PublishMessage(string(orderMessageBytes), "new")
	if err != nil {
		return err
	}

	return nil
}

// DeleteOrder implements OrderService.
func (o *OrderServiceImpl) DeleteOrder(id *uuid.UUID) error {
	var order Model.Order
	filter := bson.D{bson.E{Key: "_id", Value: id}}
	err := o.ordercollection.FindOne(o.ctx, filter).Decode(&order)
	if err != nil {
		return err
	}

	result, _ := o.ordercollection.DeleteOne(o.ctx, filter)
	if result.DeletedCount == 0 {
		return errors.New("no order found")
	}

	orderMessage := map[string]interface{}{
		"user_id":  order.UserId.String(),
		"order_id": order.Id.String(),
	}

	orderMessageBytes, err := json.Marshal(orderMessage)
	if err != nil {
		return err
	}

	err = rabbitmq.PublishMessage(string(orderMessageBytes), "delete")
	if err != nil {
		return err
	}

	return nil
}

// GetOrder implements OrderService.
func (o *OrderServiceImpl) GetOrder(id *uuid.UUID) (*Model.Order, error) {
	var order *Model.Order
	query := bson.D{bson.E{Key: "_id", Value: id}}
	err := o.ordercollection.FindOne(o.ctx, query).Decode(&order)
	return order, err
}

// GetOrders implements OrderService.
func (o *OrderServiceImpl) GetOrders(status *Model.OrderStatus, checkStatus *bool) ([]*Model.Order, error) {
	var orders []*Model.Order
	var cursor *mongo.Cursor
	var err error
	if *checkStatus {
		cursor, err = o.ordercollection.Find(o.ctx, bson.D{})
	} else {
		cursor, err = o.ordercollection.Find(o.ctx, bson.D{bson.E{Key: "status", Value: status.String()}})
	}
	if err != nil {
		return nil, err
	}

	for cursor.Next(o.ctx) {
		var order Model.Order
		err := cursor.Decode(&order)
		if err != nil {
			return nil, err
		}
		orders = append(orders, &order)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	cursor.Close(o.ctx)
	if len(orders) == 0 {
		return nil, errors.New("no orders found")
	}
	return orders, nil
}

// UpdateOrder implements OrderService.
func (o *OrderServiceImpl) UpdateOrder(id *uuid.UUID, updateData interface{}) error {
	filter := bson.D{bson.E{Key: "_id", Value: id}}
	update := bson.D{bson.E{Key: "$set", Value: updateData}}
	fmt.Println(update)
	result, _ := o.ordercollection.UpdateOne(o.ctx, filter, update)
	if result.MatchedCount != 1 {
		return errors.New("order not found")
	}
	return nil
}

// UpdateStatus implements OrderService.
func (o *OrderServiceImpl) UpdateStatus(id *uuid.UUID, status Model.OrderStatus) error {
	filter := bson.D{bson.E{Key: "_id", Value: id}}
	update := bson.D{bson.E{Key: "$set", Value: bson.D{bson.E{Key: "status", Value: status}}}}
	result, _ := o.ordercollection.UpdateOne(o.ctx, filter, update)
	if result.MatchedCount != 1 {
		return errors.New("order not found")
	}
	return nil
}

func NewOrderService(ordercollection *mongo.Collection, ctx context.Context) OrderService {
	return &OrderServiceImpl{
		ordercollection: ordercollection,
		ctx:             ctx,
	}
}
