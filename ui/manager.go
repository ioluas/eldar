package ui

import (
	"eldar/api"
	"eldar/storage"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// NewUIManager creates and initializes a new UI Manager instance.
// It handles the application's user interface state and navigation.
//
// Parameters:
//   - window: A pointer to the main application window
//   - db: A pointer to the application's database instance
//   - api: A pointer to the application's Supabase api client
//
// Returns:
//   - A new Manager instance initialized with the provided window and database
func NewUIManager(window *fyne.Window, db *storage.Database, api *api.DatabaseHTTPClient) *Manager {
	return &Manager{
		window:  window,
		db:      db,
		appPage: Unknown,
		api:     api,
	}
}

// SetAppPage updates the current application page state in the UI Manager.
// This method is used for navigation between different sections of the application.
//
// Parameters:
//   - page: The AppPage enum value representing the target page to display
func (m *Manager) SetAppPage(page AppPage) {
	m.appPage = page
}

// UpdateWindowContent updates the window's content based on the current application page state.
// It handles the navigation flow and content rendering for different pages:
// - Login: Displays the login form
// - Register: Shows the registration form
// - Config: Presents the configuration page
// - Boards, Group, Users: Shows placeholder content (TODO)
// - Unknown: Handles default navigation logic
//
// If the database configuration is incomplete (missing endpoint or anonymous key),
// it redirects to the Config page.
// If user credentials are missing, it redirects to the Login page.
// If both configurations are present, it redirects to the Boards page.
func (m *Manager) UpdateWindowContent() {
	switch m.appPage {
	case Login:
		loginLabel := widget.NewLabel("Login")
		loginLabel.Alignment = fyne.TextAlignCenter
		(*m.window).SetContent(container.NewVBox(loginLabel, m.MakeLoginForm()))
		return
	case Register:
		registerLabel := widget.NewLabel("Register")
		registerLabel.Alignment = fyne.TextAlignCenter
		(*m.window).SetContent(container.NewVBox(registerLabel, m.MakeRegisterForm()))
		return
	case Config:
		configLabel := widget.NewLabel("Config")
		configLabel.Alignment = fyne.TextAlignCenter
		(*m.window).SetContent(container.NewVBox(configLabel, m.MakeConfigPage()))
		return
	case Unknown:
		// TODO.
		break
	case Boards, Group, Users:
		(*m.window).SetContent(container.NewVBox(widget.NewLabel("TODO")))
		return
	default:
		panic("Unknown AppPage")
	}

	if m.db.Config.Endpoint == "" || m.db.Config.AnonKey == "" {
		m.SetAppPage(Config)
		m.UpdateWindowContent()
		return
	}
	if m.db.Credentials.Username == "" || m.db.Credentials.AccessToken == "" || m.db.Credentials.RefreshToken == "" {
		m.SetAppPage(Login)
		m.UpdateWindowContent()
		return
	}
	m.SetAppPage(Boards)
	m.UpdateWindowContent()
}
