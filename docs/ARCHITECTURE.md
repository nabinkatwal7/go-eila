# System Architecture

MyTrack follows a clean, modular architecture designed for maintainability and testability. The application is built using Go with a clear separation of concerns across three main layers.

## Architecture Overview

```
┌─────────────────────────────────────────┐
│         Presentation Layer              │
│         (internal/ui)                   │
│  - Fyne UI Components                   │
│  - User Input Validation                │
│  - View Logic                           │
└─────────────────┬───────────────────────┘
                  │
                  │ Repository Interface
                  │
┌─────────────────▼───────────────────────┐
│         Data Layer                      │
│         (internal/repository)           │
│  - Database Operations                  │
│  - Business Logic                       │
│  - Transaction Management               │
└─────────────────┬───────────────────────┘
                  │
                  │ SQL Queries
                  │
┌─────────────────▼───────────────────────┐
│         Domain Model                    │
│         (internal/model)                 │
│  - Pure Go Structs                      │
│  - Business Entities                    │
│  - Type Definitions                     │
└─────────────────────────────────────────┘
                  │
                  │
┌─────────────────▼───────────────────────┐
│         SQLite Database                 │
│         (mytrack.db)                    │
└─────────────────────────────────────────┘
```

## Layers

### 1. Presentation Layer (`internal/ui`)

**Purpose:** Handles all user interface interactions and rendering.

