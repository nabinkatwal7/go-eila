package model

import "time"

type AccountType string
type TransactionStatus string

const (
	AccountTypeCash      AccountType = "Cash"
	AccountTypeBank      AccountType = "Bank"
	AccountTypeCard      AccountType = "Card"
	AccountTypeInvest    AccountType = "Investment"
	AccountTypeEquity    AccountType = "Equity"
	AccountTypeLiability AccountType = "Liability"
	AccountTypeIncome    AccountType = "Income"
	AccountTypeExpense   AccountType = "Expense"

	TransactionStatusPending    TransactionStatus = "Pending"
	TransactionStatusCleared    TransactionStatus = "Cleared"
	TransactionStatusReconciled TransactionStatus = "Reconciled"
)

type Account struct {
	ID       int64
	Name     string
	Type     AccountType
	Currency string
}

type Category struct {
	ID       int64
	Name     string
	Icon     string
	Color    string
	ParentID *int64
}

type Transaction struct {
	ID          int64
	Date        time.Time
	Description string
	Note        string
	Status      TransactionStatus
	Splits      []Split
}

type Split struct {
	ID            int64
	TransactionID int64
	AccountID     int64
	CategoryID    *int64

	Amount int64

	Currency      string
	ExchangeRate  float64
}

type Budget struct {
	ID         int64
	CategoryID int64
	Amount     int64
	Period     string
}

type BudgetProgress struct {
	CategoryName string
	Budgeted     float64
	Spent        float64
	Remaining    float64
	Percent      float64
}

// Rule defines smart enrichment logic
type Rule struct {
	ID             int64
	Pattern        string // Regex or Simple match
	TargetCategoryID *int64
	TargetPayee    string
	TargetNote     string
}

// CategoryBreakdown represents spending by category
type CategoryBreakdown struct {
	CategoryID   int64
	CategoryName string
	Color        string
	Amount       float64
}

// NetWorthPoint represents net worth at a point in time
type NetWorthPoint struct {
	Month      string
	Assets     float64
	Liabilities float64
	NetWorth   float64
}