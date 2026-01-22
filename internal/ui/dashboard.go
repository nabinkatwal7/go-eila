package ui

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/nabinkatwal7/go-eila/internal/repository"
)

func NewDashboard(repo *repository.Repository) fyne.CanvasObject {
	stats, err := repo.GetDashboardStats()
	if err != nil {
		// handle error gracefully, maybe show zero stats
		stats = &repository.DashboardStats{}
	}

	// 4 Pillars
	incomeCard := createInfoCard("Income", fmt.Sprintf("%.2f", stats.TotalIncome), color.RGBA{0, 200, 0, 255})
	expenseCard := createInfoCard("Expenses", fmt.Sprintf("%.2f", stats.TotalExpense), color.RGBA{200, 0, 0, 255})
	assetCard := createInfoCard("Assets", fmt.Sprintf("%.2f", stats.TotalAssets), color.RGBA{0, 0, 200, 255})
	liabilityCard := createInfoCard("Liabilities", fmt.Sprintf("%.2f", stats.TotalLiability), color.RGBA{200, 100, 0, 255})

	// Net Worth
	netWorthLabel := widget.NewLabelWithStyle(fmt.Sprintf("Net Worth: $%.2f", stats.NetWorth), fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	// Grid for cards
	cardsGrid := container.NewGridWithColumns(2,
		incomeCard, expenseCard,
		assetCard, liabilityCard,
	)

	return container.NewVBox(
		widget.NewLabelWithStyle("Dashboard", fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Monospace: true}),
		widget.NewSeparator(),
		cardsGrid,
		widget.NewSeparator(),
		netWorthLabel,
		widget.NewLabel("Recent Activity (Pending Implementation)"),
	)
}

func createInfoCard(title, amount string, c color.Color) fyne.CanvasObject {
	rect := canvas.NewRectangle(color.Transparent)
	rect.StrokeColor = c
	rect.StrokeWidth = 2
	rect.SetMinSize(fyne.NewSize(150, 80))

	titleLabel := widget.NewLabelWithStyle(title, fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	amountLabel := widget.NewLabelWithStyle("$"+amount, fyne.TextAlignCenter, fyne.TextStyle{})

	content := container.NewVBox(titleLabel, amountLabel)

	return container.NewStack(rect, content)
}
