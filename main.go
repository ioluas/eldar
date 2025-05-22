package main

import (
	"fmt"
	"log"

	"eldar/credentials"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// updateWindowContent updates the content of the window based on the current credentials
func updateWindowContent(w fyne.Window) {
	// Get credentials from the bbolt database
	creds, err := credentials.GetCredentials()
	if err != nil {
		log.Printf("Error getting credentials: %v", err)
		// Continue with empty credentials
		creds = &credentials.Credentials{}
	}

	// Create a container for the main content
	var content *fyne.Container

	// Check if any credentials were found
	if creds.Username == "" && creds.AccessToken == "" && creds.RefreshToken == "" {
		// No credentials found
		noCredsLabel := widget.NewLabel("No credentials found")
		noCredsLabel.Alignment = fyne.TextAlignCenter

		// Add test credentials button
		addTestButton := widget.NewButton("Add Test Credentials", func() {
			err := credentials.AddTestCredentials()
			if err != nil {
				log.Printf("Error adding test credentials: %v", err)
				return
			}
			// Update the window content to show the new credentials
			updateWindowContent(w)
		})

		content = container.NewVBox(
			noCredsLabel,
			addTestButton,
		)
	} else {
		// Display the credentials
		usernameLabel := widget.NewLabel(fmt.Sprintf("Username: %s", creds.Username))
		accessTokenLabel := widget.NewLabel(fmt.Sprintf("Access Token: %s", creds.AccessToken))
		refreshTokenLabel := widget.NewLabel(fmt.Sprintf("Refresh Token: %s", creds.RefreshToken))

		// Clear credentials button
		clearButton := widget.NewButton("Clear Credentials", func() {
			err := credentials.ClearCredentials()
			if err != nil {
				log.Printf("Error clearing credentials: %v", err)
				return
			}
			// Update the window content to show the empty state
			updateWindowContent(w)
		})

		content = container.NewVBox(
			usernameLabel,
			accessTokenLabel,
			refreshTokenLabel,
			clearButton,
		)
	}

	// Set the content of the window
	w.SetContent(content)
}

func main() {
	// Create a new Fyne application for Android
	a := app.NewWithID("com.ioluas.eldar")

	// Create a window for the app
	w := a.NewWindow("Eldar Credentials")

	// Initialize the window content
	updateWindowContent(w)

	// Show and run the application
	w.ShowAndRun()
}
