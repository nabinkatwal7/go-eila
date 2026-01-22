package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/nabinkatwal7/go-eila/internal/repository"
)

type App struct {
	FyneApp fyne.App
	Window  fyne.Window
	Repo    *repository.Repository

	ContentContainer *fyne.Container
}

func NewApp(repo *repository.Repository) *App {
	a := app.New()
	w := a.NewWindow("MyTrack")
	w.Resize(fyne.NewSize(1000, 700))

	myApp := &App{
		FyneApp: a,
		Window:  w,
		Repo:    repo,
	}

	myApp.setupUI()
	return myApp
}

func (a *App) setupUI() {
	// sidebar
	sidebar := a.createSidebar()

	// content area (initial view)
	a.ContentContainer = container.NewMax(NewDashboard(a.Repo))

	// main layout
	split := container.NewHSplit(sidebar, a.ContentContainer)
	split.SetOffset(0.2)

	a.Window.SetContent(split)
}

func (a *App) createSidebar() fyne.CanvasObject {
	// Add Button
	addBtn := widget.NewButton("+ Add New", func() {
		a.ShowAddTransactionModal()
	})
	addBtn.Importance = widget.HighImportance

	// Navigation buttons
	dashBtn := widget.NewButton("Dashboard", func() {
		a.ContentContainer.Objects = []fyne.CanvasObject{NewDashboard(a.Repo)}
		a.ContentContainer.Refresh()
	})
	transBtn := widget.NewButton("Transactions", func() {
		a.ContentContainer.Objects = []fyne.CanvasObject{NewTransactionsView(a.Repo)}
		a.ContentContainer.Refresh()
	})
	// Accounts
	accountsBtn := widget.NewButton("Accounts", func() {
		// Show Accounts View
		a.ContentContainer.Objects = []fyne.CanvasObject{NewAccountsView(a.Repo, a)}
		a.ContentContainer.Refresh()
	})

	// Layout
	return container.NewVBox(
		widget.NewLabelWithStyle("MyTrack", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewSeparator(),
		addBtn,
		widget.NewSeparator(),
		dashBtn,
		transBtn,
		accountsBtn,
	)
}

func (a *App) Run() {
	a.Window.ShowAndRun()
}
