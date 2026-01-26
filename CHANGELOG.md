# Changelog

All notable changes to MyTrack will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Planned
- Transaction editing and deletion
- CSV import for bank statements
- Enhanced data visualization with real charts
- Multi-currency support with live exchange rates
- Mobile companion app

## [1.0.0] - 2024

### Added
- Initial release of MyTrack
- Double-entry accounting system with transaction validation
- Account management (Cash, Bank, Card, Investment, Liability)
- Category-based expense tracking
- Transaction creation (Simple and Split transactions)
- Budget management with monthly limits and progress tracking
- Dashboard with real-time statistics:
  - Net worth calculation
  - Income vs. expense trends
  - Account balances
- Recurring transaction detection
- Anomaly detection for large transactions
- Net worth forecasting based on historical savings
- Data export to JSON for backups
- Command palette (Ctrl+K) for quick navigation
- Cross-platform support (Windows, macOS, Linux)

### Technical
- Go 1.25+ with Fyne v2 GUI framework
- SQLite database with pure Go driver (no CGO)
- Repository pattern for clean architecture
- Local data storage (privacy-first design)

## [0.1.0] - Development Phase

### Added
- Basic transaction recording
- Account and category management
- SQLite database schema
- Fyne UI framework integration
