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
*Currently, categories are managed via the database seeding. Future versions will allow UI management.*

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
