# MyTrack API Documentation

This document outlines the core Go API provided by the `internal/repository` package for interacting with the MyTrack database.

## Repository

The `Repository` struct is the main entry point for all database operations.

```go
type Repository struct {
    DB *DB
}

func NewRepository(db *DB) *Repository
```

## Methods

### Accounts

#### `GetAllAccounts`

Retrieves all accounts from the database.

```go
func (r *Repository) GetAllAccounts() ([]model.Account, error)
```

**Returns:** A slice of all accounts, or an error if the query fails.

**Example:**

```go
accounts, err := repo.GetAllAccounts()
if err != nil {
    log.Fatal(err)
}
for _, acc := range accounts {
    fmt.Printf("%s (%s)\n", acc.Name, acc.Type)
}
```

#### `CreateAccount`

Creates a new account. Populates the `ID` field of the passed struct on success.

```go
func (r *Repository) CreateAccount(account *model.Account) error
```

**Parameters:**

- `account`: Pointer to an `Account` struct with `Name`, `Type`, and `Currency` fields set.

**Returns:** Error if the insert fails, otherwise `nil`. The `account.ID` field is populated on success.

**Example:**

```go
account := &model.Account{
    Name:     "Chase Checking",
    Type:     model.AccountTypeBank,
    Currency: "USD",
}
err := repo.CreateAccount(account)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Created account with ID: %d\n", account.ID)
```

#### `GetAccountByName`

Checks if an account exists with the given name. Returns `nil` (and no error) if not found.

```go
func (r *Repository) GetAccountByName(name string) (*model.Account, error)
```

**Parameters:**

- `name`: The account name to search for.

**Returns:** Pointer to the account if found, `nil` if not found, or an error if the query fails.

#### `GetAccountBalance`

Calculates the current balance of an account by summing all splits.

```go
func (r *Repository) GetAccountBalance(accountID int64) (float64, error)
```

**Parameters:**

- `accountID`: The ID of the account.

**Returns:** The account balance as a float64 (in dollars), or an error.

**Note:** Balances are calculated dynamically from splits. Positive amounts increase the balance, negative amounts decrease it.

### Categories

#### `GetAllCategories`

Retrieves all categories from the database.

```go
func (r *Repository) GetAllCategories() ([]model.Category, error)
```

**Returns:** A slice of all categories, or an error if the query fails.

**Example:**

```go
categories, err := repo.GetAllCategories()
if err != nil {
    log.Fatal(err)
}
for _, cat := range categories {
    fmt.Printf("%s (Color: %s)\n", cat.Name, cat.Color)
}
```

### Transactions

#### `CreateTransaction`

Creates a new balanced transaction with multiple splits. This method enforces the double-entry accounting invariant: the sum of all splits must equal zero.

```go
func (r *Repository) CreateTransaction(t *model.Transaction) error
```

**Parameters:**

- `t`: Pointer to a `Transaction` struct with `Date`, `Description`, `Note`, `Status`, and `Splits` populated.

**Returns:** Error if validation fails (unbalanced transaction) or if the database operation fails. On success, `t.ID` is populated.

**Validation:**

- The sum of all split amounts must equal zero.
- The transaction is inserted atomically (all splits are inserted in a single database transaction).

**Example:**

```go
transaction := &model.Transaction{
    Date:        time.Now(),
    Description: "Grocery Shopping",
    Note:        "Weekly groceries",
    Status:      model.TransactionStatusCleared,
    Splits: []model.Split{
        {
            AccountID:  cashAccountID,
            Amount:     -5000, // -$50.00 (credit from cash)
        },
        {
            AccountID:  expenseAccountID,
            CategoryID: &foodCategoryID,
            Amount:     5000, // +$50.00 (debit to expense)
        },
    },
}
err := repo.CreateTransaction(transaction)
if err != nil {
    log.Fatal(err)
}
```

#### `GetRecentTransactions`

Fetches the most recent transactions with their associated splits.

```go
func (r *Repository) GetRecentTransactions(limit int) ([]model.Transaction, error)
```

**Parameters:**

- `limit`: Maximum number of transactions to retrieve.

**Returns:** A slice of transactions ordered by date (most recent first), each with its splits populated.

#### `GetSplitsForTransaction`

Retrieves all splits for a specific transaction.

```go
func (r *Repository) GetSplitsForTransaction(txID int64) ([]model.Split, error)
```

**Parameters:**

- `txID`: The transaction ID.

**Returns:** A slice of splits for the transaction, or an error.

### Statistics & Analytics

#### `GetDashboardStats`

Calculates comprehensive dashboard statistics including income, expenses, assets, liabilities, and net worth.

```go
func (r *Repository) GetDashboardStats() (*DashboardStats, error)

type DashboardStats struct {
    TotalIncome    float64
    TotalExpense   float64
    TotalAssets    float64
    TotalLiability float64
    NetWorth       float64
}
```

