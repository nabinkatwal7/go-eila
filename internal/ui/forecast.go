package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/nabinkatwal7/go-eila/internal/repository"
)

func NewForecastView(repo *repository.Repository) fyne.CanvasObject {
	header := widget.NewLabelWithStyle("Net Worth Projection", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	points, avgSavings, err := repo.GetNetWorthProjection(12) // 1 Year

	content := container.NewVBox()

	if err != nil {
		content.Add(widget.NewLabel("Error: " + err.Error()))
	} else {
		// Summary
		summary := fmt.Sprintf("Based on your recent average savings of $%.2f / month...", avgSavings)
		content.Add(widget.NewLabel(summary))

		// Simple Line Chart (Simulated with bars for now as Canvas Line is tricky with many points)
		// Or just a list of future milestones

		// Let's do a List of Milestones
		for _, p := range points {
			row := container.NewHBox(
				widget.NewLabel(p.Month+":"),
				widget.NewLabelWithStyle(fmt.Sprintf("$%.2f", p.Value), fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			)
			content.Add(row)
		}

		// Chart Visualization (Simulated Line)
		// We can reuse our specific chart logic or skip for now.
		// A list is fine for "dense" info.
	}

	return container.NewVScroll(container.NewVBox(
		header,
		content,
	))
}
