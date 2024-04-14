package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage interface {
	GetAccountByUUID(string) (*Account, error)
	FindAllAccounts() ([]*Account, error)
	CreateAccount(*Account) error
	DeleteAccount(string) error
	UpdateAccount(*Account) error
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connectionString := "user=postgres dbname=postgres password=goBank123. sslmode=disable"
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (store *PostgresStore) Init() error {
	store.CreateAccountsTable()
	return store.CreatePGCryptoExtension()
}

func (store *PostgresStore) CreatePGCryptoExtension() error {
	query := `CREATE EXTENSION IF NOT EXISTS "pgcrypto";`
	_, err := store.db.Exec(query)
	return err
}

func (store *PostgresStore) CreateAccountsTable() error {
	query := `CREATE TABLE IF NOT EXISTS account (
		uuid UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
		firstName VARCHAR(255) NOT NULL,
		lastName VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL,
		clientCode SERIAL,
		balance DECIMAL(16, 2) NOT NULL DEFAULT 0.0,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`

	_, err := store.db.Exec(query)
	return err
}

func (store *PostgresStore) FindAllAccounts() ([]*Account, error) {
	rows, err := store.db.Query("SELECT * FROM account;")
	if err != nil {
		return nil, err
	}

	accountList := []*Account{}
	fmt.Printf("Rows: %+v\n", rows)
	for rows.Next() {
		account := &Account{}
		err := rows.Scan(
			&account.ID,
			&account.FirstName,
			&account.LastName,
			&account.Email,
			&account.ClientCode,
			&account.Balance,
			&account.CreatedAt,
			&account.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		accountList = append(accountList, account)
	}
	return accountList, nil
}

func (store *PostgresStore) GetAccountByUUID(uuid string) (*Account, error) {
	return nil, nil
}

func (store *PostgresStore) CreateAccount(account *Account) error {
	insertionQuery := `INSERT INTO account (
		firstName,
		lastName,
		email
	) VALUES (
		$1,
		$2,
		$3
	);`
	result, err := store.db.Exec(
		insertionQuery,
		account.FirstName,
		account.LastName,
		account.Email,
	)
	if err != nil {
		return err
	}

	fmt.Printf("Result: %+v\n", result)

	return nil
}

func (store *PostgresStore) UpdateAccount(*Account) error {
	return nil
}

func (store *PostgresStore) DeleteAccount(uuid string) error {
	return nil
}
