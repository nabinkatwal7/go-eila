# MyTrack

**A premium, EILA-style personal finance application built with Go and Fyne**

MyTrack is a privacy-focused, double-entry accounting money tracker that helps you manage your finances with professional-grade features. All data is stored locally in SQLite—no cloud, no tracking, complete privacy.

![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Go Version](https://img.shields.io/badge/go-1.20+-00ADD8.svg)

## ✨ Features

### 📊 Core Money Tracking (The 4 Pillars)
- **Expenses**: Track daily spending with categories
- **Income**: Record all income sources
- **Assets**: Monitor cash, bank accounts, investments
- **Liabilities**: Track debts and credit cards

### 💰 Double-Entry Accounting
- Professional accounting system with balanced transactions
- Every transaction has debits and credits
- Accurate balance calculations
- Support for split transactions (multiple categories in one transaction)

### 📈 Dashboard & Insights
- Real-time net worth calculation
- Income vs Expense trends (6-month chart)
- Account balances
- Category-wise spending breakdown

### 🎯 Budgets
- Set monthly budgets per category
- Real-time progress tracking
- Visual indicators for budget health
- Overspending alerts

### 🔄 Smart Features
- **Recurring Transaction Detection**: Automatically identifies subscriptions
- **Spending Anomaly Detection**: Flags unusual large transactions
- **Rule Engine**: Auto-categorize transactions based on patterns
- **Net Worth Forecasting**: Project future wealth based on trends

### 🛠️ Tools
- Debt payoff calculator
- Tax calculator
- Invoice generation
- Multi-currency support

### ⌨️ Premium UX
- **Command Palette** (Ctrl+K): Quick access to all features
- Fast transaction entry
- Keyboard shortcuts
- Clean, modern UI

## 🏗️ Architecture

### Database Schema (SQLite)

```
accounts
├── id (PK)
├── name
├── type (Cash, Bank, Card, Invest, Expense, Income, Liability, Equity)
└── currency

categories
├── id (PK)
├── name
├── icon
├── color
└── parent_id (FK)

transactions
├── id (PK)
├── date
├── description (payee)
├── note
└── status

splits (double-entry legs)
├── id (PK)
├── transaction_id (FK)
├── account_id (FK)
├── category_id (FK)
├── amount (in cents)
├── currency
└── exchange_rate

budgets
├── id (PK)
├── category_id (FK)
├── amount
└── period

rules (auto-categorization)
├── id (PK)
├── pattern
├── target_category_id (FK)
├── target_payee
└── target_note
```

### Double-Entry Accounting Example

**Expense Transaction**: Spent $50 on groceries from Cash account
```
Splits:
1. Cash (Asset)      -$50  (Credit - decrease)
2. Expenses (Expense) +$50  (Debit - increase)
Category: Food
```

**Income Transaction**: Received $1000 salary to Bank account
```
Splits:
1. Bank (Asset)      +$1000 (Debit - increase)
2. Income (Income)   -$1000 (Credit - decrease)
Category: Salary
```

The sum of all splits in a transaction must equal zero (balanced).

## 🚀 Getting Started

### Prerequisites

1. **Go** (1.20 or later)
   ```bash
   go version
   ```

2. **C Compiler** (Required for Fyne GUI framework)
   - **Windows**: Install [TDM-GCC](https://jmeubank.github.io/tdm-gcc/) or [MinGW-w64](https://www.mingw-w64.org/)
   - **Linux**:
     ```bash
     sudo apt install gcc libgl1-mesa-dev xorg-dev
     ```
   - **macOS**: Install Xcode Command Line Tools
     ```bash
     xcode-select --install
     ```

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/nabinkatwal7/go-eila.git
   cd go-eila
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Run the application**
   ```bash
   go run cmd/mytrack/main.go
   ```

   The application will:
   - Create `mytrack.db` in the current directory
   - Seed default accounts and categories
   - Open the GUI window

### Building

To create a standalone executable:

```bash
go build -o mytrack.exe cmd/mytrack/main.go
```

## 📖 Usage Guide

### Adding a Transaction

1. **Quick Add**: Click the "+ Add New" button in the sidebar
2. **Choose Mode**:
   - **Simple**: Single category transaction
   - **Split**: Multiple categories in one transaction
3. **Fill Details**:
   - Type (Expense/Income)
   - Amount
   - Date
   - Account (where money comes from/goes to)
   - Category
   - Note (optional)
4. **Save**: Transaction is saved and modal closes

### Managing Accounts

1. Navigate to **Accounts** view
2. Click **"New Account"**
3. Enter:
   - Name (e.g., "Savings Account")
   - Type (Cash, Bank, Card, Invest)
   - Currency (USD, EUR, GBP, NPR, JPY)
4. View all accounts with real-time balances

### Setting Budgets

1. Navigate to **Budgets** view
2. Click **"New Budget"**
3. Select category and set monthly limit
4. Track progress in real-time

### Using Command Palette

Press **Ctrl+K** to open the command palette for quick navigation:
- Add Transaction
- Add Account
- Go to Dashboard
- Go to Transactions
- Go to Budgets
- And more...

### Viewing Insights

- **Dashboard**: Overview of financial health
- **Transactions**: Complete transaction history
- **Recurring**: Detected subscriptions
- **Alerts**: Spending anomalies
- **Forecast**: Net worth projection

## 🗂️ Project Structure

```
mytrack/
├── cmd/
│   └── mytrack/
│       └── main.go           # Application entry point
├── internal/
│   ├── model/                # Data models
│   │   ├── models.go         # Account, Transaction, Split, Category
│   │   ├── anomaly.go        # Anomaly detection models
│   │   ├── recurring.go      # Subscription models
│   │   └── stats.go          # Statistics models
│   ├── repository/           # Database layer
│   │   ├── db.go             # Database setup & schema
│   │   └── repository.go     # CRUD operations
│   └── ui/                   # Fyne UI components
│       ├── app.go            # Main app structure
│       ├── dashboard.go      # Dashboard view
│       ├── transactions.go   # Transaction list
│       ├── add_transaction.go # Transaction modal
│       ├── accounts.go       # Accounts view
│       ├── budgets.go        # Budget management
│       ├── command_palette.go # Quick actions
│       └── ...               # Other views
├── go.mod
├── go.sum
└── README.md
```

## 🔧 Development

### Database Migrations

The schema is automatically created on first run. For a fresh start:

```bash
rm mytrack.db
go run cmd/mytrack/main.go
```

### Adding New Features

1. **Model**: Define data structure in `internal/model/`
2. **Repository**: Add CRUD methods in `internal/repository/repository.go`
3. **UI**: Create view in `internal/ui/`
4. **Wire Up**: Add navigation in `app.go` and `command_palette.go`

### Testing

```bash
go test ./...
```

## 🐛 Troubleshooting

### "C compiler 'gcc' not found"

**Solution**: Install a C compiler (see Prerequisites above)

### Application hangs on startup

**Solution**: This was fixed in recent updates. Ensure you're using the latest version.

### Transactions not saving

**Solution**: Check that:
1. You have selected an account and category
2. Amount is a valid number
3. Database file `mytrack.db` has write permissions

### Window doesn't open

**Solution**:
- On Linux, ensure X11 or Wayland is running
- On Windows, check Windows Defender isn't blocking the app
- Try running with `CGO_ENABLED=1 go run cmd/mytrack/main.go`

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🙏 Acknowledgments

- [Fyne](https://fyne.io/) - Cross-platform GUI toolkit
- [modernc.org/sqlite](https://gitlab.com/cznic/sqlite) - Pure Go SQLite driver
- Inspired by EILA and other modern finance apps

## 📧 Contact

Nabin Katwal - [@nabinkatwal7](https://github.com/nabinkatwal7)

Project Link: [https://github.com/nabinkatwal7/go-eila](https://github.com/nabinkatwal7/go-eila)

---

**Built with ❤️ using Go and Fyne**
