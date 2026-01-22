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
	Amount     int64 // Minor units (cents)
	Period     string // "Monthly"
	// Optional: Currency (assuming Base for now, or per budget)
}

// Progress struct for UI
type BudgetProgress struct {
	CategoryName string
	Budgeted     float64
	Spent        float64
	Remaining    float64
	Percent      float64
}
