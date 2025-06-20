package main

import (
	"log"

	"eldar/storage"
	"eldar/ui"

	"fyne.io/fyne/v2/app"
)

var db *storage.Database

func main() {
	tmp, err := storage.InitDatabase()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	db = tmp
	if err = db.LoadConfig(); err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	if err = db.LoadCredentials(); err != nil {
		log.Fatalf("Error loading credentials: %v", err)
	}

	a := app.NewWithID("dev.ioluas.eldar")
	w := a.NewWindow("Eldar")
	uiManager := ui.NewUIManager(&w, db)
	uiManager.UpdateWindowContent()
	w.ShowAndRun()
}
