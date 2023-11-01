package main

import (
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type TransferRequest struct {
	ToNumber int64
	Amount   int
}

type LoginRequest struct {
	Username string
	Password string
}

type LoginResponse struct {
	Token  string
	UserID string
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
	token             string
}

func makeAccount(firstname, lastname string) *Account {
	return &Account{
		ID:        rand.Intn(1000000),
		Firstname: firstname,
		Lastname:  lastname,
		Number:    int64(rand.Intn(100000)),
	}
}

func (a *Account) validPassword(pwd string) bool {
	return bcrypt.CompareHashAndPassword([]byte(a.EncryptedPassword), []byte(pwd)) == nil
}
