package ui

import (
	"fmt"
	"math/rand"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/nabinkatwal7/go-eila/internal/repository"
)

func NewInvoicingView(repo *repository.Repository) fyne.CanvasObject {
	header := widget.NewLabelWithStyle("Invoicing & Business", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	// 1. Select Client (Mocked)
	client := widget.NewEntry()
	client.PlaceHolder = "Client Name"

	// 2. Select Items (Mocked one item)
	desc := widget.NewEntry()
	desc.PlaceHolder = "Service Description"
	amount := widget.NewEntry()
	amount.PlaceHolder = "Amount"

	preview := widget.NewMultiLineEntry()
	preview.Disable()
	preview.TextStyle = fyne.TextStyle{Monospace: true}
	preview.SetMinRowsVisible(10)

	btn := widget.NewButton("Generate Invoice", func() {
		// Mock ID
		id := rand.Intn(10000)
		date := time.Now().Format("2006-01-02")

		inv := fmt.Sprintf(`IPV INVOICE #%d
Date: %s
Client: %s

----------------------------------------
Description             Amount
----------------------------------------
%s                      $%s
----------------------------------------
Total:                  $%s

Thank you for your business!
Payment due within 30 days.
`, id, date, client.Text, desc.Text, amount.Text, amount.Text)

		preview.SetText(inv)
	})

	return container.NewVScroll(container.NewVBox(
		header,
		widget.NewLabel("Create Quick Invoice"),
		client,
		desc,
		amount,
		btn,
		widget.NewSeparator(),
		widget.NewLabel("Preview:"),
		preview,
	))
}
