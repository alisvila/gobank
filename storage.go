package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
	"gopkg.in/yaml.v3"
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

type Config struct {
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Dbname   string `yaml:"dbname"`
	}
}

func newPostgress() (*PostGressStorage, error) {
	f, err := os.ReadFile("config.yaml")
	if err != nil {
		return nil, err
	}

	cfg := Config{}
	err = yaml.Unmarshal(f, &cfg)
	if err != nil {
		return nil, err
	}
	pgsql := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable", cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password)
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
		fmt.Print(err)
	}
	fmt.Println(res)

	return &PostGressStorage{
		db: db,
	}, nil
}

func (s *PostGressStorage) runQueryWithCtx(query string) (sql.Result, error) {
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	return s.db.ExecContext(ctx, query)
}

func (s *PostGressStorage) CreateAccountTable() error {
	query := `create table if not exists accounts (
		id serial primary key,
		firsname varchar(50),
		lastname varchar(50),
		number serial,
		balance serial,
		created_at timestamp
	)`
	_, err := s.runQueryWithCtx(query)
	return err
}

func (s *PostGressStorage) CreateAccount(acc *Account) error {
	query := `insert into account
	(firstname, lastname, number, balance, craeted_at)
	values ($1,$2,$3,$4,$5,$6)`
	_, err := s.db.Query(query, acc.Firstname, acc.Lastname, acc.Number, acc.Balance, acc.createdAt)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostGressStorage) DeleteAccount(id int) error {
	query := "delete from account where id = $1"
	_, err := s.db.Query(query, id)
	if err != nil {
		return err
	}
	return nil
}
func (s *PostGressStorage) GetAccount(accNumber int) (*Account, error) {
	query := "select * from account where number = $1"
	rows, err := s.db.Query(query, accNumber)
	if err != nil {
		return nil, err
	}

	account := &Account{}
	for rows.Next() {
		err := rows.Scan(&account.Firstname, &account.Lastname, &account.Balance, &account.Number, &account.createdAt)
		if err != nil {
			return nil, err
		}
	}
	return account, nil
}
func (s *PostGressStorage) GetAllAccounts() ([]*Account, error) {
	query := "select * from accounts"
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	accounts := []*Account{}
	for rows.Next() {
		account := &Account{}
		if err := rows.Scan(&account.Firstname, &account.Lastname, &account.Balance, &account.Number, &account.createdAt); err != nil {
			return accounts, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}
