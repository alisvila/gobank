package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

const (
	host     = "38.54.13.174"
	port     = 5432
	user     = "postgres"
	password = "superstrongpass"
	dbname   = ""
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	GetAccount(int) (*Account, error)
	GetAllAccounts() ([]*Account, error)
}

type PostGressStorage struct {
	db *sql.DB
}

func newPostgress() (*PostGressStorage, error) {
	pgsql := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable", host, port, user, password)
	fmt.Println(pgsql)
	db, err := sql.Open("postgres", pgsql)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	query := "CREATE DATABASE gobank"
	res, err := db.ExecContext(ctx, query)
	if err != nil {
		return nil, err
	}
	fmt.Println(res)

	return &PostGressStorage{
		db: db,
	}, nil
}

func (s *PostGressStorage) CreateAccountTable(db *PostGressStorage) error {
	query := `create table if not exists accounts (
		id serial primary key,
		firsname varchar(50),
		lastname varchar(50),
		number serial,
		balance serial,
		created_at timestamp
	)`
	_, err := s.db.Exec(query)
	return err
}

func (s *PostGressStorage) CreateAccount(*Account) error {
	return nil
}

func (s *PostGressStorage) DeleteAccount(int) error {
	return nil
}
func (s *PostGressStorage) GetAccount(int) (*Account, error) {
	return nil, nil
}
func (s *PostGressStorage) GetAllAccounts() ([]*Account, error) {
	// query := "select * from accounts"
	// res, err :=
	return nil, nil
}
