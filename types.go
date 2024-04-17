package main

import (
	"time"

	"github.com/google/uuid"
)

type UpsertAccountRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

type Account struct {
	UUID       string    `json:"uuid"`
	FirstName  string    `json:"firstName"`
	LastName   string    `json:"lastName"`
	Email      string    `json:"email"`
	ClientCode int       `json:"clientCode"`
	Balance    float32   `json:"balance"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

func NewAccount(
	firstName string,
	lastName string,
	email string,
) *Account {
	return &Account{
		UUID:       uuid.New().String(),
		FirstName:  firstName,
		LastName:   lastName,
		Email:      email,
		ClientCode: 0,
		Balance:    0.0,
	}
}
