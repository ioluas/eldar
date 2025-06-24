package main

import (
	"log"

	"eldar/api"
	"eldar/storage"
	"eldar/ui"
	"fyne.io/fyne/v2"

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

	client := api.NewClient(db.Config.Endpoint, db.Config.AnonKey)

	a := app.NewWithID("dev.ioluas.eldar")
	w := a.NewWindow("Eldar")
	w.Resize(fyne.NewSize(1024, 768))
	uiManager := ui.NewUIManager(&w, db, client)
	uiManager.UpdateWindowContent()
	w.ShowAndRun()
}
