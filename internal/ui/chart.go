package ui

import (
	"fmt"
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

// NewLineChart creates a line chart for net worth progression
func NewLineChart(points []model.NetWorthPoint) fyne.CanvasObject {
	if len(points) == 0 {
		return widget.NewLabel("No data for chart")
	}

	maxHeight := float32(200)
	pointSpacing := float32(40)
	chartWidth := float32(len(points)) * pointSpacing

	// Find max and min values for scaling
	var maxVal, minVal float64
	for _, p := range points {
		if p.NetWorth > maxVal {
			maxVal = p.NetWorth
		}
		if p.NetWorth < minVal {
			minVal = p.NetWorth
		}
	}
	rangeVal := maxVal - minVal
	if rangeVal == 0 {
		rangeVal = 1
	}

	chartContainer := container.NewWithoutLayout()
	baseY := float32(200)

	// Draw line connecting points
	for i := 0; i < len(points)-1; i++ {
		x1 := float32(i) * pointSpacing
		y1 := baseY - float32((points[i].NetWorth-minVal)/rangeVal)*maxHeight
		x2 := float32(i+1) * pointSpacing
		y2 := baseY - float32((points[i+1].NetWorth-minVal)/rangeVal)*maxHeight

		line := canvas.NewLine(color.RGBA{0, 100, 200, 255})
		line.StrokeWidth = 2
		line.Position1 = fyne.NewPos(x1, y1)
		line.Position2 = fyne.NewPos(x2, y2)
		chartContainer.Add(line)
	}

	// Draw points and labels
	for i, p := range points {
		x := float32(i) * pointSpacing
		y := baseY - float32((p.NetWorth-minVal)/rangeVal)*maxHeight

		// Point
		point := canvas.NewCircle(color.RGBA{0, 100, 200, 255})
		point.Resize(fyne.NewSize(6, 6))
		point.Move(fyne.NewPos(x-3, y-3))
		chartContainer.Add(point)

		// Label
		lbl := canvas.NewText(p.Month, color.Black)
		lbl.TextSize = 9
		lbl.Move(fyne.NewPos(x-10, baseY+5))
		chartContainer.Add(lbl)
	}

	chartContainer.Resize(fyne.NewSize(chartWidth, 250))
	return container.NewPadded(chartContainer)
}

// NewPieChart creates a simple pie chart visualization for category breakdown
func NewPieChart(breakdown []model.CategoryBreakdown) fyne.CanvasObject {
	if len(breakdown) == 0 {
		return widget.NewLabel("No data for chart")
	}

	// Calculate total
	var total float64
	for _, b := range breakdown {
		total += b.Amount
	}

	if total == 0 {
		return widget.NewLabel("No spending data")
	}

	// Create a list view showing categories with visual bars
	list := widget.NewList(
		func() int {
			return len(breakdown)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewLabel("Category"),
				widget.NewLabel("Amount"),
				widget.NewProgressBar(),
			)
		},
		func(i int, o fyne.CanvasObject) {
			b := breakdown[i]
			box := o.(*fyne.Container)
			box.Objects[0].(*widget.Label).SetText(b.CategoryName)
			box.Objects[1].(*widget.Label).SetText(fmt.Sprintf("$%.2f", b.Amount))
			box.Objects[2].(*widget.ProgressBar).SetValue(b.Amount / total)
		},
	)

	return container.NewPadded(list)
}