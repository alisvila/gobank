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
	// defer db.Close()

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

	err = CreateAccountTable(db)
	if err != nil {
		return nil, err
	}
	return &PostGressStorage{
		db: db,
	}, nil
}

func runQueryWithCtx(query string, db *sql.DB) (sql.Result, error) {
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	return db.ExecContext(ctx, query)
}

func CreateAccountTable(db *sql.DB) error {
	query := `create table if not exists account (
		id serial primary key,
		first_name varchar(100),
		last_name varchar(100),
		number serial,
		encrypted_password varchar(100),
		balance serial,
		created_at timestamp
	)`
	_, err := runQueryWithCtx(query, db)
	return err
}

func (s *PostGressStorage) CreateAccount(acc *Account) error {
	query := `insert into account 
	(first_name, last_name, number, balance, created_at)
	values ($1, $2, $3, $4, $5)`

	fmt.Printf("the calu %+v", acc)
	_, err := s.db.Query(query, acc.Firstname, acc.Lastname, acc.Number, acc.Balance, acc.createdAt)
	fmt.Println(err)
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
		err := rows.Scan(&account.ID, &account.Firstname, &account.Lastname, &account.Balance, &account.Number, &account.EncryptedPassword, &account.createdAt)
		if err != nil {
			return nil, err
		}
	}
	return account, nil
}
func (s *PostGressStorage) GetAllAccounts() ([]*Account, error) {
	query := "select * from account"
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	accounts := []*Account{}
	for rows.Next() {
		account := &Account{}
		if err := rows.Scan(&account.ID, &account.Firstname, &account.Lastname, &account.Balance, &account.EncryptedPassword, &account.Number, &account.createdAt); err != nil {
			return accounts, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}
