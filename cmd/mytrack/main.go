package main

import (
	"log"

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
	myApp := ui.NewApp(repo)

	// 4. Run
	myApp.Run()
}
