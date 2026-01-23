# MyTrack API Documentation

This document outlines the core Go API provided by the `internal/repository` package for interacting with the MyTrack database.

## Repository

The `Repository` struct is the main entry point.

```go
type Repository struct {
    DB *sql.DB
}
```

## Methods

### Accounts

#### `GetAllAccounts`
Retrieves all accounts from the database.

```go
func (r *Repository) GetAllAccounts() ([]model.Account, error)
```

#### `CreateAccount`
Creates a new account. Populates the `ID` field of the passed struct on success.

```go
func (r *Repository) CreateAccount(account *model.Account) error
```

#### `GetAccountByName`
Checks if an account exists with the given name. Returns `nil` (and no error) if not found.

```go
func (r *Repository) GetAccountByName(name string) (*model.Account, error)
```

### categories

#### `GetAllCategories`
Retrieves all categories.

```go
func (r *Repository) GetAllCategories() ([]model.Category, error)
```

### Transactions

#### `CreateTransaction`
Creates a new balanced transaction with multiple splits.
**Invariant**: The sum of all splits must be 0.

```go
func (r *Repository) CreateTransaction(t *model.Transaction) error
```

### Budgets

#### `GetBudgetsWithProgress`
Retrieves budgets for a specific month/year with calculated spent amounts.

```go
func (r *Repository) GetBudgetsWithProgress(month int, year int) ([]model.BudgetProgress, error)
```

#### `CreateOrUpdateBudget`
Sets a budget limit for a category. Upserts if it already exists for that month.

```go
func (r *Repository) CreateOrUpdateBudget(budget *model.Budget) error
```

### System

#### `ExportDataToJSON`
Exports all Accounts and Categories to a JSON file at the specified path.

```go
func (r *Repository) ExportDataToJSON(filepath string) error
```
