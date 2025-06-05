package ui

import (
	"strings"
	"testing"

	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
	"github.com/stretchr/testify/assert"
)

func TestMakeLoginForm(t *testing.T) {
	ap := Login
	updateWindowCalled := false
	updateWindow := func() {
		updateWindowCalled = true
	}

	form := MakeLoginForm(&ap, updateWindow)
	assert.NotNil(t, form)
	assert.Equal(t, 3, len(form.Items))
	assert.Equal(t, "Login", form.SubmitText)

	// Email
	emailEntry := form.Items[0].Widget.(*widget.Entry)
	assert.Equal(t, "Email", form.Items[0].Text)
	assert.Equal(t, "Enter your email address", emailEntry.PlaceHolder)

	// Password
	passwordEntry := form.Items[1].Widget.(*widget.Entry)
	assert.Equal(t, "Password", form.Items[1].Text)
	assert.Equal(t, "Enter your password", passwordEntry.PlaceHolder)
	assert.True(t, passwordEntry.Password)

	// Register
	registerButton := form.Items[2].Widget.(*widget.Button)
	assert.Equal(t, "Register", registerButton.Text)

	// Test register button click changes app page
	assert.Equal(t, Login, ap)
	registerButton.OnTapped()
	assert.Equal(t, Register, ap)
	assert.True(t, updateWindowCalled)
}

func TestMakeRegisterForm(t *testing.T) {
	form := MakeRegisterForm()
	assert.NotNil(t, form)
	assert.Equal(t, 3, len(form.Items))
	assert.Equal(t, "Register", form.SubmitText)

	// Test weak password validation
	passwordEntry := form.Items[1].Widget.(*widget.Entry)
	test.Type(passwordEntry, "weak")
	assert.Error(t, passwordEntry.Validator("weak"))
	passwordEntry.SetText("")
	// Test no upper case char validation
	test.Type(passwordEntry, "fsdfjsodijfowejf444=-")
	assert.Error(t, passwordEntry.Validator("fsdfjsodijfowejf444=-"))
	passwordEntry.SetText("")
	// Test no lower case char validation
	test.Type(passwordEntry, "SADFASDFWEF444=-")
	assert.Error(t, passwordEntry.Validator("SADFASDFWEF444=-"))
	passwordEntry.SetText("")

	// Test no digit char validation
	test.Type(passwordEntry, "SADFASDFWEFsdfasdf=-")
	assert.Error(t, passwordEntry.Validator("SADFASDFWEFsdfasdf=-"))
	passwordEntry.SetText("")

	// Test no special char validation
	test.Type(passwordEntry, "SADFASDFWEFsdfasdf33")
	assert.Error(t, passwordEntry.Validator("SADFASDFWEFsdfasdf33"))
	passwordEntry.SetText("")

	// Test too short password validation
	test.Type(passwordEntry, "Sd2-")
	assert.Error(t, passwordEntry.Validator("Sd2-"))
	passwordEntry.SetText("")

	// Test too long password validation
	test.Type(passwordEntry, strings.Repeat("Sd2-", 200))
	assert.Error(t, passwordEntry.Validator(strings.Repeat("Sd2-", 200)))
	passwordEntry.SetText("")

	// Test valid password
	test.Type(passwordEntry, "StrongP@ss123")
	assert.Nil(t, passwordEntry.Validator("StrongP@ss123"))

	// Test password confirmation
	confirmEntry := form.Items[2].Widget.(*widget.Entry)
	test.Type(confirmEntry, "different")
	assert.Error(t, confirmEntry.Validator("different"))
	confirmEntry.SetText("")
	test.Type(confirmEntry, "StrongP@ss123")
	assert.Nil(t, confirmEntry.Validator("StrongP@ss123"))
}
