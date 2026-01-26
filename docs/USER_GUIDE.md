# MyTrack User Guide

Welcome to the MyTrack User Guide. This document provides detailed instructions on how to use the application effectively.

## 1. Getting Started

### Installation

Please refer to the [README](../README.md#installation) for installation instructions.

### Initial Setup

When you first launch MyTrack, it automatically creates a database file `mytrack.db` in the application directory. It also seeds some default accounts and categories to help you get started immediately.

## 2. Managing Data

### Accounts

Accounts represent where your money is (Assets) or what you owe (Liabilities).

- **Assets**: Cash, Bank Accounts, Savings, Investments.
- **Liabilities**: Credit Cards, Loans, Mortgages.
- **Income/Expense**: These are system accounts used for double-entry calculation; you typically don't manage these directly.

**To Create an Account:**

1. Navigate to the **Accounts** view.
2. Click **New Account**.
3. Enter the **Name** (e.g., "Chase Checking").
4. Select the **Type** (e.g., Bank).
5. Choose the **Currency**.
6. Click **Create**.

### Categories

Categories help you organize your spending.
_Currently, categories are managed via the database seeding. Future versions will allow UI management._

## 3. Transactions

Transactions are the core of MyTrack. Every movement of money is a transaction.

### Simple Transaction

Use this for everyday purchases.

1. Click **+ Add New** in the sidebar.
2. Select the **Simple** tab.
3. Enter the **Amount**.
4. Enter the **Date** (YYYY-MM-DD).
5. Select the **Account** (where money comes from/goes to).
6. Select the **Category** (what it was for).
7. Click **Save**.

### Split Transaction

Use this when a single payment covers multiple categories (e.g., a "Target" receipt with groceries and clothes).

1. Click **+ Add New**.
2. Select the **Split** tab.
3. Enter the total amount and top-level details.
4. Add splits for each category/account.
5. Ensure the **Unassigned** amount is 0.00.
6. Click **Save**.

## 4. Budgets

Budgets help you control spending.

1. Go to **Budgets**.
2. Click **Set/Update Budget**.
3. Select a **Category** (e.g., "Food").
4. Enter the monthly limit (e.g., "500.00").
5. Click **Save**.

The bar will turn red if you exceed this amount in the current month.

## 5. Tools & Settings

### Command Palette

Press `Ctrl+K` (or `Cmd+K` on Mac) to open the Command Palette for quick navigation.

### Data Backup

1. Go to **Settings**.
2. Click **Export Data to JSON**.
3. Select a location to save your backup file.

**Note:** The current export includes accounts and categories. Transaction data export is planned for a future release.

## 6. Understanding Double-Entry Accounting

MyTrack uses double-entry accounting, which means every transaction affects at least two accounts and the total debits equal the total credits.

### How It Works

- When you spend money, it decreases an asset (like Cash) and increases an expense
- When you earn money, it increases an asset (like Bank) and increases income
- When you transfer money, it decreases one account and increases another

### Why This Matters

- **Accuracy:** The system prevents unbalanced transactions
- **Completeness:** Every dollar is accounted for
- **Reliability:** Your account balances are always mathematically correct

## 7. Dashboard Overview

The dashboard provides a comprehensive view of your financial health:

- **Net Worth:** Total assets minus total liabilities
- **Income vs. Expense:** 6-month trend showing your spending patterns
- **Account Balances:** Real-time balances for all your accounts
- **Budget Progress:** Visual indicators showing how much of your budgets you've used

## 8. Keyboard Shortcuts

- `Ctrl+K` (or `Cmd+K` on Mac): Open command palette
- `Esc`: Close dialogs/modals
- `Enter`: Submit forms (when focused)

## 9. Troubleshooting

### Transaction Won't Save

- **Check:** Ensure all required fields are filled
- **For Split Transactions:** Make sure the "Unassigned" amount is exactly 0.00
- **Error Message:** Read the error dialog for specific guidance

### Account Balance Looks Wrong

- **Recalculate:** Account balances are calculated from all transactions
- **Check Transactions:** Review recent transactions for that account
- **Verify:** Ensure transactions are properly categorized

### Budget Not Showing

- **Check Date:** Budgets are shown for the current month by default
- **Verify Category:** Ensure you've set a budget for that category
- **Check Transactions:** Make sure transactions are categorized correctly

### Application Won't Start

- **Check Database:** Ensure `mytrack.db` is not corrupted (try backing up and deleting it to start fresh)
- **Check Permissions:** Ensure the application has write permissions in its directory
- **Check Dependencies:** Verify all Go dependencies are installed (`go mod download`)

## 10. Best Practices

1. **Regular Entry:** Enter transactions daily or weekly to maintain accurate records
2. **Categorize Consistently:** Use the same categories for similar expenses
3. **Review Budgets Monthly:** Adjust budgets based on your actual spending patterns
4. **Backup Regularly:** Export your data periodically to prevent data loss
5. **Reconcile Accounts:** Periodically verify account balances match your bank statements
6. **Use Split Transactions:** For complex purchases (like grocery stores that sell multiple item types), use split transactions for better categorization
