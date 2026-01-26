package ui

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/nabinkatwal7/go-eila/internal/model"
	"github.com/nabinkatwal7/go-eila/internal/repository"
)

func NewTransactionsView(repo *repository.Repository, app *App) fyne.CanvasObject {
	header := widget.NewLabelWithStyle("Transactions", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	// Search and filter controls
	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("Search by description...")

	startDateEntry := widget.NewEntry()
	startDateEntry.SetPlaceHolder("Start date (YYYY-MM-DD)")
	endDateEntry := widget.NewEntry()
	endDateEntry.SetPlaceHolder("End date (YYYY-MM-DD)")

	minAmountEntry := widget.NewEntry()
	minAmountEntry.SetPlaceHolder("Min amount")
	maxAmountEntry := widget.NewEntry()
	maxAmountEntry.SetPlaceHolder("Max amount")

	searchBtn := widget.NewButton("Search", nil)
	clearBtn := widget.NewButton("Clear", nil)

	// Table to display transactions
	var transactions []model.Transaction
	var table *widget.Table

	refreshTable := func() {
		var err error
		var startDate, endDate *time.Time
		var minAmount, maxAmount *float64

		// Parse filters
		if startDateEntry.Text != "" {
			if d, err := time.Parse("2006-01-02", startDateEntry.Text); err == nil {
				startDate = &d
			}
		}
		if endDateEntry.Text != "" {
			if d, err := time.Parse("2006-01-02", endDateEntry.Text); err == nil {
				endDate = &d
			}
		}
		if minAmountEntry.Text != "" {
			if amt, err := ValidateAmount(minAmountEntry.Text); err == nil {
				minAmount = &amt
			}
		}
		if maxAmountEntry.Text != "" {
			if amt, err := ValidateAmount(maxAmountEntry.Text); err == nil {
				maxAmount = &amt
			}
		}

		// Perform search
		if searchEntry.Text != "" || startDate != nil || endDate != nil || minAmount != nil || maxAmount != nil {
			transactions, err = repo.SearchTransactions(searchEntry.Text, startDate, endDate, minAmount, maxAmount, 1000)
		} else {
			transactions, err = repo.GetRecentTransactions(100)
		}

		if err != nil {
			dialog.ShowError(err, app.Window)
			return
		}

		table.Refresh()
	}

	searchBtn.OnTapped = refreshTable
	clearBtn.OnTapped = func() {
		searchEntry.SetText("")
		startDateEntry.SetText("")
		endDateEntry.SetText("")
		minAmountEntry.SetText("")
		maxAmountEntry.SetText("")
		refreshTable()
	}

	// Initial load
	transactions, _ = repo.GetRecentTransactions(100)

	// Get categories for display
	categories, _ := repo.GetAllCategories()
	categoryIDToName := make(map[int64]string)
	for _, cat := range categories {
		categoryIDToName[cat.ID] = cat.Name
	}

	tableHeader := []string{"Date", "Payee", "Amount", "Category", "Actions"}
	table = widget.NewTable(
		func() (int, int) {
			return len(transactions), len(tableHeader)
		},
		func() fyne.CanvasObject {
			if len(tableHeader) > 4 {
				// Actions column
				return container.NewHBox(
					widget.NewButton("Edit", nil),
					widget.NewButton("Delete", nil),
				)
			}
			return widget.NewLabelWithStyle("Cell", fyne.TextAlignLeading, fyne.TextStyle{})
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			if i.Row >= len(transactions) {
				return
			}

			t := transactions[i.Row]

			if i.Col == 4 { // Actions column
				box := o.(*fyne.Container)
				editBtn := box.Objects[0].(*widget.Button)
				deleteBtn := box.Objects[1].(*widget.Button)

				editBtn.OnTapped = func() {
					app.ShowEditTransactionModal(t.ID)
				}

				deleteBtn.OnTapped = func() {
					dialog.ShowConfirm("Delete Transaction",
						fmt.Sprintf("Are you sure you want to delete transaction '%s'?", t.Description),
						func(confirmed bool) {
							if confirmed {
								if err := repo.DeleteTransaction(t.ID); err != nil {
									dialog.ShowError(err, app.Window)
								} else {
									dialog.ShowInformation("Success", "Transaction deleted", app.Window)
									refreshTable()
									app.ContentContainer.Refresh()
								}
							}
						}, app.Window)
				}
				return
			}

			label := o.(*widget.Label)
			switch i.Col {
			case 0: // Date
				label.SetText(t.Date.Format("2006-01-02"))
			case 1: // Payee
				label.SetText(t.Description)
			case 2: // Amount
				var amt int64
				for _, s := range t.Splits {
					if s.Amount > 0 {
						amt += s.Amount
					}
				}
				label.SetText(fmt.Sprintf("$%.2f", float64(amt)/100.0))
				label.TextStyle = fyne.TextStyle{Bold: true}
				label.Alignment = fyne.TextAlignTrailing
			case 3: // Category
				var catName string
				for _, s := range t.Splits {
					if s.CategoryID != nil {
						if name, ok := categoryIDToName[*s.CategoryID]; ok {
							catName = name
							break
						}
					}
				}
				label.SetText(catName)
			}
		},
	)

	table.SetColumnWidth(0, 100) // Date
	table.SetColumnWidth(1, 250) // Payee
	table.SetColumnWidth(2, 100) // Amount
	table.SetColumnWidth(3, 120) // Category
	table.SetColumnWidth(4, 150) // Actions

	// Filter controls
	filterBox := container.NewVBox(
		container.NewHBox(
			searchEntry,
			searchBtn,
			clearBtn,
		),
		container.NewHBox(
			widget.NewLabel("Date Range:"),
			startDateEntry,
			widget.NewLabel("-"),
			endDateEntry,
		),
		container.NewHBox(
			widget.NewLabel("Amount Range:"),
			minAmountEntry,
			widget.NewLabel("-"),
			maxAmountEntry,
		),
	)

	return container.NewBorder(
		container.NewVBox(header, filterBox),
		nil, nil, nil,
		table,
	)
}
