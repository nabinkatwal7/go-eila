package ui

import (
	"fmt"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/nabinkatwal7/go-eila/internal/model"
)

func (a *App) ShowAddTransactionModal() {
	// We will use Tabs for "Simple" vs "Split" modes

	simpleContent := a.createSimpleForm()
	splitContent := a.createSplitForm()

	tabs := container.NewAppTabs(
		container.NewTabItem("Simple", simpleContent),
		container.NewTabItem("Split", splitContent),
	)

	// Wrap in a custom dialog or just show a window?
	// Fyne standard dialogs take content.

	// Create a custom dialog window to handle the resizing better than standard dialog
	w := a.FyneApp.NewWindow("Add Transaction")
	w.Resize(fyne.NewSize(500, 600))
	w.SetContent(container.NewPadded(tabs))
	w.Show()
}

func (a *App) createSimpleForm() fyne.CanvasObject {
	// Inputs
	amountEntry := widget.NewEntry()
	amountEntry.PlaceHolder = "Amount"

	noteEntry := widget.NewEntry()
	noteEntry.PlaceHolder = "Note (e.g., Lunch)"

	typeSelect := widget.NewSelect([]string{"Expense", "Income"}, nil)
	typeSelect.Selected = "Expense"

	dateEntry := widget.NewEntry()
	dateEntry.SetText(time.Now().Format("2006-01-02"))

	// Mocks
	accountSelect := widget.NewSelect([]string{"Cash", "Bank"}, nil)
	accountSelect.Selected = "Cash"

	categorySelect := widget.NewSelect([]string{"Food", "Transport", "Salary"}, nil)
	categorySelect.Selected = "Food"

	form := widget.NewForm(
		widget.NewFormItem("Type", typeSelect),
		widget.NewFormItem("Amount", amountEntry),
		widget.NewFormItem("Date", dateEntry),
		widget.NewFormItem("Account", accountSelect),
		widget.NewFormItem("Category", categorySelect),
		widget.NewFormItem("Note", noteEntry),
	)

	saveBtn := widget.NewButtonWithIcon("Save Simple", theme.DocumentSaveIcon(), func() {
		priceFloat, _ := strconv.ParseFloat(amountEntry.Text, 64)
		amountCents := int64(priceFloat * 100)
		date, _ := time.Parse("2006-01-02", dateEntry.Text)

		t := &model.Transaction{
			Date:        date,
			Description: categorySelect.Selected,
			Note:        noteEntry.Text,
			Status:      model.TransactionStatusPending,
		}

		// Mock IDs
		accountID := int64(1)
		categoryID := int64(1)

		var splits []model.Split

		if typeSelect.Selected == "Expense" {
			splits = append(splits,
				model.Split{ // Asset Leg
					AccountID: accountID,
					Amount:    -amountCents,
					Currency: "USD",
					ExchangeRate: 1.0,
				},
				model.Split{ // Expense Leg
					AccountID: 2, // Mock Expense Account
					CategoryID: &categoryID,
					Amount:    amountCents,
					Currency: "USD",
					ExchangeRate: 1.0,
				},
			)
		} else {
			splits = append(splits,
				model.Split{
					AccountID: accountID,
					Amount:    amountCents,
					Currency: "USD",
					ExchangeRate: 1.0,
				},
				model.Split{
					AccountID: 2,
					CategoryID: &categoryID,
					Amount:    -amountCents,
					Currency: "USD",
					ExchangeRate: 1.0,
				},
			)
		}

		t.Splits = splits

		if err := a.Repo.CreateTransaction(t); err != nil {
			dialog.ShowError(err, a.Window)
		} else {
			a.ContentContainer.Refresh()
			// Close the modal window somehow?
			// We created a new window 'w' in ShowAddTransactionModal but don't have ref here easily.
			// Passing 'w' or using a callback would be better.
			// For this MVP refactor, let's just Refresh.
			// User has to close window manually or we change architecture.
		}
	})

	return container.NewVBox(form, saveBtn)
}

