# MyTrack

A premium, EILA-style money tracking application built with Go and Fyne.

## Features (Planned)
- 4 Pillars: Expenses, Income, Assets, Liabilities
- Fast Add Flow
- Dashboard with Insights
- Budgets & recurring transactions
- Privacy-focused (Local SQLite)

## Prerequisites
- **Go** (1.20+)
- **C Compiler (GCC)**: Required for Fyne.
    - **Windows**: Install [TDM-GCC](https://jmeubank.github.io/tdm-gcc/) or [MinGW-w64](https://www.mingw-w64.org/).
    - **Linux**: `sudo apt install gcc libgl1-mesa-dev xorg-dev`
    - **macOS**: Install Xcode Command Line Tools.

## Running
```bash
# Ensure CGO is enabled
$env:CGO_ENABLED="1"
go run cmd/mytrack/main.go
```
