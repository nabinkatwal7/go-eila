# MyTrack Quick Reference

A quick reference guide for common tasks and features in MyTrack.

## Navigation

| Action | Shortcut/Location |
|--------|------------------|
| Open Command Palette | `Ctrl+K` (Windows/Linux) or `Cmd+K` (Mac) |
| Add Transaction | Click **+ Add New** in sidebar |
| View Dashboard | Click **Dashboard** in sidebar |
| View Transactions | Click **Transactions** in sidebar |
| View Accounts | Click **Accounts** in sidebar |
| View Budgets | Click **Budgets** in sidebar |
| Settings | Click **Settings** in sidebar |

## Account Types

| Type | Description | Use Case |
|------|-------------|----------|
| **Cash** | Physical cash | Wallet, petty cash |
| **Bank** | Bank accounts | Checking, savings accounts |
| **Card** | Credit/debit cards | Credit cards, prepaid cards |
| **Investment** | Investment accounts | Stocks, bonds, retirement accounts |
| **Liability** | Debts | Loans, mortgages, credit card debt |
| **Income** | Income accounts | Salary, freelance income (system) |
| **Expense** | Expense accounts | Spending categories (system) |

## Transaction Types

### Simple Transaction
- **Use for:** Single-category purchases
- **Example:** $50 grocery purchase from Cash account
- **Steps:** Amount → Date → Account → Category → Save

### Split Transaction
- **Use for:** Multi-category purchases
- **Example:** $100 Target purchase: $60 groceries, $40 clothes
- **Steps:** Total amount → Add splits → Ensure balance = 0 → Save

## Common Workflows

### Recording Daily Expenses
1. Click **+ Add New**
2. Select **Simple** tab
3. Enter amount (e.g., "25.50")
4. Select date (defaults to today)
5. Choose account (e.g., "Cash" or "Chase Checking")
6. Choose category (e.g., "Food")
7. Click **Save**

### Setting Up a Budget
1. Navigate to **Budgets**
2. Click **Set/Update Budget**
3. Select category (e.g., "Food")
4. Enter monthly limit (e.g., "500.00")
5. Click **Save**

### Creating a New Account
1. Navigate to **Accounts**
2. Click **New Account**
3. Enter name (e.g., "Chase Savings")
4. Select type (e.g., "Bank")
5. Select currency (default: "USD")
6. Click **Create**

### Viewing Financial Health
1. Navigate to **Dashboard**
2. Review:
   - Net Worth (top)
   - Income vs. Expense chart (6 months)
   - Account balances (list)
   - Budget progress (if budgets are set)

## Data Format

- **Amounts:** Enter as decimal (e.g., "25.50" for $25.50)
- **Dates:** Format YYYY-MM-DD (e.g., "2024-01-15")
- **Currency:** Currently supports USD (multi-currency planned)

## Tips & Tricks

1. **Quick Entry:** Use the command palette (`Ctrl+K`) to quickly navigate to any feature
2. **Split Transactions:** When shopping at stores with multiple categories, use split transactions for better tracking
3. **Regular Reviews:** Check the dashboard weekly to monitor your financial health
4. **Budget Alerts:** Red progress bars indicate you've exceeded your budget
5. **Anomaly Detection:** Large transactions (>$200) are automatically flagged for review

## Keyboard Shortcuts

| Shortcut | Action |
|----------|--------|
| `Ctrl+K` / `Cmd+K` | Open command palette |
| `Esc` | Close dialog/modal |
| `Enter` | Submit form (when focused) |

## Error Messages

| Error | Cause | Solution |
|-------|-------|----------|
| "Transaction is not balanced" | Split amounts don't sum to zero | Adjust split amounts until total is 0.00 |
| "Account not found" | Invalid account selected | Select a valid account from the dropdown |
| "Category not found" | Invalid category selected | Select a valid category from the dropdown |
| "Invalid date format" | Date not in YYYY-MM-DD format | Enter date as YYYY-MM-DD (e.g., 2024-01-15) |

## Account Balance Calculation

Account balances are calculated automatically from all transactions:
- **Assets:** Positive amounts increase balance, negative decrease
- **Liabilities:** Negative amounts represent debt (balance increases with more debt)
- **Income/Expense:** Used for tracking, not traditional balances

## Budget Progress

Budget progress is calculated monthly:
- **Green bar:** Under budget (< 100%)
- **Yellow bar:** Approaching limit (80-100%)
- **Red bar:** Over budget (> 100%)

## Data Storage

- **Location:** `mytrack.db` in the application directory
- **Format:** SQLite database
- **Backup:** Export to JSON via Settings → Export Data to JSON
- **Privacy:** All data stored locally, never sent to cloud

## Getting Help

1. Check the [User Guide](USER_GUIDE.md) for detailed instructions
2. Review the [Architecture Documentation](ARCHITECTURE.md) for technical details
3. See the [Roadmap](ROADMAP.md) for planned features
4. Check [Contributing Guidelines](../CONTRIBUTING.md) if you want to contribute
