package ui

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type Command struct {
	Title       string
	Description string
	Action      func()
}

func (a *App) SetupCommandPalette() {
	// commands list
	commands := []Command{
		{"Add Transaction", "Open the transaction creation modal", a.ShowAddTransactionModal},
		{"Add Account", "Create a new financial account", a.ShowCreateAccountModal},
		{"Go to Dashboard", "View financial overview", func() { a.ContentContainer.Objects = []fyne.CanvasObject{NewDashboard(a.Repo)}; a.ContentContainer.Refresh() }},
		{"Go to Transactions", "View transaction history", func() { a.ContentContainer.Objects = []fyne.CanvasObject{NewTransactionsView(a.Repo)}; a.ContentContainer.Refresh() }},
		{"Go to Accounts", "Manage accounts", func() { a.ContentContainer.Objects = []fyne.CanvasObject{NewAccountsView(a.Repo, a)}; a.ContentContainer.Refresh() }},
		{"Go to Budgets", "Manage spending limits", func() { a.ContentContainer.Objects = []fyne.CanvasObject{NewBudgetsView(a.Repo, a)}; a.ContentContainer.Refresh() }},
		{"Go to Recurring", "View detected subscriptions", func() { a.ContentContainer.Objects = []fyne.CanvasObject{NewRecurringView(a.Repo)}; a.ContentContainer.Refresh() }},
		{"Go to Alerts", "View spending anomalies", func() { a.ContentContainer.Objects = []fyne.CanvasObject{NewAnomaliesView(a.Repo)}; a.ContentContainer.Refresh() }},
		{"Go to Settings", "Backup and Data options", func() { a.ContentContainer.Objects = []fyne.CanvasObject{NewSettingsView(a.Repo, a.Window)}; a.ContentContainer.Refresh() }},
		{"Go to Forecast", "Project future net worth", func() { a.ContentContainer.Objects = []fyne.CanvasObject{NewForecastView(a.Repo)}; a.ContentContainer.Refresh() }},
	}

	// Register Shortcut
	ctrlK := &desktop.CustomShortcut{KeyName: fyne.KeyK, Modifier: fyne.KeyModifierControl}
	a.Window.Canvas().AddShortcut(ctrlK, func(shortcut fyne.Shortcut) {
		a.showPalette(commands)
	})
}

func (a *App) showPalette(commands []Command) {
	// Create a modal dialog

	entry := widget.NewEntry()
	entry.PlaceHolder = "Type a command..."

	listData := commands // Filtered copy

	list := widget.NewList(
		func() int { return len(listData) },
		func() fyne.CanvasObject {
			return container.NewVBox(
				widget.NewLabelWithStyle("Title", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				widget.NewLabel("Description"),
			)
		},
		func(i int, o fyne.CanvasObject) {
			box := o.(*fyne.Container)
			box.Objects[0].(*widget.Label).SetText(listData[i].Title)
			box.Objects[1].(*widget.Label).SetText(listData[i].Description)
		},
	)

	// Modal window? Or Overlay?
	// Fyne Overlay is lightweight.
	// But List needs size.

	popupContent := container.NewBorder(
		container.NewPadded(entry),
		nil, nil, nil,
		list,
	)

	// Create a PopUp
	var popup *widget.PopUp

	execute := func(cmd Command) {
		popup.Hide()
		cmd.Action()
	}

	list.OnSelected = func(id int) {
		execute(listData[id])
	}

	entry.OnChanged = func(s string) {
		s = strings.ToLower(s)
		var filtered []Command
		for _, c := range commands {
			if strings.Contains(strings.ToLower(c.Title), s) || strings.Contains(strings.ToLower(c.Description), s) {
				filtered = append(filtered, c)
			}
		}
		listData = filtered
		list.Refresh()

		// Auto select first?
		if len(listData) > 0 {
			list.Select(0) // Visual select
		}
	}

	// Handle Enter key in Entry to execute first item
	// Fyne Entry doesn't have OnSubmitted easily for this unless we subclass.
	// We can trust user checking list.

	popup = widget.NewModalPopUp(container.NewGridWithColumns(1,
		container.NewPadded(widget.NewLabelWithStyle("Command Palette", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})),
		container.NewPadded(popupContent),
		widget.NewButton("Close", func() { popup.Hide() }),
	), a.Window.Canvas())

	// Resize popup
	popup.Resize(fyne.NewSize(400, 300))
	popup.Show()

	a.Window.Canvas().Focus(entry)
}
