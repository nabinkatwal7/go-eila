package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

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

func (r *Repository) GetAllCategories() ([]model.Category, error) {
	rows, err := r.DB.Query("SELECT id, name, icon, color FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []model.Category
	for rows.Next() {
		var c model.Category
		var icon, color sql.NullString
		if err := rows.Scan(&c.ID, &c.Name, &icon, &color); err != nil {
			return nil, err
		}
		if icon.Valid {
			c.Icon = icon.String
		}
		if color.Valid {
			c.Color = color.String
		}
		categories = append(categories, c)
	}
	return categories, nil
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

// --- Budgets ---

func (r *Repository) CreateBudget(b *model.Budget) error {
	query := `INSERT INTO budgets (category_id, amount, period) VALUES (?, ?, ?)`
	res, err := r.DB.Exec(query, b.CategoryID, b.Amount, b.Period)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	b.ID = id
	return nil
}

func (r *Repository) GetBudgetsWithProgress(month int, year int) ([]model.BudgetProgress, error) {
	// For each budget, calculate spent amount in that category for the given month/year.
	// We need to join splits -> transactions to filter by date.

	// Get all budgets first
	rows, err := r.DB.Query("SELECT b.category_id, c.name, b.amount FROM budgets b JOIN categories c ON b.category_id = c.id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var progress []model.BudgetProgress

	type budgetItem struct {
		CatID int64
		Name  string
		Amount int64
	}
	var items []budgetItem

	for rows.Next() {
		var bi budgetItem
		if err := rows.Scan(&bi.CatID, &bi.Name, &bi.Amount); err != nil {
			return nil, err
		}
		items = append(items, bi)
	}
	rows.Close()

	// For each, calculate spent
	// Start/End of month
	// Date string filter: start <= date < end
	// SQLite date format YYYY-MM-DD
	// Simple string match 'YYYY-MM%' works for month
	dateFilter := fmt.Sprintf("%04d-%02d%%", year, month) // e.g. 2025-01%

	for _, item := range items {
		var spentCents sql.NullInt64
		// Join splits -> transactions
		// Filter by category_id AND date
		// Amount is Debit (Positive) for Expenses.
		query := `
			SELECT SUM(s.amount)
			FROM splits s
			JOIN transactions t ON s.transaction_id = t.id
			WHERE s.category_id = ?
			AND t.date LIKE ?
			AND s.amount > 0 -- Only sum debits (expenses)
		`
		err := r.DB.QueryRow(query, item.CatID, dateFilter).Scan(&spentCents)
		if err != nil { return nil, err }

		spent := 0.0
		if spentCents.Valid {
			spent = float64(spentCents.Int64) / 100.0
		}

		budgeted := float64(item.Amount) / 100.0
		remaining := budgeted - spent
		percent := 0.0
		if budgeted > 0 {
			percent = spent / budgeted
		}

		progress = append(progress, model.BudgetProgress{
			CategoryName: item.Name,
			Budgeted:     budgeted,
			Spent:        spent,
			Remaining:    remaining,
			Percent:      percent,
		})
	}
	return progress, nil
}

// --- Recurring Logic (Simple Heuristic) ---

func (r *Repository) DetectRecurringPatterns() ([]model.Subscription, error) {
	// 1. Group transactions by Description (Payee)
	// 2. If count >= 2 and amounts are similar -> Candidate

	// This is analytical, can be heavy.
	query := `
		SELECT t.description, COUNT(*) as cnt, AVG(s.amount) as avg_amt, MAX(t.date) as last_date
		FROM transactions t
		JOIN splits s ON s.transaction_id = t.id
		WHERE s.amount > 0 -- Expenses only (debits)
		AND t.date > date('now', '-3 months') -- Look back 3 months
		GROUP BY t.description
		HAVING cnt >= 2
	`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []model.Subscription

	for rows.Next() {
		var desc string
		var cnt int
		var avgAmt float64
		var lastDate time.Time

		if err := rows.Scan(&desc, &cnt, &avgAmt, &lastDate); err != nil {
			continue
		}

		// Naive: If found, assume Monthly for now
		subs = append(subs, model.Subscription{
			Name:      desc,
			Amount:    avgAmt / 100.0,
			Frequency: "Monthly?", // Heuristic needed for real frequency
			NextDueDate: lastDate.AddDate(0, 1, 0).Format("2006-01-02"), // Assume +1 month
		})
	}

	return subs, nil
}

// --- Anomaly Detection ---

func (r *Repository) DetectAnomalies() ([]model.Anomaly, error) {
	// 1. Large Transactions (Threshold: > $500 or just > $200 for demo)
	// In production, this should be (Avg + 2*StdDev) per user.
	// We'll use a hard threshold of $200.00 (20000 cents) for now.

	query := `
		SELECT t.date, t.description, s.amount
		FROM transactions t
		JOIN splits s ON s.transaction_id = t.id
		WHERE s.amount > 20000 -- $200
		AND t.date > date('now', '-1 month') -- Recent only
		ORDER BY t.date DESC
	`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var anomalies []model.Anomaly

	for rows.Next() {
		var date time.Time
		var desc string
		var amt int64

		if err := rows.Scan(&date, &desc, &amt); err != nil {
			continue
		}

		anomalies = append(anomalies, model.Anomaly{
			Type:        "Large Transaction",
			Description: fmt.Sprintf("%s: $%.2f", desc, float64(amt)/100.0),
			Severity:    model.SeverityMedium,
			Date:        date.Format("2006-01-02"),
		})
	}

	return anomalies, nil
}

func (r *Repository) GetMonthlyStats(months int) ([]model.MonthlyStat, error) {
	// Aggregate Income vs Expense for last N months.
	// We need 12 rows (or N), with 0 if no data.
	// SQLite date grouping: strftime('%Y-%m', date)

	query := `
		SELECT strftime('%Y-%m', t.date) as month_key,
			   SUM(CASE WHEN a.type = 'Income' THEN ABS(s.amount) ELSE 0 END) as income,
			   SUM(CASE WHEN a.type = 'Expense' THEN s.amount ELSE 0 END) as expense
		FROM transactions t
		JOIN splits s ON s.transaction_id = t.id
		JOIN accounts a ON s.account_id = a.id
		WHERE t.date > date('now', ?)
		GROUP BY month_key
		ORDER BY month_key ASC
	`
	dateParam := fmt.Sprintf("-%d months", months)

	rows, err := r.DB.Query(query, dateParam)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []model.MonthlyStat

	for rows.Next() {
		var mKey string
		var inc, exp sql.NullFloat64
		if err := rows.Scan(&mKey, &inc, &exp); err != nil {
			continue
		}

		// Parse month key to "Jan"
		t, _ := time.Parse("2006-01", mKey)
		monthName := t.Format("Jan")

		stats = append(stats, model.MonthlyStat{
			Month: monthName,
			Income: inc.Float64 / 100.0,
			Expense: exp.Float64 / 100.0,
		})
	}

	return stats, nil
}

// --- Forecasting ---

type ProjectionPoint struct {
	Month string
	Value float64
}

func (r *Repository) GetNetWorthProjection(monthsAhead int) ([]ProjectionPoint, float64, error) {
	// 1. Get Current Net Worth
	stats, err := r.GetDashboardStats()
	if err != nil { return nil, 0, err }
	currentNW := stats.NetWorth

	// 2. Calculate Avg Monthly Savings (Last 3 months)
	// Query already exists essentially in GetMonthlyStats
	monthlyData, err := r.GetMonthlyStats(3)
	if err != nil { return nil, 0, err }

	var totalSavings float64
	var count float64
	for _, m := range monthlyData {
		savings := m.Income - m.Expense
		totalSavings += savings
		count++
	}

	avgSavings := 0.0
	if count > 0 {
		avgSavings = totalSavings / count
	}

	// 3. Project
	var points []ProjectionPoint
	// Start from next month
	now := time.Now()
	runningNW := currentNW

	for i := 1; i <= monthsAhead; i++ {
		runningNW += avgSavings
		futureDate := now.AddDate(0, i, 0)
		points = append(points, ProjectionPoint{
			Month: futureDate.Format("Jan 06"),
			Value: runningNW,
		})
	}

	return points, avgSavings, nil
}


func (r *Repository) CreateRule(rule *model.Rule) error {
	query := `INSERT INTO rules (pattern, target_category_id, target_payee, target_note) VALUES (?, ?, ?, ?)`
	res, err := r.DB.Exec(query, rule.Pattern, rule.TargetCategoryID, rule.TargetPayee, rule.TargetNote)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	rule.ID = id
	return nil
}

func (r *Repository) EnrichTransaction(description string) (string, *int64, string) {
	// Returns: (NormalizedPayee, CategoryID, Note)
	// Naive implementation: fetch all rules and regex match.
	// For small rule set this is fine.

	rows, err := r.DB.Query("SELECT pattern, target_category_id, target_payee, target_note FROM rules")
	if err != nil {
		return description, nil, ""
	}
	defer rows.Close()

	for rows.Next() {
		var pattern string
		var catID *int64
		var payee sql.NullString
		var note sql.NullString

		if err := rows.Scan(&pattern, &catID, &payee, &note); err != nil {
			continue
		}

		// Simple substring match for now (strings.Contains)
		// Ideally use Regex if pattern starts with ^ or similar
		if contains(description, pattern) {
			newDesc := description
			if payee.Valid && payee.String != "" {
				newDesc = payee.String
			}
			newNote := ""
			if note.Valid && note.String != "" {
				newNote = note.String
			}
			return newDesc, catID, newNote
		}
	}

	return description, nil, ""
}

// Helper
func contains(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

func (r *Repository) GetAccountByName(name string) (*model.Account, error) {
	query := `SELECT id, name, type, currency FROM accounts WHERE name = ?`
	var a model.Account
	err := r.DB.QueryRow(query, name).Scan(&a.ID, &a.Name, &a.Type, &a.Currency)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Not found
		}
		return nil, err
	}
	return &a, nil
}
