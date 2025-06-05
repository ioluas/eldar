package main

import (
	"log"

	"eldar/credentials"
	"eldar/ui"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const (
	Register ui.AppPage = iota
	Login
	Group
	Boards
	Users
	Unknown
)

var appPage ui.AppPage
var w fyne.Window

// updateWindowContent updates the content of the window based on the current credentials
func updateWindowContent() {
	if appPage == Login {
		noCredsLabel := widget.NewLabel("Login")
		noCredsLabel.Alignment = fyne.TextAlignCenter
		w.SetContent(container.NewVBox(noCredsLabel, ui.MakeLoginForm(&appPage, updateWindowContent)))
		return
	}

	if appPage == Register {
		noCredsLabel := widget.NewLabel("Register")
		noCredsLabel.Alignment = fyne.TextAlignCenter
		w.SetContent(container.NewVBox(noCredsLabel, ui.MakeRegisterForm()))
		return
	}

	creds, err := credentials.GetCredentials()
	if err != nil {
		log.Printf("Error getting credentials: %v", err)
		creds = &credentials.Credentials{}
	}

	// No creds exist, we need to ask user to login
	if creds.Username == "" && creds.AccessToken == "" && creds.RefreshToken == "" {
		appPage = Login
		updateWindowContent()
		return
	}

	// Credentials are not empty, display todo for now
	w.SetContent(container.NewVBox(widget.NewLabel("TODO")))
}

func main() {
	appPage = Unknown
	a := app.NewWithID("dev.ioluas.eldar")
	w = a.NewWindow("Eldar")
	updateWindowContent()
	w.ShowAndRun()
}
