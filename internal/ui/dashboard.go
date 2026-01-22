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
		stats = &repository.DashboardStats{} // Empty on error
	}

	// 4 Pillars
	incomeCard := createInfoCard("Income", fmt.Sprintf("$%.2f", stats.TotalIncome), color.RGBA{0, 200, 0, 255})
	expenseCard := createInfoCard("Expenses", fmt.Sprintf("$%.2f", stats.TotalExpense), color.RGBA{200, 0, 0, 255})
	assetCard := createInfoCard("Assets", fmt.Sprintf("$%.2f", stats.TotalAssets), color.RGBA{0, 0, 200, 255})
	// Net Worth
	netWorthCard := createInfoCard("Net Worth", fmt.Sprintf("$%.2f", stats.NetWorth), color.RGBA{100, 100, 100, 255})

	pillars := container.NewGridWithColumns(4, incomeCard, expenseCard, assetCard, netWorthCard)

	// Chart
	chartStats, _ := repo.GetMonthlyStats(6)
	chart := NewBarChart(chartStats)

	chartArea := container.NewVBox(
		widget.NewLabelWithStyle("Income vs Expense (6 Months)", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		container.NewCenter(chart),
	)

	// Placeholder for header, assuming it will be defined elsewhere or is a global variable
	// For now, let's define a simple header
	header := widget.NewLabelWithStyle("Dashboard", fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Monospace: true})

	return container.NewVScroll(container.NewVBox(
		header,
		pillars,
		widget.NewSeparator(),
		chartArea,
	))
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
