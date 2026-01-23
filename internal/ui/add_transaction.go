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
	// Create window first so we can pass it to form creators
	w := a.FyneApp.NewWindow("Add Transaction")

	// We will use Tabs for "Simple" vs "Split" modes
	simpleContent := a.createSimpleForm(w)
	splitContent := a.createSplitForm(w)

	tabs := container.NewAppTabs(
		container.NewTabItem("Simple", simpleContent),
		container.NewTabItem("Split", splitContent),
	)

	w.Resize(fyne.NewSize(500, 600))
	w.SetContent(container.NewPadded(tabs))
	w.Show()
}

func (a *App) createSimpleForm(w fyne.Window) fyne.CanvasObject {
	// Load real accounts and categories from database
	accounts, err := a.Repo.GetAllAccounts()
	if err != nil {
		return widget.NewLabel("Error loading accounts: " + err.Error())
	}

	categories, err := a.Repo.GetAllCategories()
	if err != nil {
		return widget.NewLabel("Error loading categories: " + err.Error())
	}

	// Filter accounts by type for dropdowns
	var assetAccounts []model.Account
	for _, acc := range accounts {
		if acc.Type == model.AccountTypeCash || acc.Type == model.AccountTypeBank ||
		   acc.Type == model.AccountTypeCard || acc.Type == model.AccountTypeInvest {
			assetAccounts = append(assetAccounts, acc)
		}
	}

	// Create dropdown options and ID maps
	accountNames := make([]string, len(assetAccounts))
	accountNameToID := make(map[string]int64)
	for i, acc := range assetAccounts {
		accountNames[i] = acc.Name
		accountNameToID[acc.Name] = acc.ID
	}

	categoryNames := make([]string, len(categories))
	categoryNameToID := make(map[string]int64)
	for i, cat := range categories {
		categoryNames[i] = cat.Name
		categoryNameToID[cat.Name] = cat.ID
	}

	// Inputs
	amountEntry := widget.NewEntry()
	amountEntry.PlaceHolder = "Amount"

	noteEntry := widget.NewEntry()
	noteEntry.PlaceHolder = "Note (e.g., Lunch)"

	typeSelect := widget.NewSelect([]string{"Expense", "Income"}, nil)
	typeSelect.Selected = "Expense"

	dateEntry := widget.NewEntry()
	dateEntry.SetText(time.Now().Format("2006-01-02"))

	accountSelect := widget.NewSelect(accountNames, nil)
	if len(accountNames) > 0 {
		accountSelect.Selected = accountNames[0]
	}

	categorySelect := widget.NewSelect(categoryNames, nil)
	if len(categoryNames) > 0 {
		categorySelect.Selected = categoryNames[0]
	}

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

		// Get actual IDs from selections
		accountID := accountNameToID[accountSelect.Selected]
		categoryID := categoryNameToID[categorySelect.Selected]

		// Rule Engine Hook
		desc := categorySelect.Selected
		if noteEntry.Text != "" {
			desc = noteEntry.Text
		}

		enrichedPayee, enrichedCatID, enrichedNote := a.Repo.EnrichTransaction(desc)

		finalDesc := desc
		if enrichedPayee != "" && enrichedPayee != desc {
			finalDesc = enrichedPayee
		}

		finalNote := noteEntry.Text
		if enrichedNote != "" {
			finalNote = enrichedNote
		}

		// If Rule found a category, use it
		if enrichedCatID != nil {
			categoryID = *enrichedCatID
		}

		t := &model.Transaction{
			Date:        date,
			Description: finalDesc,
			Note:        finalNote,
			Status:      model.TransactionStatusPending,
		}

		// Find Expense/Income account IDs
		var expenseAccountID, incomeAccountID int64
		for _, acc := range accounts {
			if acc.Type == model.AccountTypeExpense {
				expenseAccountID = acc.ID
			}
			if acc.Type == model.AccountTypeIncome {
				incomeAccountID = acc.ID
			}
		}

		var splits []model.Split

		if typeSelect.Selected == "Expense" {
			splits = append(splits,
				model.Split{ // Asset Leg (decrease)
					AccountID: accountID,
					Amount:    -amountCents,
					Currency: "USD",
					ExchangeRate: 1.0,
				},
				model.Split{ // Expense Leg (increase)
					AccountID: expenseAccountID,
					CategoryID: &categoryID,
					Amount:    amountCents,
					Currency: "USD",
					ExchangeRate: 1.0,
				},
			)
		} else {
			splits = append(splits,
				model.Split{ // Asset Leg (increase)
					AccountID: accountID,
					Amount:    amountCents,
					Currency: "USD",
					ExchangeRate: 1.0,
				},
				model.Split{ // Income Leg (decrease)
					AccountID: incomeAccountID,
					CategoryID: &categoryID,
					Amount:    -amountCents,
					Currency: "USD",
					ExchangeRate: 1.0,
				},
			)
		}

		t.Splits = splits

		if err := a.Repo.CreateTransaction(t); err != nil {
			dialog.ShowError(err, w)
		} else {
			a.ContentContainer.Refresh()
			w.Close() // Close the modal window
			dialog.ShowInformation("Success", "Transaction added successfully!", a.Window)
		}
	})

	return container.NewVBox(form, saveBtn)
}

