package Model

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	Id           uuid.UUID          `json:"id" bson:"_id,omitempty"`
	UserId       uuid.UUID          `json:"user_id" bson:"user_id,omitempty"`
	ProductId    uuid.UUID          `json:"product_id" bson:"product_id,omitempty"`
	Quantity     int                `json:"quantity" bson:"quantity,omitempty"`
	Price        float64            `json:"price" bson:"price,omitempty"`
	Status       OrderStatus        `json:"status" bson:"status,omitempty"`
	OrderDate    primitive.DateTime `json:"order_date" bson:"order_date,omitempty"`
	DeliveryDate *time.Time         `json:"delivery_date" bson:"delivery_date,omitempty"`
}

func NewOrder(userId uuid.UUID, quantity int, price float64) Order {
	return Order{
		Id:           uuid.New(),
		UserId:       userId,
		ProductId:    uuid.New(),
		Quantity:     quantity,
		Price:        price,
		Status:       Pending,
		OrderDate:    primitive.NewDateTimeFromTime(time.Now()),
		DeliveryDate: nil,
	}
}

type OrderStatus int

const (
	Pending OrderStatus = iota
	Completed
	Shipped
	Cancelled
)

func (s OrderStatus) String() string {
	return [...]string{"Pending", "Completed", "Shipped", "Cancelled"}[s]
}

func ParseOrderStatus(status string) (os OrderStatus, err error, checkStatus bool) {
	if status == "" {
		return 0, nil, true
	}
	switch status {
	case "Pending":
		return Pending, nil, false
	case "Completed":
		return Completed, nil, false
	case "Shipped":
		return Shipped, nil, false
	case "Cancelled":
		return Cancelled, nil, false
	default:
		return 0, errors.New("invalid status"), false
	}
}

func (s OrderStatus) MarshalBSONJSON() (bsontype.Type, []byte, error) {
	return bson.MarshalValue(s.String())
}

func (s *OrderStatus) UnmarshalBsonValue(t bsontype.Type, data []byte) error {
	var status string
	err := bson.UnmarshalValue(t, data, &status)
	if err != nil {
		return err
	}

	switch status {
	case "Pending":
		*s = Pending
	case "Completed":
		*s = Completed
	case "Shipped":
		*s = Shipped
	case "Cancelled":
		*s = Cancelled
	default:
		return errors.New("invalid status")
	}

	return nil
}
