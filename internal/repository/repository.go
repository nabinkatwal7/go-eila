package repository

import (
	"database/sql"
	"errors"
	"math"

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
	query := `INSERT INTO accounts (name, type, currency) VALUES (?, ?, ?)`
	res, err := r.DB.Exec(query, account.Name, account.Type, account.Currency)
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
	rows, err := r.DB.Query("SELECT id, name, type, currency FROM accounts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []model.Account
	for rows.Next() {
		var a model.Account
		if err := rows.Scan(&a.ID, &a.Name, &a.Type, &a.Currency); err != nil {
			return nil, err
		}
		accounts = append(accounts, a)
	}
	return accounts, nil
}

func (r *Repository) GetAccountBalance(accountID int64) (float64, error) {
	// Dynamically calculate balance from splits
	// Sum(amount) where account_id = ?
	// Since we store in cents, convert to float (div 100) for display.
	var balanceCents sql.NullInt64
	err := r.DB.QueryRow("SELECT SUM(amount) FROM splits WHERE account_id = ?", accountID).Scan(&balanceCents)
	if err != nil {
		return 0, err
	}
	if !balanceCents.Valid {
		return 0, nil
	}
	return float64(balanceCents.Int64) / 100.0, nil
}

// --- Transactions (Double Entry) ---

// CreateTransaction inserts a header and its splits transactionally.
// It explicitly checks that debits match credits (Sum of amounts == 0).
func (r *Repository) CreateTransaction(t *model.Transaction) error {
	// 1. Validate Balance
	var sum int64 = 0
	for _, s := range t.Splits {
		sum += s.Amount
	}
	if sum != 0 {
		return errors.New("transaction is not balanced (splits sum != 0)")
	}

	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 2. Insert Header
	query := `INSERT INTO transactions (date, description, note, status) VALUES (?, ?, ?, ?)`
	res, err := tx.Exec(query, t.Date, t.Description, t.Note, t.Status)
	if err != nil {
		return err
	}
	txID, err := res.LastInsertId()
	if err != nil {
		return err
	}
	t.ID = txID

	// 3. Insert Splits
	splitQuery := `INSERT INTO splits (transaction_id, account_id, category_id, amount, currency, exchange_rate) VALUES (?, ?, ?, ?, ?, ?)`
	for i := range t.Splits {
		s := &t.Splits[i]
		s.TransactionID = txID
		_, err = tx.Exec(splitQuery, s.TransactionID, s.AccountID, s.CategoryID, s.Amount, s.Currency, s.ExchangeRate)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// GetRecentTransactions fetches transactions with their splits (simple version loading only headers first).
// For a full UI, we often need the splits to know "Amount" (which is ambiguous in double entry)
// Usually we show the sum of positive splits as "Amount" or specific account impact.
func (r *Repository) GetRecentTransactions(limit int) ([]model.Transaction, error) {
	// Fetch Headers
	query := `SELECT id, date, description, note, status FROM transactions ORDER BY date DESC LIMIT ?`
	rows, err := r.DB.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []model.Transaction
	for rows.Next() {
		var t model.Transaction
		if err := rows.Scan(&t.ID, &t.Date, &t.Description, &t.Note, &t.Status); err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}

	// Fetch splits for these transactions (N+1 query for simplicity now, optimize later with JOIN)
	for i := range transactions {
		splits, err := r.GetSplitsForTransaction(transactions[i].ID)
		if err != nil {
			return nil, err
		}
		transactions[i].Splits = splits
	}

	return transactions, nil
}

func (r *Repository) GetSplitsForTransaction(txID int64) ([]model.Split, error) {
	query := `SELECT id, transaction_id, account_id, category_id, amount, currency, exchange_rate FROM splits WHERE transaction_id = ?`
	rows, err := r.DB.Query(query, txID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var splits []model.Split
	for rows.Next() {
		var s model.Split
		if err := rows.Scan(&s.ID, &s.TransactionID, &s.AccountID, &s.CategoryID, &s.Amount, &s.Currency, &s.ExchangeRate); err != nil {
			return nil, err
		}
		splits = append(splits, s)
	}
	return splits, nil
}

// Stats (Updated for Double Entry)

type DashboardStats struct {
	TotalIncome    float64
	TotalExpense   float64
	TotalAssets    float64
	TotalLiability float64
	NetWorth       float64
}

func (r *Repository) GetDashboardStats() (*DashboardStats, error) {
	stats := &DashboardStats{}

	// We need to join splits with accounts to check account type
	// Income: Sum of splits where Account Type = Income. (Usually negative credit, so abs or invert)
	// Expense: Sum of splits where Account Type = Expense. (Positive debit)

	// Helper query
	sumByType := func(accType model.AccountType) (float64, error) {
		var val sql.NullFloat64
		// SUM(amount * exchange_rate)
		// Note: Splits amount is Integer (cents). Exchange Rate is Real.
		// Result is Real float representing total cents in Base Currency.
		query := `
			SELECT SUM(s.amount * s.exchange_rate)
			FROM splits s
			JOIN accounts a ON s.account_id = a.id
			WHERE a.type = ?
		`
		err := r.DB.QueryRow(query, accType).Scan(&val)
		if err != nil {
			return 0, err
		}
		if !val.Valid { return 0, nil }
		return val.Float64 / 100.0, nil
	}

	inc, err := sumByType(model.AccountTypeIncome)
	if err != nil { return nil, err }
	stats.TotalIncome = math.Abs(inc) // Display positive

	exp, err := sumByType(model.AccountTypeExpense)
	if err != nil { return nil, err }
	stats.TotalExpense = exp

	// Assets (Sum of Cash, Bank, Invest)
	// We can't use the simple helper for multiple types easily without modification.
	// Let's do a manual query for assets.
	var assetVal sql.NullFloat64
	assetQuery := `
		SELECT SUM(s.amount * s.exchange_rate)
		FROM splits s
		JOIN accounts a ON s.account_id = a.id
		WHERE a.type IN (?, ?, ?)
	`
	err = r.DB.QueryRow(assetQuery, model.AccountTypeCash, model.AccountTypeBank, model.AccountTypeInvest).Scan(&assetVal)
	if err != nil { return nil, err }
	stats.TotalAssets = 0
	if assetVal.Valid {
		stats.TotalAssets = assetVal.Float64 / 100.0
	}

	liab, err := sumByType(model.AccountTypeLiability)
	if err != nil { return nil, err }
	stats.TotalLiability = math.Abs(liab) // Usually negative balance implies debt, but standard liability is Credit normal.

	stats.NetWorth = stats.TotalAssets - stats.TotalLiability

	return stats, nil
}
