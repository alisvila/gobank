package main

import "math/rand"

type Account struct {
	ID        int
	Firstname string
	Lastname  string
	Number    int64
	Balance   int64
}

func makeAccount(firstname, lastname string) *Account {
	return &Account{
		ID:        rand.Intn(1000000),
		Firstname: firstname,
		Lastname:  lastname,
		Number:    rand.Int63(),
	}
}
