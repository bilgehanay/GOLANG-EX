package service

import (
	"deneme.com/bng-go/Model"
	"github.com/google/uuid"
)

type UserService interface {
	CreateUser(*Model.User) error
	GetUser(*uuid.UUID) (*Model.User, error)
	GetUsers() ([]*Model.User, error)
	UpdateUser(*Model.User) error
	DeleteUser(*uuid.UUID) error
}
