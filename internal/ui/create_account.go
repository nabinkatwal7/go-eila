package ui

import (
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"github.com/nabinkatwal7/go-eila/internal/model"
)

func (a *App) ShowCreateAccountModal() {
	nameEntry := widget.NewEntry()
	nameEntry.PlaceHolder = "Account Name (e.g., 'Travel Fund')"

	typeSelect := widget.NewSelect([]string{
		string(model.AccountTypeCash),
		string(model.AccountTypeBank),
		string(model.AccountTypeCard),
		string(model.AccountTypeInvest),
		string(model.AccountTypeEquity),
	}, nil)
	typeSelect.Selected = string(model.AccountTypeBank)

	// Currency Select
	currencySelect := widget.NewSelect([]string{"USD", "EUR", "GBP", "NPR", "JPY"}, nil)
	currencySelect.Selected = "USD"

	items := []*widget.FormItem{
		widget.NewFormItem("Name", nameEntry),
		widget.NewFormItem("Type", typeSelect),
		widget.NewFormItem("Currency", currencySelect),
	}

	dialog.ShowForm("New Account", "Create", "Cancel", items, func(confirm bool) {
		if confirm {
			acc := &model.Account{
				Name: nameEntry.Text,
				Type: model.AccountType(typeSelect.Selected),
				Currency: currencySelect.Selected,
			}
			if err := a.Repo.CreateAccount(acc); err != nil {
				dialog.ShowError(err, a.Window)
			} else {
				a.ContentContainer.Refresh()
			}
		}
	}, a.Window)
}
