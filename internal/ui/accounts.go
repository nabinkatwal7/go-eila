package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/nabinkatwal7/go-eila/internal/repository"
)

func NewAccountsView(repo *repository.Repository, a *App) fyne.CanvasObject {
	header := widget.NewLabelWithStyle("Accounts", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	// Add Button
	addBtn := widget.NewButton("+ New Account", func() {
		a.ShowCreateAccountModal()
	})

	// List
	accounts, err := repo.GetAllAccounts()
	if err != nil {
		return widget.NewLabel("Error loading accounts: " + err.Error())
	}

	list := widget.NewList(
		func() int {
			return len(accounts)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabel("Name"),
				widget.NewLabel("Type"),
				widget.NewLabel("Currency"),
				widget.NewLabel("Balance"),
			)
		},
		func(i int, o fyne.CanvasObject) {
			ac := accounts[i]
			box := o.(*fyne.Container)
			box.Objects[0].(*widget.Label).SetText(ac.Name)
			box.Objects[1].(*widget.Label).SetText(string(ac.Type))
			box.Objects[2].(*widget.Label).SetText(ac.Currency)

			// Balance check (Need separate query or eager load)
			bal, _ := repo.GetAccountBalance(ac.ID)
			box.Objects[3].(*widget.Label).SetText(fmt.Sprintf("%.2f", bal))
		},
	)

	return container.NewBorder(container.NewHBox(header, addBtn), nil, nil, nil, list)
}