**Technology:** Built using the [Fyne](https://fyne.io) toolkit v2, which provides cross-platform native GUI capabilities.

**Responsibilities:**

- Rendering UI components (buttons, forms, tables, charts)
- User input validation before data submission
- Event handling (button clicks, form submissions)
- Navigation between views
- Display formatting (currency, dates, percentages)

**Key Components:**

- `app.go`: Main application structure and window setup
- `dashboard.go`: Main dashboard view with statistics
- `transactions.go`: Transaction listing and management
- `add_transaction.go`: Transaction creation forms
- `accounts.go`: Account management interface
- `budgets.go`: Budget viewing and configuration
- `command_palette.go`: Quick navigation feature (Ctrl+K)
- `validation.go`: Input validation utilities

**Example Structure:**

```go
type App struct {
    Repo             *repository.Repository
    ContentContainer *fyne.Container
    FyneApp          fyne.App
    Window           fyne.Window
}

func NewApp(fyneApp fyne.App, w fyne.Window, repo *repository.Repository) *App
```

**Communication:** The UI layer communicates exclusively with the Data Layer through the `Repository` interface. It never directly accesses the database.

### 2. Data Layer (`internal/repository`)

**Purpose:** Abstracts database operations and provides business logic.

**Technology:** Uses `database/sql` with the `modernc.org/sqlite` pure Go SQLite driver (no CGO dependencies).

**Responsibilities:**

- Database connection management
- CRUD operations for all entities
- Complex queries (joins, aggregations)
- Transaction atomicity (ensuring double-entry balance)
- Data validation and business rules
- Statistics and analytics calculations

**Key Components:**

- `db.go`: Database connection and schema initialization
- `repository.go`: Main repository with all CRUD methods
- `backup.go`: Data export functionality

**Repository Pattern:**

```go
type Repository struct {
    DB *DB
}

func NewRepository(db *DB) *Repository {
    return &Repository{DB: db}
}
```

**Key Features:**

- All database access is centralized in this layer
- Methods are strongly typed using domain models
- Complex operations (like `CreateTransaction`) are atomic
- Error handling is consistent and meaningful

### 3. Domain Model (`internal/model`)

**Purpose:** Defines the core business entities and their relationships.

**Technology:** Pure Go structs with no external dependencies.

**Key Entities:**

**Account:**

```go
type Account struct {
    ID       int64
    Name     string
    Type     AccountType  // Cash, Bank, Card, Investment, Liability, Income, Expense
    Currency string
}
```

**Transaction:**

```go
type Transaction struct {
    ID          int64
    Date        time.Time
    Description string
    Note        string
    Status      TransactionStatus  // Pending, Cleared, Reconciled
    Splits      []Split
}
```

**Split:**

```go
type Split struct {
    ID            int64
    TransactionID int64
    AccountID     int64
    CategoryID    *int64  // Optional
    Amount        int64   // Stored in cents (minor units)
    Currency      string
    ExchangeRate  float64
}
```

**Design Principles:**

- No dependencies on UI or database packages
- Types are self-documenting
- Constants define valid values (e.g., `AccountType`, `TransactionStatus`)

## Double-Entry Accounting System

The core feature of MyTrack is its double-entry accounting engine, which ensures mathematical accuracy and prevents data inconsistencies.

### Core Concepts

**Double-Entry Accounting:** Every financial transaction affects at least two accounts, and the total debits must equal the total credits.

**Splits:** A single atomic movement of money. Each split represents one side of a transaction.

**Transaction:** A grouping of 2 or more splits that together represent a complete financial event. The sum of all split amounts in a transaction must equal zero.

### Schema Design

**`accounts` table:**

```sql
CREATE TABLE accounts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    type TEXT NOT NULL,           -- Cash, Bank, Card, Investment, Liability, Income, Expense
    currency TEXT DEFAULT 'USD',
    is_closed BOOLEAN DEFAULT 0
);
```

**`transactions` table:**

```sql
CREATE TABLE transactions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date DATETIME NOT NULL,
    description TEXT NOT NULL,
    note TEXT,
    status TEXT DEFAULT 'Pending',  -- Pending, Cleared, Reconciled
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

**`splits` table:**

```sql
CREATE TABLE splits (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    transaction_id INTEGER NOT NULL,
    account_id INTEGER NOT NULL,
    category_id INTEGER,           -- Optional, only for Income/Expense accounts
    amount INTEGER NOT NULL,       -- Stored in cents (minor units)
    currency TEXT DEFAULT 'USD',
    exchange_rate REAL DEFAULT 1.0,
    FOREIGN KEY(transaction_id) REFERENCES transactions(id) ON DELETE CASCADE,
    FOREIGN KEY(account_id) REFERENCES accounts(id),
    FOREIGN KEY(category_id) REFERENCES categories(id)
);
```

**`categories` table:**

```sql
CREATE TABLE categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    icon TEXT,
    color TEXT,
    parent_id INTEGER,             -- For hierarchical categories
    FOREIGN KEY(parent_id) REFERENCES categories(id)
);
```

**`budgets` table:**

```sql
CREATE TABLE budgets (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    category_id INTEGER NOT NULL,
    amount INTEGER NOT NULL,       -- Stored in cents
    period TEXT DEFAULT 'Monthly',
    FOREIGN KEY(category_id) REFERENCES categories(id)
);
```

### Transaction Flows

#### Expense Flow

When a user spends money, the system creates a balanced transaction:

1. User spends $50.00 from Cash account on Food category
2. System creates Transaction with 2 Splits:
   - Split 1: Account=Cash, Amount=-5000 (credit, money leaving)
   - Split 2: Account=Expense, Category=Food, Amount=+5000 (debit, expense recorded)
3. Sum of splits: -5000 + 5000 = 0 ✓

**Code Example:**

```go
transaction := &model.Transaction{
    Date:        time.Now(),
    Description: "Grocery Store",
    Splits: []model.Split{
        {AccountID: cashAccountID, Amount: -5000},
        {AccountID: expenseAccountID, CategoryID: &foodCategoryID, Amount: 5000},
    },
}
err := repo.CreateTransaction(transaction)
```

#### Income Flow

When a user receives income:

1. User receives $2000.00 salary into Bank account
2. System creates Transaction with 2 Splits:
   - Split 1: Account=Bank, Amount=+200000 (debit, money entering)
   - Split 2: Account=Income, Category=Salary, Amount=-200000 (credit, income recorded)
3. Sum of splits: 200000 - 200000 = 0 ✓

#### Transfer Flow

When moving money between accounts:

1. User transfers $100.00 from Bank to Cash
2. System creates Transaction with 2 Splits:
   - Split 1: Account=Bank, Amount=-10000 (credit, money leaving)
   - Split 2: Account=Cash, Amount=+10000 (debit, money entering)
3. Sum of splits: -10000 + 10000 = 0 ✓

### Balance Calculation

Account balances are calculated dynamically by summing all splits:

```go
func (r *Repository) GetAccountBalance(accountID int64) (float64, error) {
    var balanceCents sql.NullInt64
    err := r.DB.QueryRow(
        "SELECT SUM(amount) FROM splits WHERE account_id = ?",
        accountID,
    ).Scan(&balanceCents)

    if !balanceCents.Valid {
        return 0, nil
    }
    return float64(balanceCents.Int64) / 100.0, nil
}
```

**Key Points:**

- Balances are never stored; they're always calculated from splits
- This ensures consistency and prevents balance drift
- Positive amounts increase asset balances, negative amounts decrease them
- For liabilities, the sign convention may be inverted depending on accounting perspective

## Data Flow Example

Here's a complete example of adding a transaction:

1. **User Action:** User clicks "Add Transaction" and fills out the form
2. **UI Layer:** `add_transaction.go` validates input and creates a `Transaction` struct
3. **Repository Call:** UI calls `repo.CreateTransaction(transaction)`
4. **Validation:** Repository checks that splits sum to zero
5. **Database Transaction:** Repository begins a SQL transaction
6. **Insert Header:** Insert into `transactions` table
7. **Insert Splits:** Insert all splits into `splits` table
8. **Commit:** If all succeeds, commit; otherwise rollback
9. **UI Update:** UI refreshes to show the new transaction

## Error Handling

Errors flow from the database layer up to the UI:

- **Database Errors:** Caught in repository, wrapped with context
- **Validation Errors:** Returned immediately (e.g., "transaction not balanced")
- **UI Errors:** Displayed to user via Fyne dialogs

## Testing Considerations

The architecture supports testing at multiple levels:

- **Unit Tests:** Test repository methods with an in-memory SQLite database
- **Integration Tests:** Test complete flows (UI → Repository → Database)
- **UI Tests:** Manual testing required (Fyne doesn't have automated UI testing yet)

## Future Improvements

- **Migration System:** Replace schema creation with proper migrations (e.g., golang-migrate)
- **Caching Layer:** Add caching for frequently accessed data (account balances, stats)
- **Background Jobs:** Move heavy operations (recurring detection, anomaly analysis) to background goroutines
- **API Layer:** Add REST API for potential mobile app integration
