package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/nabinkatwal7/go-eila/internal/repository"
)

func NewTransactionsView(repo *repository.Repository) fyne.CanvasObject {
	header := widget.NewLabelWithStyle("Transactions (Double-Entry)", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

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
				widget.NewLabel("Desc"),
				widget.NewLabel("Amount"),
			)
		},
		func(i int, o fyne.CanvasObject) {
			t := txs[i]
			box := o.(*fyne.Container)
			box.Objects[0].(*widget.Label).SetText(t.Date.Format("2006-01-02"))
			box.Objects[1].(*widget.Label).SetText(t.Description)

			// Amount: Sum of positive splits? Or just the first split?
			// Display the first split that is NOT the main account... tricky.
			// Just display sum of positive splits for now (Total Debit)
			var sum int64
			for _, s := range t.Splits {
				if s.Amount > 0 {
					sum += s.Amount
				}
			}
			box.Objects[2].(*widget.Label).SetText(fmt.Sprintf("$%.2f", float64(sum)/100.0))
		},
	)

	return container.NewBorder(header, nil, nil, nil, list)
}
