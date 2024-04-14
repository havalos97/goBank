package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Storage interface {
	GetAccountByUUID(string) (*Account, error)
	FindAllAccounts() ([]*Account, error)
	CreateAccount(*Account) (*Account, error)
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
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    client_code SERIAL,
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
	for rows.Next() {
		account := new(Account)
		err := rows.Scan(
			&account.UUID,
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
	rows, err := store.db.Query("SELECT * FROM account WHERE uuid=$1 LIMIT 1;", uuid)
	if err != nil {
		return nil, err
	}

	foundAcc := new(Account)
	for rows.Next() {
		account := new(Account)
		err := rows.Scan(
			&account.UUID,
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
		foundAcc = account
	}
	return foundAcc, nil
}

func (store *PostgresStore) CreateAccount(account *Account) (*Account, error) {
	insertionQuery := `INSERT INTO account (
    uuid,
    first_name,
    last_name,
    email
  ) VALUES (
    $1,
    $2,
    $3,
    $4
  ) RETURNING
    uuid,
    first_name,
    last_name,
    email,
    client_code,
    balance,
    created_at,
    updated_at
  ;`

	newAcc := new(Account)
	err := store.db.QueryRow(
		insertionQuery,
		account.UUID,
		account.FirstName,
		account.LastName,
		account.Email,
	).Scan(
		&newAcc.UUID,
		&newAcc.FirstName,
		&newAcc.LastName,
		&newAcc.Email,
		&newAcc.ClientCode,
		&newAcc.Balance,
		&newAcc.CreatedAt,
		&newAcc.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (store *PostgresStore) UpdateAccount(*Account) error {
	return nil
}

func (store *PostgresStore) DeleteAccount(uuid string) error {
	return nil
}
