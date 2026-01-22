package repository

import (
	"github.com/nabinkatwal7/go-eila/internal/model"
)

type Repository struct {
	DB *DB
}

func NewRepository(db *DB) *Repository {
	return &Repository{DB: db}
}

// --- Accounts ---

func (r *Repository) CreateAccount(account *model.Account) error {
	query := `INSERT INTO accounts (name, type, balance, currency) VALUES (?, ?, ?, ?)`
	res, err := r.DB.Exec(query, account.Name, account.Type, account.Balance, account.Currency)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	account.ID = id
	return nil
}

func (r *Repository) GetAllAccounts() ([]model.Account, error) {
	rows, err := r.DB.Query("SELECT id, name, type, balance, currency FROM accounts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []model.Account
	for rows.Next() {
		var a model.Account
		if err := rows.Scan(&a.ID, &a.Name, &a.Type, &a.Balance, &a.Currency); err != nil {
			return nil, err
		}
		accounts = append(accounts, a)
	}
	return accounts, nil
}

// --- Categories ---

func (r *Repository) CreateCategory(category *model.Category) error {
	query := `INSERT INTO categories (name, icon, color, parent_id, type) VALUES (?, ?, ?, ?, ?)`
	res, err := r.DB.Exec(query, category.Name, category.Icon, category.Color, category.ParentID, category.Type)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	category.ID = id
	return nil
}

func (r *Repository) GetAllCategories() ([]model.Category, error) {
	rows, err := r.DB.Query("SELECT id, name, icon, color, parent_id, type FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []model.Category
	for rows.Next() {
		var c model.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Icon, &c.Color, &c.ParentID, &c.Type); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

// --- Transactions ---

func (r *Repository) CreateTransaction(t *model.Transaction) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Insert Transaction
	query := `INSERT INTO transactions (amount, date, note, account_id, category_id, target_account_id, type, tags) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	res, err := tx.Exec(query, t.Amount, t.Date, t.Note, t.AccountID, t.CategoryID, t.TargetAccountID, t.Type, t.Tags)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	t.ID = id

	// 2. Update Account Balance
	// Income adds to balance, Expense subtracts. Setup for Transfer/Assets later.
	var updateAccountQuery string
	if t.Type == model.TransactionTypeIncome {
		updateAccountQuery = `UPDATE accounts SET balance = balance + ? WHERE id = ?`
		_, err = tx.Exec(updateAccountQuery, t.Amount, t.AccountID)
	} else if t.Type == model.TransactionTypeExpense {
		updateAccountQuery = `UPDATE accounts SET balance = balance - ? WHERE id = ?`
		_, err = tx.Exec(updateAccountQuery, t.Amount, t.AccountID)
	}
	// TODO: Handle transfers and other types

	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *Repository) GetRecentTransactions(limit int) ([]model.Transaction, error) {
	query := `SELECT id, amount, date, note, account_id, category_id, target_account_id, type, tags FROM transactions ORDER BY date DESC LIMIT ?`
	rows, err := r.DB.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []model.Transaction
	for rows.Next() {
		var t model.Transaction
		if err := rows.Scan(&t.ID, &t.Amount, &t.Date, &t.Note, &t.AccountID, &t.CategoryID, &t.TargetAccountID, &t.Type, &t.Tags); err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}
	return transactions, nil
}

// --- Stats ---

type DashboardStats struct {
	TotalIncome   float64
	TotalExpense  float64
	TotalAssets   float64
	TotalLiability float64
	NetWorth      float64
}

func (r *Repository) GetDashboardStats() (*DashboardStats, error) {
	stats := &DashboardStats{}

	// simple sum queries
	// Optimally, do this in one query or cache it.

	// Income
	row := r.DB.QueryRow("SELECT COALESCE(SUM(amount), 0) FROM transactions WHERE type = ?", model.TransactionTypeIncome)
	if err := row.Scan(&stats.TotalIncome); err != nil {
		return nil, err
	}

	// Expense
	row = r.DB.QueryRow("SELECT COALESCE(SUM(amount), 0) FROM transactions WHERE type = ?", model.TransactionTypeExpense)
	if err := row.Scan(&stats.TotalExpense); err != nil {
		return nil, err
	}

	// Assets/Liabilities should come from Accounts balance sum based on account type generally,
	// or Transaction types if we are strictly following EILA 4 pillars as "flows".
	// But usually Assets = Sum of Asset Accounts.
	// For this implementation, let's assume Assets = Sum(Bank + Cash + Investment)
	// Liabilities = Sum(Credit Cards + Loans) (which would be negative balance or positive liability?)
	// Let's stick to simple Account Type sums for Assets/Liabilities

	// Assets
	// query types: Cash, Bank, Investment
	row = r.DB.QueryRow("SELECT COALESCE(SUM(balance), 0) FROM accounts WHERE type IN (?, ?, ?)",
		model.AccountTypeCash, model.AccountTypeBank, model.AccountTypeInvest)
	if err := row.Scan(&stats.TotalAssets); err != nil {
		return nil, err
	}

	// Liabilities
	row = r.DB.QueryRow("SELECT COALESCE(SUM(balance), 0) FROM accounts WHERE type = ?", model.AccountTypeLiability)
	if err := row.Scan(&stats.TotalLiability); err != nil {
		return nil, err
	}

	stats.NetWorth = stats.TotalAssets - stats.TotalLiability

	return stats, nil
}
