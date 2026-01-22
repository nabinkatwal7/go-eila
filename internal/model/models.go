package model

import "time"

type AccountType string
type CategoryType string
type TransactionType string

const (
	AccountTypeCash   AccountType = "Cash"
	AccountTypeBank   AccountType = "Bank"
	AccountTypeCard   AccountType = "Card"
	AccountTypeInvest AccountType = "Investment"
	AccountTypeLiability AccountType = "Liability"

	TransactionTypeExpense   TransactionType = "Expense"
	TransactionTypeIncome    TransactionType = "Income"
	TransactionTypeTransfer  TransactionType = "Transfer"
	TransactionTypeAsset     TransactionType = "Asset"     // Adjustment
	TransactionTypeLiability TransactionType = "Liability" // Adjustment
)

type Account struct {
	ID       int64
	Name     string
	Type     AccountType
	Balance  float64
	Currency string
}

type Category struct {
	ID       int64
	Name     string
	Icon     string // E.g., unicode char or resource name
	Color    string // Hex code
	ParentID *int64 // Nullable for root categories
	Type     TransactionType
}

type Transaction struct {
	ID          int64
	Amount      float64
	Date        time.Time
	Note        string
	AccountID   int64
	CategoryID  *int64 // Nullable for transfers
	TargetAccountID *int64 // Only for transfers
	Type        TransactionType
	Tags        string // Comma separated for simplicity in SQLite
}

type Budget struct {
	ID         int64
	CategoryID int64
	Amount     float64
	Period     string // e.g., "Monthly"
}
