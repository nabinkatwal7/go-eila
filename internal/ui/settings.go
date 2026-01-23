package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/nabinkatwal7/go-eila/internal/repository"
)

func NewSettingsView(repo *repository.Repository, w fyne.Window) fyne.CanvasObject {
	header := widget.NewLabelWithStyle("Settings & Data", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	exportBtn := widget.NewButton("Export Data (JSON)", func() {
		dlg := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if writer == nil { return } // Cancelled

			err = repo.ExportDataToJSON(writer.URI().Path())
			if err != nil {
				dialog.ShowError(err, w)
			} else {
				dialog.ShowInformation("Success", "Data exported successfully.", w)
			}
		}, w)
		dlg.SetFileName("mytrack_backup.json")
		dlg.SetFilter(storage.NewExtensionFileFilter([]string{".json"}))
		dlg.Show()
	})

	snapshotBtn := widget.NewButton("Create Time-Travel Snapshot", func() {
		// Mock snapshot creation
		dialog.ShowInformation("Snapshot Created", "A snapshot of your database has been saved to 'snapshots/req_id.db'. (Mock)", w)
	})

	rolloverBtn := widget.NewButton("Run Budget Rollover", func() {
		// Mock rollover logic
		dialog.ShowInformation("Rollover Complete", "Unused budget from last month has been carried forward to current month. (Mock)", w)
	})

	return container.NewVBox(
		header,
		widget.NewSeparator(),
		widget.NewLabel("Data Integrity"),
		exportBtn,
		snapshotBtn,
		widget.NewSeparator(),
		widget.NewLabel("Advanced Features"),
		rolloverBtn,
		widget.NewLabel("More settings coming soon..."),
	)
}
