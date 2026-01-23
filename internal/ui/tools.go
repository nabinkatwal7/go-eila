package ui

import (
	"fmt"
	"math"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/nabinkatwal7/go-eila/internal/repository"
)

func NewToolsView(repo *repository.Repository) fyne.CanvasObject {
	header := widget.NewLabelWithStyle("Financial Tools", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	// Tab 1: Debt Payoff Simulator
	debtContent := createDebtCalculator()

	// Tab 2: Tax Estimator
	taxContent := createTaxEstimator()

	tabs := container.NewAppTabs(
		container.NewTabItem("Debt Payoff", debtContent),
		container.NewTabItem("Tax Estimator", taxContent),
	)

	return container.NewBorder(header, nil, nil, nil, tabs)
}

func createDebtCalculator() fyne.CanvasObject {
	principal := widget.NewEntry()
	principal.PlaceHolder = "Total Debt (e.g. 5000)"
	rate := widget.NewEntry()
	rate.PlaceHolder = "Annual Interest Rate % (e.g. 15)"
	payment := widget.NewEntry()
	payment.PlaceHolder = "Monthly Payment (e.g. 200)"

	result := widget.NewLabel("Enter details to calculate payoff time.")

	btn := widget.NewButton("Calculate", func() {
		p, _ := strconv.ParseFloat(principal.Text, 64)
		r, _ := strconv.ParseFloat(rate.Text, 64)
		pay, _ := strconv.ParseFloat(payment.Text, 64)

		if p <= 0 || pay <= 0 {
			result.SetText("Invalid values.")
			return
		}

		// N = -log(1 - (r/n * P) / A) / log(1 + r/n)
		// Simple iterative:
		balance := p
		months := 0
		monthlyRate := r / 100 / 12
		totalInterest := 0.0

		// Fail safe for infinite loop
		if balance * monthlyRate >= pay {
			result.SetText("Payment is too low to cover interest! You will be in debt forever.")
			return
		}

		for balance > 0 && months < 1000 {
			interest := balance * monthlyRate
			totalInterest += interest
			balance = balance + interest - pay
			months++
		}

		result.SetText(fmt.Sprintf("Debt Free in: %d months (%.1f years)\nTotal Interest Paid: $%.2f", months, float64(months)/12.0, totalInterest))
	})

	return container.NewVBox(
		widget.NewLabel("Simulate Debt Payoff Strategy"),
		principal,
		rate,
		payment,
		btn,
		widget.NewSeparator(),
		result,
	)
}

func createTaxEstimator() fyne.CanvasObject {
	incomeEntry := widget.NewEntry()
	incomeEntry.PlaceHolder = "Annual Taxable Income"

	result := widget.NewLabel("Enter income to estimate tax (Simplified US brackets).")

	btn := widget.NewButton("Estimate Tax", func() {
		inc, _ := strconv.ParseFloat(incomeEntry.Text, 64)

		// Simplified 2024 Single Filer Brackets (Approx)
		// 0 - 11,600: 10%
		// 11,601 - 47,150: 12%
		// 47,151 - 100,525: 22%
		// 100,526 - 191,950: 24%
		// 191,951+: 32% (capped for simplicity)

		var tax float64
		rem := inc

		brackets := []struct{ limit, Rate float64 }{
			{11600, 0.10},
			{35550, 0.12}, // 47150 - 11600
			{53375, 0.22}, // 100525 - 47150
			{91425, 0.24}, // 191950 - 100525
		}

		for _, b := range brackets {
			if rem <= 0 { break }
			taxable := math.Min(rem, b.limit)
			tax += taxable * b.Rate
			rem -= taxable
		}
		if rem > 0 {
			tax += rem * 0.32
		}

		effective := (tax / inc) * 100
		result.SetText(fmt.Sprintf("Estimated Tax: $%.2f\nEffective Rate: %.1f%%", tax, effective))
	})

	return container.NewVBox(
		widget.NewLabel("Simple Tax Estimator (Single Filer)"),
		incomeEntry,
		btn,
		widget.NewSeparator(),
		result,
	)
}
