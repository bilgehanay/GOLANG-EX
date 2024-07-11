package service

import (
	"deneme.com/bng-go/Model"
	"github.com/google/uuid"
)

type OrderService interface {
	CreateOrder(*Model.Order) error
	GetOrder(*uuid.UUID) (*Model.Order, error)
	GetOrders(status *Model.OrderStatus, checkStatus *bool) ([]*Model.Order, error)
	UpdateOrder(id *uuid.UUID, updateData interface{}) error
	DeleteOrder(*uuid.UUID) error

	UpdateStatus(*uuid.UUID, Model.OrderStatus) error
}
