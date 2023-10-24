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

func errorHandler() {
	if r := recover(); r != nil {
		fmt.Println(r)
	}
}

func newPostgress() (*PostGressStorage, error) {
	defer errorHandler()

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
