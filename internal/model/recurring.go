package model

// ... existing models ...

// RecurringRule defines a detected or manual recurring transaction
type RecurringRule struct {
	ID        int64
	Pattern   string // Payee name pattern
	Amount    int64  // Expected amount (cents)
	Interval  string // "Monthly", "Weekly"
	LastDate  string // YYYY-MM-DD
	Confirmed bool   // User confirmed this rule?
}

// Subscription is a high-level view of a recurring expense
type Subscription struct {
	Name        string
	Amount      float64
	Frequency   string
	NextDueDate string
}
