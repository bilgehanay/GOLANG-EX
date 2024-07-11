package Model

import "github.com/google/uuid"

type Address struct {
	Street  string `json:"street" bson:"street,omitempty"`
	State   string `json:"state" bson:"state,omitempty"`
	City    string `json:"city" bson:"city,omitempty"`
	ZipCode string `json:"zip_code" bson:"zip_code,omitempty"`
	Country string `json:"country" bson:"country,omitempty"`
}

type User struct {
	Id       uuid.UUID `json:"id" bson:"_id"`
	Name     string    `json:"name" bson:"name" validate:"required,min=2,max=32"`
	Surname  string    `json:"surname" bson:"surname" validate:"required,min=2,max=32"`
	Email    string    `json:"email" bson:"email" validate:"required,email"`
	Password string    `json:"password" bson:"password" validate:"required,min=6,max=32"`
	Age      int       `json:"age" bson:"age" validate:"required,min=18,max=120"`
	Address  Address   `json:"address" bson:"address,omitempty" validate:"required"`
}

func NewUser(name, surname, email, password string, age int, address Address) User {
	return User{
		Id:       uuid.New(),
		Name:     name,
		Surname:  surname,
		Email:    email,
		Password: password,
		Age:      age,
		Address:  address,
	}
}
