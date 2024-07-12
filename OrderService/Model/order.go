package Model

import (
	"encoding/json"
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

type OrderStatus string

const (
	Pending   OrderStatus = "pending"
	Completed OrderStatus = "completed"
	Shipped   OrderStatus = "shipped"
	Cancelled OrderStatus = "cancelled"
)

func (s OrderStatus) String() string {
	return string(s)
}

func ParseOrderStatus(status string) (os OrderStatus, err error, checkStatus bool) {
	if status == "" {
		return "", nil, true
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
		return "", errors.New("invalid status"), false
	}
}

func (s OrderStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s *OrderStatus) UnmarshalJSON(data []byte) error {
	var status string
	if err := json.Unmarshal(data, &status); err != nil {
		return err
	}
	parsedStatus, err, _ := ParseOrderStatus(status)
	if err != nil {
		return err
	}
	*s = parsedStatus
	return nil
}

func (s OrderStatus) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bson.MarshalValue(s.String())
}

func (s *OrderStatus) UnmarshalBsonValue(t bsontype.Type, data []byte) error {
	var status string
	err := bson.UnmarshalValue(t, data, &status)
	if err != nil {
		return err
	}
	parsedStatus, err, _ := ParseOrderStatus(status)
	if err != nil {
		return err
	}
	*s = parsedStatus
	return nil
}
