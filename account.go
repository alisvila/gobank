package main

import (
	"math/rand"
	"time"
)

type TransferRequest struct {
	ToNumber int64
	Amount   int
}
type AccountAPI struct {
	Firstname string
	Lastname  string
}
type Account struct {
	ID                int
	Firstname         string
	Lastname          string
	Number            int64
	Balance           int64
	EncryptedPassword string
	createdAt         time.Time
}

func makeAccount(firstname, lastname string) *Account {
	return &Account{
		ID:        rand.Intn(1000000),
		Firstname: firstname,
		Lastname:  lastname,
		Number:    int64(rand.Intn(100000)),
	}
}
