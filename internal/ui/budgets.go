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

	// Add Budget Button
	addBtn := widget.NewButton("Set/Update Budget", func() {
		showAddBudgetModal(repo, a)
	})

	// Fetch Progress
	now := time.Now()
	progress, err := repo.GetBudgetsWithProgress(int(now.Month()), now.Year())

	content := container.NewVBox()

	if err != nil {
		content.Add(widget.NewLabelWithStyle("Error loading budgets: "+err.Error(), fyne.TextAlignLeading, fyne.TextStyle{TabWidth: 2}))
	} else {
		for _, p := range progress {
			// Row: Name --- Spending / Limit
			// Progress Bar
			info := fmt.Sprintf("%s: $%.2f / $%.2f", p.CategoryName, p.Spent, p.Budgeted)
			label := widget.NewLabelWithStyle(info, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

			bar := widget.NewProgressBar()
			bar.Value = p.Percent
			// Visual indication for over budget
			if p.Spent > p.Budgeted {
				label.TextStyle = fyne.TextStyle{Bold: true, Italic: true} // Just style for now, color needs canvas
				// Format to show overage
				label.SetText(fmt.Sprintf("%s: $%.2f / $%.2f (OVER BUDGET!)", p.CategoryName, p.Spent, p.Budgeted))
			}

			row := container.NewVBox(label, bar)
			content.Add(row)
		}

		if len(progress) == 0 {
			content.Add(widget.NewLabel("No budgets set for this month."))
		}
	}

	return container.NewBorder(container.NewHBox(header, addBtn), nil, nil, nil, container.NewVScroll(content))
}

func showAddBudgetModal(repo *repository.Repository, a *App) {
	// Load real categories
	categories, err := repo.GetAllCategories()
	if err != nil {
		dialog.ShowError(err, a.Window)
		return
	}

	if len(categories) == 0 {
		dialog.ShowInformation("No Categories", "Please create some categories first (not yet implemented in UI).", a.Window)
		return
	}

	categoryNames := make([]string, len(categories))
	categoryNameToID := make(map[string]int64)
	for i, c := range categories {
		categoryNames[i] = c.Name
		categoryNameToID[c.Name] = c.ID
	}

	catSelect := widget.NewSelect(categoryNames, nil)
	if len(categoryNames) > 0 {
		catSelect.Selected = categoryNames[0]
	}

	amtEntry := widget.NewEntry()
	amtEntry.PlaceHolder = "Monthly Limit (e.g. 500.00)"

	items := []*widget.FormItem{
		widget.NewFormItem("Category", catSelect),
		widget.NewFormItem("Limit ($)", amtEntry),
	}

	d := dialog.NewForm("Set Budget", "Save", "Cancel", items, func(confirm bool) {
		if confirm {
			// Validate inputs
			amt, err := ValidateAmount(amtEntry.Text)
			if err != nil {
				dialog.ShowError(err, a.Window)
				return // Fyne dialog callback doesn't support easy 'prevent close', checking here is tricky
				// Ideal UX: Validation before submission or reopen dialog.
				// For now: Show error dialog.
			}

			catID := categoryNameToID[catSelect.Selected]

			b := &model.Budget{
				CategoryID: catID,
				Amount: int64(amt * 100),
				Period: "Monthly",
			}

			// Note: CreateBudget in repo currently just inserts.
			// Ideally should upsert (update if exists for category).
			// Assuming repo handles it or we just add new row.
			if err := repo.CreateBudget(b); err != nil {
				dialog.ShowError(err, a.Window)
			} else {
				a.ContentContainer.Refresh()
				dialog.ShowInformation("Success", "Budget set successfully.", a.Window)
			}
		}
	}, a.Window)

	d.Resize(fyne.NewSize(400, 200))
	d.Show()
}
