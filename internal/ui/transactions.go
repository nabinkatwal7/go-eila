package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/nabinkatwal7/go-eila/internal/repository"
)

func NewTransactionsView(repo *repository.Repository) fyne.CanvasObject {
	// Header
	header := widget.NewLabelWithStyle("Transactions", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	// Table (List of transactions)
	// For now, dragging in recent transactions
	txs, err := repo.GetRecentTransactions(50)
	if err != nil {
		return widget.NewLabel("Error loading transactions: " + err.Error())
	}

	list := widget.NewList(
		func() int {
			return len(txs)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabel("Date"),
				widget.NewLabel("Note"),
				widget.NewLabel("Amount"),
			)
		},
		func(i int, o fyne.CanvasObject) {
			t := txs[i]
			box := o.(*fyne.Container)
			box.Objects[0].(*widget.Label).SetText(t.Date.Format("2006-01-02"))
			box.Objects[1].(*widget.Label).SetText(t.Note)
			box.Objects[2].(*widget.Label).SetText(fmt.Sprintf("%.2f", t.Amount))
		},
	)

	return container.NewBorder(header, nil, nil, nil, list)
}
