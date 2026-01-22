package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/nabinkatwal7/go-eila/internal/repository"
)

func NewTransactionsView(repo *repository.Repository) fyne.CanvasObject {
	txs, err := repo.GetRecentTransactions(50)
	if err != nil {
		return widget.NewLabel("Error loading transactions: " + err.Error())
	}

	// Table Data
	// [row][col]
	var data [][]string
	tableHeader := []string{"Date", "Payee", "Amount"} // Simple columns
	// Maybe add Category but getting category name requires join or fetch.
	// For now, let's stick to simple props available in Transaction struct + logic.

	for _, t := range txs {
		// Calculate amount (sum of positive splits? or negative if expense?)
		// Logic from CreateTransaction: Expenses are debits (positive amount).
		// Let's sum splits.
		var amt int64
		for _, s := range t.Splits {
			if s.Amount > 0 {
				amt += s.Amount
			}
		}

		amtStr := fmt.Sprintf("$%.2f", float64(amt)/100.0)
		dateStr := t.Date.Format("2006-01-02")

		data = append(data, []string{dateStr, t.Description, amtStr})
	}

	table := widget.NewTable(
		func() (int, int) {
			return len(data), len(tableHeader)
		},
		func() fyne.CanvasObject {
			return widget.NewLabelWithStyle("Cell", fyne.TextAlignLeading, fyne.TextStyle{})
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			label := o.(*widget.Label)
			if i.Row < len(data) && i.Col < len(data[i.Row]) {
				label.SetText(data[i.Row][i.Col])
				// Styling
				if i.Col == 2 { // Amount
					label.TextStyle = fyne.TextStyle{Bold: true}
					label.Alignment = fyne.TextAlignTrailing
				} else {
					label.Alignment = fyne.TextAlignLeading
				}
			}
		},
	)

	// Column widths
	table.SetColumnWidth(0, 100) // Date
	table.SetColumnWidth(1, 400) // Payee
	table.SetColumnWidth(2, 100) // Amount

	return container.NewBorder(
		widget.NewLabelWithStyle("Transactions (Dense View)", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		nil, nil, nil,
		table,
	)
}
