package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage interface {
	GetAccountByUUID(string) (*Account, error)
	FindAllAccounts() ([]*Account, error)
	CreateAccount(*Account) (*Account, error)
	DeleteAccount(string) error
	UpdateAccount(*Account) (*Account, error)
}

type PostgresStore struct {
	db *sql.DB
}

func rowToAccount(rows *sql.Rows) (*Account, error) {
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
	return account, nil
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
		account, err := rowToAccount(rows)
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

	for rows.Next() {
		return rowToAccount(rows)
	}
	return nil, fmt.Errorf("not found")
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
  ) RETURNING *
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

func (store *PostgresStore) UpdateAccount(account *Account) (*Account, error) {
	updateQuery := `UPDATE account SET
    first_name = $1,
    last_name = $2,
    email = $3
    WHERE uuid = $4
    RETURNING *
  ;`

	newAcc := new(Account)
	err := store.db.QueryRow(
		updateQuery,
		account.FirstName,
		account.LastName,
		account.Email,
		account.UUID,
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

func (store *PostgresStore) DeleteAccount(uuid string) error {
	deleteQuery := `DELETE FROM account WHERE uuid = $1;`

	_, err := store.db.Exec(
		deleteQuery,
		uuid,
	)
	return err
}
