package main

import (
	"github.com/google/uuid"
	"math/rand"
)

type Account struct {
	ID         string
	FirstName  string
	LastName   string
	Email      string
	ClientCode int
	Balance    float32
}

func NewAccount(
	firstName string,
	lastName string,
	email string,
) *Account {
	return &Account{
		ID:         uuid.New().String(),
		FirstName:  firstName,
		LastName:   lastName,
		Email:      email,
		ClientCode: rand.Intn(100),
		Balance:    0.0,
	}
}
