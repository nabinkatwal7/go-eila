package ui

import (
	"strconv"
	"time"

	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/nabinkatwal7/go-eila/internal/model"
)

func (a *App) ShowAddTransactionModal() {
	// Inputs
	amountEntry := widget.NewEntry()
	amountEntry.PlaceHolder = "Amount"

	noteEntry := widget.NewEntry()
	noteEntry.PlaceHolder = "Note (e.g., Lunch)"

	typeSelect := widget.NewSelect([]string{"Expense", "Income", "Transfer"}, nil)
	typeSelect.Selected = "Expense"

	// Date (default today)
	// For simplicity using Entry for now, could use DatePicker later
	dateEntry := widget.NewEntry()
	dateEntry.SetText(time.Now().Format("2006-01-02"))

	// Accounts & Categories (Mocked for now, needs Repo fetch)
	// In real app, fetch these from a.Repo.GetAllAccounts()
	accountSelect := widget.NewSelect([]string{"Cash", "Bank"}, nil)
	accountSelect.Selected = "Cash"

	categorySelect := widget.NewSelect([]string{"Food", "Transport", "Salary"}, nil)
	categorySelect.Selected = "Food"

	// Form Items
	items := []*widget.FormItem{
		widget.NewFormItem("Type", typeSelect),
		widget.NewFormItem("Amount", amountEntry),
		widget.NewFormItem("Date", dateEntry),
		widget.NewFormItem("Account", accountSelect),
		widget.NewFormItem("Category", categorySelect),
		widget.NewFormItem("Note", noteEntry),
	}

	dialog.ShowForm("Add Transaction", "Save", "Cancel", items, func(confirm bool) {
		if confirm {
			price, _ := strconv.ParseFloat(amountEntry.Text, 64)
			date, _ := time.Parse("2006-01-02", dateEntry.Text)

			// Map string to proper IDs (Mock logic)
			// In real app utilize the selected account/category to find ID

			t := &model.Transaction{
				Amount:     price,
				Date:       date,
				Note:       noteEntry.Text,
				Type:       model.TransactionType(typeSelect.Selected),
				AccountID:  1, // Default/Mock ID
				CategoryID: nil, // Handle this
			}

			// Adjust CategoryID
			catID := int64(1) // Mock
			t.CategoryID = &catID

			if err := a.Repo.CreateTransaction(t); err != nil {
				dialog.ShowError(err, a.Window)
			} else {
				// Refresh Dashboard/Transactions if visible
				// Verify: How to trigger refresh?
				// For now simple weak refresh if we are on that page
				// ideally use detailed callback or event bus
				a.ContentContainer.Refresh()
			}
		}
	}, a.Window)
}
