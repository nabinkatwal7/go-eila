package ui

import (
	"fmt"
	"image/color"
	"time"

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

	// Charts
	chartStats, _ := repo.GetMonthlyStats(6)
	incomeExpenseChart := NewBarChart(chartStats)

	// Net Worth Progression
	netWorthHistory, _ := repo.GetNetWorthHistory(12)
	netWorthChart := NewLineChart(netWorthHistory)

	// Category Breakdown (last 30 days)
	now := time.Now()
	thirtyDaysAgo := now.AddDate(0, 0, -30)
	categoryBreakdown, _ := repo.GetCategoryBreakdown(&thirtyDaysAgo, &now)
	categoryChart := NewPieChart(categoryBreakdown)

	incomeExpenseArea := container.NewVBox(
		widget.NewLabelWithStyle("Income vs Expense (6 Months)", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		container.NewCenter(incomeExpenseChart),
	)

	netWorthArea := container.NewVBox(
		widget.NewLabelWithStyle("Net Worth Progression (12 Months)", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		container.NewCenter(netWorthChart),
	)

	categoryArea := container.NewVBox(
		widget.NewLabelWithStyle("Category Breakdown (Last 30 Days)", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		container.NewCenter(categoryChart),
	)

	header := widget.NewLabelWithStyle("Dashboard", fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Monospace: true})

	return container.NewVScroll(container.NewVBox(
		header,
		pillars,
		widget.NewSeparator(),
		incomeExpenseArea,
		widget.NewSeparator(),
		netWorthArea,
		widget.NewSeparator(),
		categoryArea,
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
