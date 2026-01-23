package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/nabinkatwal7/go-eila/internal/model"
	"github.com/nabinkatwal7/go-eila/internal/repository"
)

func NewInvestmentView(repo *repository.Repository) fyne.CanvasObject {
	header := widget.NewLabelWithStyle("Investment Portfolio", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	// Get Investment Accounts
	accounts, err := repo.GetAllAccounts()
	if err != nil {
		return widget.NewLabel("Error: " + err.Error())
	}

	content := container.NewVBox()

	totalValue := 0.0

	for _, a := range accounts {
		if a.Type == model.AccountTypeInvest {
			bal, _ := repo.GetAccountBalance(a.ID)
			totalValue += bal

			// Mocking Performance for now as we don't have separate 'Cost Basis' tracking in splits yet
			// In real app, we'd query transfers vs income.
			// Let's assume 10% gain for demo visual.
			gain := bal * 0.10
			cost := bal - gain

			card := createInvestCard(a.Name, bal, cost, gain)
			content.Add(card)
		}
	}

	summary := widget.NewLabelWithStyle(fmt.Sprintf("Total Portfolio Value: $%.2f", totalValue), fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	return container.NewVScroll(container.NewVBox(
		header,
		summary,
		widget.NewSeparator(),
		content,
	))
}

func createInvestCard(name string, value, cost, gain float64) fyne.CanvasObject {
	valStr := fmt.Sprintf("Value: $%.2f", value)
	costStr := fmt.Sprintf("Cost: $%.2f", cost)
	gainStr := fmt.Sprintf("Unrealized Gain: $%.2f (+10%%)", gain)

	return widget.NewCard(name, valStr, container.NewVBox(
		widget.NewLabel(costStr),
		widget.NewLabelWithStyle(gainStr, fyne.TextAlignLeading, fyne.TextStyle{}),
	))
}
