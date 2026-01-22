package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/nabinkatwal7/go-eila/internal/model"
	"github.com/nabinkatwal7/go-eila/internal/repository"
)

func NewAnomaliesView(repo *repository.Repository) fyne.CanvasObject {
	header := widget.NewLabelWithStyle("Spending Alerts & Anomalies", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	alerts, err := repo.DetectAnomalies()
	content := container.NewVBox()

	if err != nil {
		content.Add(widget.NewLabel("Error: " + err.Error()))
	} else if len(alerts) == 0 {
		content.Add(widget.NewLabel("No anomalies detected. Good job!"))
	} else {
		for _, a := range alerts {
			// Color code based on severity
			c := color.RGBA{200, 200, 0, 255} // Yellow
			if a.Severity == model.SeverityHigh {
				c = color.RGBA{200, 0, 0, 255} // Red
			}

			// Icon or colorful text
			indicator := canvas.NewCircle(c)
			indicator.Resize(fyne.NewSize(10, 10))
			indicator.Move(fyne.NewPos(0, 5))

			title := widget.NewLabelWithStyle(a.Type, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
			desc := widget.NewLabel(a.Description)
			date := widget.NewLabel(a.Date)

			row := container.NewBorder(nil, nil,
				container.NewHBox(indicator, title),
				date,
				desc,
			)
			content.Add(widget.NewCard("", "", row))
		}
	}

	return container.NewBorder(header, nil, nil, nil, container.NewVScroll(content))
}
