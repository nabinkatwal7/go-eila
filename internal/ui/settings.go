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
			if writer == nil {
				return // Cancelled
			}

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

	importBtn := widget.NewButton("Import Data (JSON)", func() {
		dlg := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if reader == nil {
				return // Cancelled
			}
			defer reader.Close()

			// Confirm before importing
			dialog.ShowConfirm("Import Data",
				"This will add data from the backup file. Continue?",
				func(confirmed bool) {
					if confirmed {
						err = repo.ImportDataFromJSON(reader.URI().Path())
						if err != nil {
							dialog.ShowError(err, w)
						} else {
							dialog.ShowInformation("Success", "Data imported successfully.", w)
						}
					}
				}, w)
		}, w)
		dlg.SetFilter(storage.NewExtensionFileFilter([]string{".json"}))
		dlg.Show()
	})

	csvImportBtn := widget.NewButton("Import Transactions (CSV)", func() {
		dlg := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if reader == nil {
				return // Cancelled
			}
			defer reader.Close()

			// Show dialog for account/category selection
			accountEntry := widget.NewEntry()
			accountEntry.SetPlaceHolder("Default account name (optional)")
			categoryEntry := widget.NewEntry()
			categoryEntry.SetPlaceHolder("Default category name (optional)")

			form := widget.NewForm(
				widget.NewFormItem("Default Account", accountEntry),
				widget.NewFormItem("Default Category", categoryEntry),
			)

			importDlg := dialog.NewCustomConfirm("CSV Import", "Import", "Cancel", form,
				func(confirmed bool) {
					if confirmed {
						accountName := accountEntry.Text
						categoryName := categoryEntry.Text
						err = repo.ImportTransactionsFromCSV(reader.URI().Path(), accountName, categoryName)
						if err != nil {
							dialog.ShowError(err, w)
						} else {
							dialog.ShowInformation("Success", "Transactions imported successfully.", w)
						}
					}
				}, w)
			importDlg.Resize(fyne.NewSize(400, 200))
			importDlg.Show()
		}, w)
		dlg.SetFilter(storage.NewExtensionFileFilter([]string{".csv"}))
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
		widget.NewLabel("Data Export/Import"),
		exportBtn,
		importBtn,
		csvImportBtn,
		widget.NewSeparator(),
		widget.NewLabel("Advanced Features"),
		snapshotBtn,
		rolloverBtn,
		widget.NewLabel("More settings coming soon..."),
	)
}
