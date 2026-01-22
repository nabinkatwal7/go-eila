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

func NewBudgetsView(repo *repository.Repository, a *App) fyne.CanvasObject {
	header := widget.NewLabelWithStyle("Budgets (This Month)", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	// Add Budget Button (Quick mock)
	addBtn := widget.NewButton("Set Budget", func() {
		showAddBudgetModal(repo, a)
	})

	// Fetch Progress
	now := time.Now()
	progress, err := repo.GetBudgetsWithProgress(int(now.Month()), now.Year())
	if err != nil {
		return widget.NewLabel("Error: " + err.Error())
	}

	content := container.NewVBox()

	for _, p := range progress {
		// Row: Name --- Spending / Limit
		// Progress Bar

		info := fmt.Sprintf("%s: $%.2f / $%.2f", p.CategoryName, p.Spent, p.Budgeted)
		label := widget.NewLabelWithStyle(info, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

		bar := widget.NewProgressBar()
		bar.Value = p.Percent
		// Customize bar color based heavily on usage? Fyne default is blue.

		row := container.NewVBox(label, bar)
		content.Add(row)
	}

	if len(progress) == 0 {
		content.Add(widget.NewLabel("No budgets set for this month."))
	}

	return container.NewBorder(container.NewHBox(header, addBtn), nil, nil, nil, container.NewVScroll(content))
}

func showAddBudgetModal(repo *repository.Repository, a *App) {
	// Simple modal to set budget on a category
	// Mocks category selection
	catSelect := widget.NewSelect([]string{"Food", "Transport"}, nil)
	catSelect.Selected = "Food"

	amtEntry := widget.NewEntry()
	amtEntry.PlaceHolder = "Monthly Limit"

	items := []*widget.FormItem{
		widget.NewFormItem("Category", catSelect),
		widget.NewFormItem("Limit", amtEntry),
	}

	dialog.ShowForm("Set Budget", "Save", "Cancel", items, func(confirm bool) {
		if confirm {
			// Save
			// Mock Category ID lookup
			catID := int64(1)

			// Parse amount
			var amt float64
			fmt.Sscanf(amtEntry.Text, "%f", &amt)

			b := &model.Budget{
				CategoryID: catID,
				Amount: int64(amt * 100),
				Period: "Monthly",
			}
			repo.CreateBudget(b)
			a.ContentContainer.Refresh()
		}
	}, a.Window)
}
