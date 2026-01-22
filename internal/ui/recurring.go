package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/nabinkatwal7/go-eila/internal/repository"
)

func NewRecurringView(repo *repository.Repository) fyne.CanvasObject {
	header := widget.NewLabelWithStyle("Smart Recurring Detection", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	subs, err := repo.DetectRecurringPatterns()
	content := container.NewVBox()

	if err != nil {
		content.Add(widget.NewLabel("Error: " + err.Error()))
	} else if len(subs) == 0 {
		content.Add(widget.NewLabel("No recurring patterns detected yet. Add more history!"))
	} else {
		for _, s := range subs {
			card := widget.NewCard(s.Name, fmt.Sprintf("Est. $%.2f / %s", s.Amount, s.Frequency),
				widget.NewLabel("Next Due: " + s.NextDueDate),
			)
			content.Add(card)
		}
	}

	return container.NewBorder(header, nil, nil, nil, container.NewVScroll(content))
}
