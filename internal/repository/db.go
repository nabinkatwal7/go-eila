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

	// For dev, strict foreign keys
	if _, err := db.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		return nil, err
	}

	if err := createSchema(db); err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

func createSchema(db *sql.DB) error {
	// Note: In a production app, we would use proper migrations (golang-migrate etc).
	// For this refactor, we are redefining the schema.
	// Tables: accounts, categories, transactions (header), splits.

	schema := `
	CREATE TABLE IF NOT EXISTS accounts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		type TEXT NOT NULL,
		currency TEXT DEFAULT 'USD',
		is_closed BOOLEAN DEFAULT 0
	);

	CREATE TABLE IF NOT EXISTS categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		icon TEXT,
		color TEXT,
		parent_id INTEGER,
		FOREIGN KEY(parent_id) REFERENCES categories(id)
	);

	CREATE TABLE IF NOT EXISTS transactions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date DATETIME NOT NULL,
		description TEXT NOT NULL,
		note TEXT,
		status TEXT DEFAULT 'Pending',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS splits (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		transaction_id INTEGER NOT NULL,
		account_id INTEGER NOT NULL,
		category_id INTEGER,
		amount INTEGER NOT NULL, -- Stored in minor units (cents)
		currency TEXT DEFAULT 'USD',
		exchange_rate REAL DEFAULT 1.0,
		FOREIGN KEY(transaction_id) REFERENCES transactions(id) ON DELETE CASCADE,
		FOREIGN KEY(account_id) REFERENCES accounts(id),
		FOREIGN KEY(category_id) REFERENCES categories(id)
	);
	`
	// Note: We might need to Drop tables if they exist with old schema,
	// but for now relying on user starting fresh or manual cleanup since this is a dev phase.

	_, err := db.Exec(schema)
	if err != nil {
		return err
	}

	go seedDefaults(db)

	return nil
}

func seedDefaults(db *sql.DB) {
	var count int
	row := db.QueryRow("SELECT COUNT(*) FROM categories")
	if err := row.Scan(&count); err == nil && count == 0 {
		log.Println("Seeding default categories...")
		// TODO: Seed logic
	}
}
