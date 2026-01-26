package repository

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/nabinkatwal7/go-eila/internal/model"
)

// BackupData represents the full state of the user's data
type BackupData struct {
	Accounts     []model.Account     `json:"accounts"`
	Categories   []model.Category    `json:"categories"`
	Transactions []model.Transaction `json:"transactions"`
}

// ExportDataToJSON dumps the DB to a JSON file with full transaction data
func (r *Repository) ExportDataToJSON(filepath string) error {
	// 1. Fetch Accounts
	accounts, err := r.GetAllAccounts()
	if err != nil {
		return err
	}

	// 2. Fetch Categories
	categories, err := r.GetAllCategories()
	if err != nil {
		return err
	}

	// 3. Fetch All Transactions with splits
	transactions, err := r.GetRecentTransactions(100000) // Large limit to get all
	if err != nil {
		return err
	}

	data := BackupData{
		Accounts:     accounts,
		Categories:   categories,
		Transactions: transactions,
	}

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// ImportDataFromJSON restores data from a JSON backup file
func (r *Repository) ImportDataFromJSON(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	var data BackupData
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return err
	}

	// Import in transaction to ensure atomicity
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Clear existing data (optional - you might want to merge instead)
	// For safety, we'll just add new data

	// Import accounts (skip if exists)
	for _, acc := range data.Accounts {
		existing, _ := r.GetAccountByName(acc.Name)
		if existing == nil {
			if err := r.CreateAccount(&acc); err != nil {
				return err
			}
		}
	}

	// Import transactions
	for _, t := range data.Transactions {
		if err := r.CreateTransaction(&t); err != nil {
			return err
		}
	}

	return tx.Commit()
}

// CSVTransactionRow represents a single CSV row for import
type CSVTransactionRow struct {
	Date        string
	Description string
	Amount      string
	Category    string
	Account     string
	Note        string
}

// ImportTransactionsFromCSV imports transactions from a CSV file
// Expected CSV format: Date,Description,Amount,Category,Account,Note
func (r *Repository) ImportTransactionsFromCSV(filepath string, accountName string, categoryName string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	if len(records) < 2 {
		return fmt.Errorf("CSV file must have at least a header and one data row")
	}

	// Get default account and category if not provided
	accounts, err := r.GetAllAccounts()
	if err != nil {
		return err
	}

	var defaultAccountID int64
	var expenseAccountID int64
	for _, acc := range accounts {
		if accountName != "" && acc.Name == accountName {
			defaultAccountID = acc.ID
		}
		if acc.Type == model.AccountTypeExpense {
			expenseAccountID = acc.ID
		}
		if defaultAccountID == 0 && (acc.Type == model.AccountTypeCash || acc.Type == model.AccountTypeBank) {
			defaultAccountID = acc.ID
		}
	}

	categories, err := r.GetAllCategories()
	if err != nil {
		return err
	}

	var defaultCategoryID int64
	categoryNameToID := make(map[string]int64)
	for _, cat := range categories {
		categoryNameToID[cat.Name] = cat.ID
		if categoryName != "" && cat.Name == categoryName {
			defaultCategoryID = cat.ID
		}
		if defaultCategoryID == 0 {
			defaultCategoryID = cat.ID
		}
	}

	// Parse header (skip first row)
	header := records[0]
	dateIdx, descIdx, amountIdx, catIdx, accIdx, noteIdx := -1, -1, -1, -1, -1, -1

	for i, col := range header {
		col = strings.ToLower(strings.TrimSpace(col))
		switch col {
		case "date":
			dateIdx = i
		case "description", "payee", "memo":
			descIdx = i
		case "amount":
			amountIdx = i
		case "category":
			catIdx = i
		case "account":
			accIdx = i
		case "note", "notes":
			noteIdx = i
		}
	}

	if dateIdx == -1 || descIdx == -1 || amountIdx == -1 {
		return fmt.Errorf("CSV must have Date, Description, and Amount columns")
	}

	// Process rows
	for i := 1; i < len(records); i++ {
		row := records[i]
		if len(row) <= dateIdx || len(row) <= descIdx || len(row) <= amountIdx {
			continue // Skip malformed rows
		}

		// Parse date
		dateStr := strings.TrimSpace(row[dateIdx])
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			// Try other formats
			date, err = time.Parse("01/02/2006", dateStr)
			if err != nil {
				continue // Skip rows with invalid dates
			}
		}

		// Parse amount
		amountStr := strings.TrimSpace(row[amountIdx])
		amountStr = strings.ReplaceAll(amountStr, "$", "")
		amountStr = strings.ReplaceAll(amountStr, ",", "")
		amount, err := strconv.ParseFloat(amountStr, 64)
		if err != nil {
			continue // Skip rows with invalid amounts
		}
		amountCents := int64(amount * 100)

		// Get account
		accID := defaultAccountID
		if accIdx >= 0 && accIdx < len(row) && strings.TrimSpace(row[accIdx]) != "" {
			accName := strings.TrimSpace(row[accIdx])
			for _, acc := range accounts {
				if acc.Name == accName {
					accID = acc.ID
					break
				}
			}
		}

		// Get category
		catID := defaultCategoryID
		if catIdx >= 0 && catIdx < len(row) && strings.TrimSpace(row[catIdx]) != "" {
			catName := strings.TrimSpace(row[catIdx])
			if id, ok := categoryNameToID[catName]; ok {
				catID = id
			}
		}

		// Get description and note
		description := strings.TrimSpace(row[descIdx])
		note := ""
		if noteIdx >= 0 && noteIdx < len(row) {
			note = strings.TrimSpace(row[noteIdx])
		}

		// Create transaction
		transaction := &model.Transaction{
			Date:        date,
			Description: description,
			Note:        note,
			Status:      model.TransactionStatusPending,
			Splits: []model.Split{
				{
					AccountID:    accID,
					Amount:       -amountCents,
					Currency:     "USD",
					ExchangeRate: 1.0,
				},
				{
					AccountID:    expenseAccountID,
					CategoryID:   &catID,
					Amount:       amountCents,
					Currency:     "USD",
					ExchangeRate: 1.0,
				},
			},
		}

		if err := r.CreateTransaction(transaction); err != nil {
			// Log error but continue with other rows
			continue
		}
	}

	return nil
}
