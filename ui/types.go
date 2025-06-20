// Package ui provides user interface components and functionality for the Eldar application.
// It includes forms for authentication, navigation, and other UI elements.
package ui

import (
	"eldar/api"
	"eldar/storage"

	"fyne.io/fyne/v2"
)

// AppPage represents the different pages/screens in the application.
// It's used for navigation and determining which UI components to display.
type AppPage int

// Application page constants
const (
	Register AppPage = iota // Registration page
	Config                  // Config page, endpoint and anon-key config
	Login                   // Login page
	Group                   // Group management page
	Boards                  // Boards/dashboard page
	Users                   // User management page
	Unknown                 // Unknown/default page
)

// String returns the string representation of an AppPage.
// This is useful for debugging and logging purposes.s
func (ap AppPage) String() string {
	switch ap {
	case Register:
		return "Register"
	case Config:
		return "Config"
	case Login:
		return "Login"
	case Group:
		return "Group"
	case Boards:
		return "Boards"
	case Users:
		return "Users"
	case Unknown:
		return "Unknown"
	default:
		panic("Unknown AppPage")
	}
}

type Manager struct {
	appPage AppPage
	window  *fyne.Window
	db      *storage.Database
	api     *api.DatabaseHTTPClient
}
