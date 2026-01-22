package model

import "time"

type AccountType string
type TransactionStatus string

const (
	AccountTypeCash      AccountType = "Cash"
	AccountTypeBank      AccountType = "Bank"
	AccountTypeCard      AccountType = "Card"
	AccountTypeInvest    AccountType = "Investment"
	AccountTypeEquity    AccountType = "Equity" // Added in double entry
	AccountTypeLiability AccountType = "Liability"
	AccountTypeIncome    AccountType = "Income"
	AccountTypeExpense   AccountType = "Expense"

	// Sub-types can be handled via tags or a separate field if needed,
	// for now mapping high-level types.

	TransactionStatusPending    TransactionStatus = "Pending"
	TransactionStatusCleared    TransactionStatus = "Cleared"
	TransactionStatusReconciled TransactionStatus = "Reconciled"
)

type Account struct {
	ID       int64
	Name     string
	Type     AccountType
	Currency string // ISO 4217 e.g. "USD", "EUR"
	// Balance is now calculated dynamically from splits,
	// but caching it might be useful later.
	// For strict double-entry, we query sum(splits).
}

type Category struct {
	ID       int64
	Name     string
	Icon     string
	Color    string
	ParentID *int64
}

// Transaction represents the "Header" of a double-entry event.
type Transaction struct {
	ID          int64
	Date        time.Time
	Description string // Payee or short desc
	Note        string // Detailed memo
	Status      TransactionStatus
	Splits      []Split // Loaded on demand usually, but helpful here
}

// Split represents one leg of the transaction.
// Sum of Amount must be 0 for a valid transaction.
type Split struct {
	ID            int64
	TransactionID int64
	AccountID     int64
	CategoryID    *int64 // Optional: mainly for Expense/Income accounts

	// Amount in minor units (e.g. cents) to avoid float errors.
	// Positive = Increase Asset/Expense (Debit)
	// Negative = Increase Liability/Income (Credit)
	// PROVISO: This sign convention varies.
	// Let's use:
	// Debit is Positive (+), Credit is Negative (-)
	// Asset (+) -> Debit increases
	// Liability (-) -> Credit increases
	// Income (-) -> Credit increases
	// Expense (+) -> Debit increases
	Amount int64

	Currency      string  // Currency of this split
	ExchangeRate  float64 // 1.0 if same as base.
}

// Helper for UI "Simple Mode"
type SimpleTransactionInput struct {
	Date        time.Time
	Description string
	Amount      float64 // Input as float, convert to int64
	Type        string // Expense, Income, Transfer
	FromAccountID int64
	ToAccountID   *int64 // For transfers
	CategoryID    *int64
}
