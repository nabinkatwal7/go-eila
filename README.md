# MyTrack

**A premium, professional-grade personal finance application built with Go and Fyne.**

MyTrack is a privacy-focused, double-entry accounting money tracker that helps you manage your finances with professional-grade features. All data is stored locally in SQLite - no cloud, no tracking, complete privacy.

![MyTrack Dashboard](assets/screenshot.png)

## Overview

MyTrack is designed for users who want granular control over their finances. It combines the rigor of double-entry accounting with a modern, fast, and intuitive user interface.

## Highlights

*   **Privacy First**: Your financial data never leaves your device.
*   **Double-Entry Accuracy**: Every transaction is balanced, ensuring zero discrepancies.
*   **Performance**: Native Go application with minimal resource footprint.
*   **Cross-Platform**: Runs on Windows, macOS, and Linux.

## Features

### Core Money Tracking
*   **Expenses**: Track daily spending with detailed categorization.
*   **Income**: Record various income sources.
*   **Assets**: Monitor cash, bank accounts, and investment portfolios.
*   **Liabilities**: Track debts, credit cards, and loans.

### Dashboard & Insights
*   **Real-time Net Worth**: Instant calculation of your financial health.
*   **Trend Analysis**: 6-month income vs. expense charts.
*   **Account Balances**: Live updates of all asset and liability accounts.

### Budget Management
*   **Monthly Budgets**: Set spending limits per category.
*   **Visual Tracking**: Progress bars indicating budget utilization.
*   **Alerts**: Visual indicators when budgets are exceeded.

### Smart Operations
*   **Split Transactions**: Categorize a single transaction across multiple categories.
*   **Recurring Detection**: Automatically identify subscription patterns.
*   **Anomaly Detection**: Flag unusual large transactions.
*   **Data Export**: Backup your data to JSON for portability.

### Power User Tools
*   **Command Palette (Ctrl+K)**: Rapid navigation to any feature.
*   **Keyboard Shortcuts**: Fast data entry.

## Architecture

MyTrack relies on a robust stack:

*   **Language**: Go (Golang) 1.20+
*   **GUI Framework**: Fyne v2
*   **Database**: SQLite (via `modernc.org/sqlite` pure Go driver)
*   **Architecture Pattern**: Repository Pattern with MVC-like separation.

### Data Model

The application uses a strict double-entry system. A `Transaction` consists of multiple `Splits`. The sum of all splits in a transaction must always equal zero.

```
Transaction
├── ID
├── Date
├── Description
└── Splits []
    ├── AccountID (Source/Destination)
    ├── CategoryID (Optional)
    └── Amount (Positive for Debit/Expense, Negative for Credit/Income)
```

## Installation

### Prerequisites

1.  **Go**: Version 1.20 or later.
2.  **C Compiler**: Required for Fyne's graphical drivers (OpenGL).
    *   **Windows**: TDM-GCC or MinGW-w64.
    *   **Linux**: GCC (`sudo apt install gcc libgl1-mesa-dev xorg-dev`).
    *   **macOS**: Xcode Command Line Tools.

### Steps

1.  **Clone the Repository**
    ```bash
    git clone https://github.com/nabinkatwal7/go-eila.git
    cd go-eila
    ```

2.  **Install Dependencies**
    ```bash
    go mod download
    ```

3.  **Run Application**
    ```bash
    go run cmd/mytrack/main.go
    ```

## Usage Guide

### First Run
On the first launch, MyTrack will create a local `mytrack.db` file and seed it with default accounts (Cash, General Expenses) and categories (Food, Transport).

### Adding Transactions
1.  Click **+ Add New** in the sidebar.
2.  Select **Simple** for one-off expenses or **Split** for complex receipts.
3.  Enter amount, date, and select Account/Category.
4.  Click Save.

### Managing Budgets
1.  Navigate to the **Budgets** view.
2.  Click **Set/Update Budget**.
3.  Select a category and define the monthly limit.

## Development

### Project Structure
*   `cmd/mytrack/`: Entry point.
*   `internal/model/`: Domain models (structs).
*   `internal/repository/`: Database logic and queries.
*   `internal/ui/`: Fyne UI components and view logic.

### Building for Production
To create a standalone binary:
```bash
go build -ldflags="-s -w" -o mytrack.exe cmd/mytrack/main.go
```

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

This project is open source and available under the [MIT License](LICENSE).
