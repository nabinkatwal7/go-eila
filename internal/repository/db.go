package repository

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

type DB struct {
	*sql.DB
}

func NewDB(dbPath string) (*DB, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if err := createSchema(db); err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

func createSchema(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS accounts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		type TEXT NOT NULL,
		balance REAL DEFAULT 0,
		currency TEXT DEFAULT 'USD'
	);

	CREATE TABLE IF NOT EXISTS categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		icon TEXT,
		color TEXT,
		parent_id INTEGER,
		type TEXT,
		FOREIGN KEY(parent_id) REFERENCES categories(id)
	);

	CREATE TABLE IF NOT EXISTS transactions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		amount REAL NOT NULL,
		date DATETIME NOT NULL,
		note TEXT,
		account_id INTEGER NOT NULL,
		category_id INTEGER,
		target_account_id INTEGER,
		type TEXT NOT NULL,
		tags TEXT,
		FOREIGN KEY(account_id) REFERENCES accounts(id),
		FOREIGN KEY(category_id) REFERENCES categories(id),
		FOREIGN KEY(target_account_id) REFERENCES accounts(id)
	);

	CREATE TABLE IF NOT EXISTS budgets (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		category_id INTEGER NOT NULL,
		amount REAL NOT NULL,
		period TEXT DEFAULT 'Monthly',
		FOREIGN KEY(category_id) REFERENCES categories(id)
	);
	`

	_, err := db.Exec(schema)
	if err != nil {
		return err
	}

	// Seed default categories if empty
	go seedDefaults(db)

	return nil
}

func seedDefaults(db *sql.DB) {
	var count int
	row := db.QueryRow("SELECT COUNT(*) FROM categories")
	if err := row.Scan(&count); err == nil && count == 0 {
		// Basic seeding
		log.Println("Seeding default categories...")
		// TODO: Add comprehensive seed data
	}
}