func (a *App) createSplitForm(w fyne.Window) fyne.CanvasObject {
	// Load real accounts and categories from database
	accounts, err := a.Repo.GetAllAccounts()
	if err != nil {
		return widget.NewLabel("Error loading accounts: " + err.Error())
	}

	categories, err := a.Repo.GetAllCategories()
	if err != nil {
		return widget.NewLabel("Error loading categories: " + err.Error())
	}

	// Filter asset accounts
	var assetAccounts []model.Account
	for _, acc := range accounts {
		if acc.Type == model.AccountTypeCash || acc.Type == model.AccountTypeBank ||
		   acc.Type == model.AccountTypeCard || acc.Type == model.AccountTypeInvest {
			assetAccounts = append(assetAccounts, acc)
		}
	}

	// Create dropdown options and ID maps
	accountNames := make([]string, len(assetAccounts))
	accountNameToID := make(map[string]int64)
	for i, acc := range assetAccounts {
		accountNames[i] = acc.Name
		accountNameToID[acc.Name] = acc.ID
	}

	categoryNames := make([]string, len(categories))
	categoryNameToID := make(map[string]int64)
	for i, cat := range categories {
		categoryNames[i] = cat.Name
		categoryNameToID[cat.Name] = cat.ID
	}

	// Header
	dateEntry := widget.NewEntry()
	dateEntry.SetText(time.Now().Format("2006-01-02"))

	descEntry := widget.NewEntry()
	descEntry.PlaceHolder = "Payee (e.g. Supermarket)"

	sourceAccount := widget.NewSelect(accountNames, nil)
	if len(accountNames) > 0 {
		sourceAccount.Selected = accountNames[0]
	}

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
		row := NewSplitRow(categoryNames, updateTotal)
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

		// Get source account ID
		srcID := accountNameToID[sourceAccount.Selected]

		// Find Expense account ID
		var expenseAccountID int64
		for _, acc := range accounts {
			if acc.Type == model.AccountTypeExpense {
				expenseAccountID = acc.ID
				break
			}
		}

		var totalCents int64
		var splits []model.Split

		// 1. Process Split Rows (Expenses usually)
		for _, row := range splitRows {
			amtFloat, _ := strconv.ParseFloat(row.AmountEntry.Text, 64)
			amtCents := int64(amtFloat * 100)
			if amtCents <= 0 { continue }

			totalCents += amtCents

			// Get category ID from selection
			catID := categoryNameToID[row.CategorySelect.Selected]
			splits = append(splits, model.Split{
				AccountID: expenseAccountID,
				CategoryID: &catID,
				Amount: amtCents, // Debit Expense
				Currency: "USD", ExchangeRate: 1.0,
			})
		}

		// 2. Add Source Account Split (Asset/Liability)
		splits = append(splits, model.Split{
			AccountID: srcID,
			Amount: -totalCents, // Credit Source
			Currency: "USD", ExchangeRate: 1.0,
		})

		t.Splits = splits

		if err := a.Repo.CreateTransaction(t); err != nil {
			dialog.ShowError(err, w)
		} else {
			a.ContentContainer.Refresh()
			w.Close() // Close the modal window
			dialog.ShowInformation("Success", "Split transaction added successfully!", a.Window)
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

func NewSplitRow(categoryNames []string, onChange func()) *SplitRow {
	cat := widget.NewSelect(categoryNames, nil)
	if len(categoryNames) > 0 {
		cat.Selected = categoryNames[0]
	}

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
