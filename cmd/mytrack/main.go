package main

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"github.com/nabinkatwal7/go-eila/internal/repository"
	"github.com/nabinkatwal7/go-eila/internal/ui"
)

func main() {
	// 1. Init DB
	db, err := repository.NewDB("mytrack.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 2. Init Repository
	repo := repository.NewRepository(db)

	// 3. Init UI
	myFyneApp := app.New()
	myWindow := myFyneApp.NewWindow("MyTrack")
	myWindow.Resize(fyne.NewSize(800, 600))

	myApp := ui.NewApp(myFyneApp, myWindow, repo)
	myApp.Init()
	myApp.ContentContainer = container.NewStack(ui.NewDashboard(repo))

	// 4. Run
	myApp.Run()
}
