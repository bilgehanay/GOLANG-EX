package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"deneme.com/bng-go/Model"
	rabbitmq "deneme.com/bng-go/RabbitMQ"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserServiceImpl struct {
	usercollection *mongo.Collection
	ctx            context.Context
}

// ListenOrderMessage implements UserService.
func (u *UserServiceImpl) ListenOrderMessage() {
	rabbitmq.ConsumeMessages(func(d amqp.Delivery) {
		var orderMessage map[string]string
		if err := json.Unmarshal(d.Body, &orderMessage); err != nil {
			fmt.Println("Error decoding JSON: ", err)
			return
		}

		userid, err := uuid.Parse(orderMessage["user_id"])
		if err != nil {
			fmt.Println("Error parsing UUID: ", err)
			return
		}

		orderid, err := uuid.Parse(orderMessage["order_id"])
		if err != nil {
			fmt.Println("Error parsing UUID: ", err)
			return
		}

		filter := bson.D{bson.E{Key: "_id", Value: userid}}
		update := bson.D{bson.E{Key: "$push", Value: bson.D{bson.E{Key: "orders", Value: orderid}}}}

		_, err = u.usercollection.UpdateOne(u.ctx, filter, update)
		if err != nil {
			fmt.Println("Error updating user: ", err)
			return
		} else {
			fmt.Println("Order added to user")
		}
	})
}

func NewUserService(usercollection *mongo.Collection, ctx context.Context) UserService {
	return &UserServiceImpl{
		usercollection: usercollection,
		ctx:            ctx,
	}
}

func (u *UserServiceImpl) CreateUser(user *Model.User) error {
	newuser := Model.NewUser(user.Name, user.Surname, user.Email, user.Password, user.Age, user.Address)
	_, err := u.usercollection.InsertOne(u.ctx, newuser)
	return err
}

func (u *UserServiceImpl) GetUser(id *uuid.UUID) (*Model.User, error) {
	var user *Model.User
	query := bson.D{bson.E{Key: "_id", Value: id}}
	err := u.usercollection.FindOne(u.ctx, query).Decode(&user)
	return user, err
}

func (u *UserServiceImpl) GetUsers() ([]*Model.User, error) {
	var users []*Model.User
	cursor, err := u.usercollection.Find(u.ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	for cursor.Next(u.ctx) {
		var user Model.User
		err := cursor.Decode(&user)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}
	cursor.Close(u.ctx)
	if len(users) == 0 {
		return nil, errors.New("No users found")
	}

	return users, nil
}

func (u *UserServiceImpl) UpdateUser(user *Model.User) error {
	filter := bson.D{bson.E{Key: "_id", Value: user.Id}}
	update := bson.D{bson.E{Key: "$set", Value: bson.D{
		bson.E{Key: "email", Value: user.Email},
		bson.E{Key: "password", Value: user.Password},
		bson.E{Key: "address", Value: user.Address},
	}}}
	result, _ := u.usercollection.UpdateOne(u.ctx, filter, update)
	if result.MatchedCount != 1 {
		return errors.New("User not found")
	}
	return nil
}

func (u *UserServiceImpl) DeleteUser(id *uuid.UUID) error {
	filter := bson.D{bson.E{Key: "_id", Value: id}}
	result, _ := u.usercollection.DeleteOne(u.ctx, filter)
	if result.DeletedCount != 1 {
		return errors.New("User not found")
	}
	return nil
}
