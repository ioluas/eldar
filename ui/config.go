package ui

import (
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func (m *Manager) MakeConfigPage() *widget.Form {
	form := widget.NewForm()
	endpointInput := widget.NewEntry()
	endpointInput.SetPlaceHolder("Enter your endpoint for Supabase project")
	anonKeyInput := widget.NewEntry()
	anonKeyInput.SetPlaceHolder("Enter your anonymous key for Supabase project")
	form.AppendItem(widget.NewFormItem("Endpoint", endpointInput))
	form.AppendItem(widget.NewFormItem("Anonymous Key", anonKeyInput))
	form.SubmitText = "Save"
	form.OnSubmit = func() {
		if err := m.db.SaveConfig(endpointInput.Text, anonKeyInput.Text); err != nil {
			dialog.ShowError(err, *m.window)
			return
		}
		m.SetAppPage(Unknown)
		m.UpdateWindowContent()
	}
	return form
}
