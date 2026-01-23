# System Architecture

MyTrack follows a clean, modular architecture designed for maintainability and testability.

## Layers

### 1. Presentation Layer (`internal/ui`)
*   Built using the **Fyne** toolkit.
*   Handles all user interactions, rendering, and input validation.
*   Contains sub-packages/files for each major view (Dashboard, Transactions, Accounts, etc.).
*   Communicates with the Data Layer via the Repository interface.

### 2. Data Layer (`internal/repository`)
*    abstracts the underlying SQLite database.
*   Provides strongly-typed methods for CRUD operations (e.g., `CreateTransaction`, `GetAllAccounts`).
*   Handles complex logic like multi-table joins for reports.
*   Manages database transactions for atomic operations.

### 3. Domain Model (`internal/model`)
*   Contains pure Go structs representing the business entities.
*   Examples: `Transaction`, `Account`, `Split`, `Budget`.
*   No dependencies on UI or heavy database logic.

## Double-Entry Accounting System

The core feature of MyTrack is its double-entry engine.

### Concepts
*   **Splits**: A single atomic movement of money.
*   **Transaction**: A grouping of 2 or more splits that sum to zero.

### Schema Design

`transactions` table:
*   Header information (Date, Payee, Note)

`splits` table:
*   `transaction_id`: FK to transactions
*   `account_id`: FK to accounts
*   `amount`: Integer (cents)
*   `category_id`: Optional FK to categories (only for Income/Expense legs)

### Flows

**Expense Flow**:
1.  User spends money from an Asset account (e.g., Cash).
2.  System creates a Split: Credit Asset (-Amount).
3.  System creates a Split: Debit Expense (+Amount).
4.  `Sum(Splits) == 0`.

**Income Flow**:
1.  User receives money into an Asset account.
2.  System creates a Split: Debit Asset (+Amount).
3.  System creates a Split: Credit Income (-Amount).