func (a *App) createSplitForm() fyne.CanvasObject {
	// Header
	dateEntry := widget.NewEntry()
	dateEntry.SetText(time.Now().Format("2006-01-02"))

	descEntry := widget.NewEntry()
	descEntry.PlaceHolder = "Payee (e.g. Supermarket)"

	sourceAccount := widget.NewSelect([]string{"Cash", "Bank", "Credit Card"}, nil)
	sourceAccount.Selected = "Cash"

	// Splits container
	splitsContainer := container.NewVBox()

	// Helper to add row
	totalLabel := widget.NewLabel("Total: $0.00")
	var splitRows []*SplitRow

	updateTotal := func() {
		var sum float64
		for _, row := range splitRows {
			v, _ := strconv.ParseFloat(row.AmountEntry.Text, 64)
			sum += v
		}
		totalLabel.SetText(fmt.Sprintf("Total: $%.2f", sum))
	}

	addSplitRow := func() {
		row := NewSplitRow(updateTotal)
		splitRows = append(splitRows, row)
		splitsContainer.Add(row.Container)
	}

	// Add initial rows
	addSplitRow()
	addSplitRow()

	addBtn := widget.NewButtonWithIcon("Add Split", theme.ContentAddIcon(), addSplitRow)

	saveBtn := widget.NewButtonWithIcon("Save Split Transaction", theme.DocumentSaveIcon(), func() {
		// Logic to save
		date, _ := time.Parse("2006-01-02", dateEntry.Text)
		t := &model.Transaction{
			Date: date,
			Description: descEntry.Text,
			Status: model.TransactionStatusPending,
		}

		var totalCents int64
		var splits []model.Split

		// 1. Process Split Rows (Expenses usually)
		for _, row := range splitRows {
			amtFloat, _ := strconv.ParseFloat(row.AmountEntry.Text, 64)
			amtCents := int64(amtFloat * 100)
			if amtCents <= 0 { continue }

			totalCents += amtCents

			// Add Expense Split
			catID := int64(1) // Mock lookup row.CategorySelect.Selected
			splits = append(splits, model.Split{
				AccountID: 2, // Mock Expense Account
				CategoryID: &catID,
				Amount: amtCents, // Debit Expense
				Currency: "USD", ExchangeRate: 1.0,
			})
		}

		// 2. Add Source Account Split (Asset/Liability)
		srcID := int64(1) // Mock sourceAccount.Selected
		splits = append(splits, model.Split{
			AccountID: srcID,
			Amount: -totalCents, // Credit Source
			Currency: "USD", ExchangeRate: 1.0,
		})

		t.Splits = splits

		if err := a.Repo.CreateTransaction(t); err != nil {
			dialog.ShowError(err, a.Window) // This dialog might need the parent window passed correctly
		} else {
			a.ContentContainer.Refresh()
		}
	})

	return container.NewVBox(
		widget.NewLabelWithStyle("Split Transaction", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewForm(
			widget.NewFormItem("Date", dateEntry),
			widget.NewFormItem("Payee", descEntry),
			widget.NewFormItem("Source", sourceAccount),
		),
		widget.NewSeparator(),
		widget.NewLabel("Category Splits:"),
		splitsContainer,
		addBtn,
		widget.NewSeparator(),
		totalLabel,
		saveBtn,
	)
}

type SplitRow struct {
	Container *fyne.Container
	CategorySelect *widget.Select
	AmountEntry *widget.Entry
}

func NewSplitRow(onChange func()) *SplitRow {
	cat := widget.NewSelect([]string{"Food", "Home", "Baby", "Ent."}, nil)
	cat.Selected = "Food"

	amt := widget.NewEntry()
	amt.PlaceHolder = "0.00"
	amt.OnChanged = func(s string) { onChange() }

	// Layout: [Category (Expand)] [Amount (Fixed)]
	// Using Grid or HBox

	row := container.NewGridWithColumns(2, cat, amt)

	return &SplitRow{
		Container: row,
		CategorySelect: cat,
		AmountEntry: amt,
	}
}
