package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/nabinkatwal7/go-eila/internal/model"
)

// Simple Bar Chart
func NewBarChart(stats []model.MonthlyStat) fyne.CanvasObject {
	if len(stats) == 0 {
		return widget.NewLabel("No data for chart")
	}

	// Constants
	barWidth := float32(30)
	spacing := float32(10)
	maxHeight := float32(150)

	// Find max value for scaling
	var maxVal float64
	for _, s := range stats {
		if s.Income > maxVal { maxVal = s.Income }
		if s.Expense > maxVal { maxVal = s.Expense }
	}
	if maxVal == 0 { maxVal = 100 } // Avoid div by zero

	chartContainer := container.NewWithoutLayout()

	for i, s := range stats {
		xBase := float32(i) * (barWidth*2 + spacing)

		// Income Bar (Green)
		incHeight := float32(s.Income / maxVal) * maxHeight
		incBar := canvas.NewRectangle(color.RGBA{0, 200, 0, 200})
		incBar.Resize(fyne.NewSize(barWidth, incHeight))
		// Position: y starts from bottom.
		// In Fyne WithoutLayout, (0,0) is top-left usually?
		// Actually relative positioning is manual.
		// Let's assume bottom is at Y=200.
		baseY := float32(200)
		incBar.Move(fyne.NewPos(xBase, baseY-incHeight))

		// Expense Bar (Red)
		expHeight := float32(s.Expense / maxVal) * maxHeight
		expBar := canvas.NewRectangle(color.RGBA{200, 0, 0, 200})
		expBar.Resize(fyne.NewSize(barWidth, expHeight))
		expBar.Move(fyne.NewPos(xBase+barWidth, baseY-expHeight))

		// Label
		lbl := canvas.NewText(s.Month, color.Black)
		lbl.TextSize = 10
		lbl.Move(fyne.NewPos(xBase + barWidth/2, baseY+5))

		chartContainer.Add(incBar)
		chartContainer.Add(expBar)
		chartContainer.Add(lbl)
	}

	// Wrap in a sized container
	chartContainer.Resize(fyne.NewSize(float32(len(stats))*(barWidth*2+spacing), 220))

	// Scrollable if too many months?
	return container.NewPadded(chartContainer)
}
