package service

import (
	"context"
	"errors"
	"fmt"

	"deneme.com/bng-go/Model"
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
	_, err := o.ordercollection.InsertOne(o.ctx, neworder)
	return err
}

// DeleteOrder implements OrderService.
func (o *OrderServiceImpl) DeleteOrder(id *uuid.UUID) error {
	filter := bson.D{bson.E{Key: "_id", Value: id}}
	result, _ := o.ordercollection.DeleteOne(o.ctx, filter)
	if result.DeletedCount == 0 {
		return errors.New("no order found")
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
		cursor, err = o.ordercollection.Find(o.ctx, bson.D{bson.E{Key: "status", Value: status}})
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