**Returns:** A pointer to `DashboardStats` with all calculated values, or an error.

**Note:** Statistics are calculated by joining splits with accounts and aggregating by account type.

#### `GetMonthlyStats`

Retrieves monthly income and expense statistics for the last N months.

```go
func (r *Repository) GetMonthlyStats(months int) ([]model.MonthlyStat, error)
```

**Parameters:**

- `months`: Number of months to look back.

**Returns:** A slice of `MonthlyStat` structs, one per month.

**Example:**

```go
stats, err := repo.GetMonthlyStats(6) // Last 6 months
if err != nil {
    log.Fatal(err)
}
for _, stat := range stats {
    fmt.Printf("%s: Income $%.2f, Expense $%.2f\n",
        stat.Month, stat.Income, stat.Expense)
}
```

### Budgets

#### `CreateBudget`

Creates a new budget for a category.

```go
func (r *Repository) CreateBudget(b *model.Budget) error
```

**Parameters:**

- `b`: Pointer to a `Budget` struct with `CategoryID`, `Amount` (in cents), and `Period` set.

**Returns:** Error if the insert fails. On success, `b.ID` is populated.

#### `GetBudgetsWithProgress`

Retrieves budgets for a specific month/year with calculated spent amounts and progress indicators.

```go
func (r *Repository) GetBudgetsWithProgress(month int, year int) ([]model.BudgetProgress, error)
```

**Parameters:**

- `month`: Month number (1-12).
- `year`: Year (e.g., 2024).

**Returns:** A slice of `BudgetProgress` structs containing:

- `CategoryName`: Name of the category
- `Budgeted`: Budgeted amount for the period
- `Spent`: Actual amount spent
- `Remaining`: Remaining budget
- `Percent`: Percentage of budget used (0.0 to 1.0+)

**Example:**

```go
progress, err := repo.GetBudgetsWithProgress(1, 2024) // January 2024
if err != nil {
    log.Fatal(err)
}
for _, p := range progress {
    fmt.Printf("%s: $%.2f / $%.2f (%.1f%%)\n",
        p.CategoryName, p.Spent, p.Budgeted, p.Percent*100)
}
```

### Recurring Transactions

#### `DetectRecurringPatterns`

Analyzes transaction history to identify potential recurring subscriptions or bills.

```go
func (r *Repository) DetectRecurringPatterns() ([]model.Subscription, error)
```

**Returns:** A slice of `Subscription` structs representing detected recurring patterns.

**Note:** This is a heuristic-based detection that looks for transactions with the same description occurring multiple times in the last 3 months.

### Anomaly Detection

#### `DetectAnomalies`

Identifies unusual transactions that may require attention.

```go
func (r *Repository) DetectAnomalies() ([]model.Anomaly, error)
```

**Returns:** A slice of `Anomaly` structs representing detected anomalies.

**Current Implementation:** Flags transactions over $200 in the last month as "Large Transaction" anomalies.

### Forecasting

#### `GetNetWorthProjection`

Projects future net worth based on historical savings patterns.

```go
func (r *Repository) GetNetWorthProjection(monthsAhead int) ([]ProjectionPoint, float64, error)

type ProjectionPoint struct {
    Month string
    Value float64
}
```

**Parameters:**

- `monthsAhead`: Number of months to project into the future.

**Returns:**

- A slice of projection points (one per month)
- Average monthly savings rate
- Error if calculation fails

**Note:** Projection is based on average monthly savings from the last 3 months.

### Rules & Enrichment

#### `CreateRule`

Creates a rule for automatically enriching transactions based on description patterns.

```go
func (r *Repository) CreateRule(rule *model.Rule) error
```

**Parameters:**

- `rule`: Pointer to a `Rule` struct with pattern matching and target enrichment fields.

**Example:**

```go
rule := &model.Rule{
    Pattern:        "STARBUCKS",
    TargetCategoryID: &coffeeCategoryID,
    TargetPayee:    "Starbucks",
    TargetNote:     "Coffee purchase",
}
err := repo.CreateRule(rule)
```

#### `EnrichTransaction`

Applies rules to enrich a transaction description with category, payee, and note.

```go
func (r *Repository) EnrichTransaction(description string) (string, *int64, string)
```

**Parameters:**

- `description`: The original transaction description.

**Returns:**

- Normalized payee name (or original if no match)
- Category ID (or nil if no match)
- Note (or empty string if no match)

### Data Export

#### `ExportDataToJSON`

Exports all accounts and categories to a JSON file for backup purposes.

```go
func (r *Repository) ExportDataToJSON(filepath string) error
```

**Parameters:**

- `filepath`: Path where the JSON file should be saved.

**Returns:** Error if the export fails.

**Note:** This method exports accounts and categories only. Transaction data is not included in the current implementation.
