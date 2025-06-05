// Package ui provides user interface components and functionality for the Eldar application.
// It includes forms for authentication, navigation, and other UI elements.
package ui

import (
	"errors"
	"log"
	"regexp"

	"fyne.io/fyne/v2/widget"
)

// Regular expressions used for password validation
var (
	// upperRe matches any uppercase letter
	upperRe = regexp.MustCompile(`[A-Z]`)
	// lowerRe matches any lowercase letter
	lowerRe = regexp.MustCompile(`[a-z]`)
	// digitRe matches any digit
	digitRe = regexp.MustCompile(`[0-9]`)
	// specialRe matches any special character
	specialRe = regexp.MustCompile(`[!"£$%^&*()\-_=+\][{}'@#~/?.>,<|]`)
)

// MakeLoginForm creates and returns a login form widget.
// It includes fields for email and password, along with a button to navigate to the registration page.
//
// Parameters:
//   - ap: A pointer to the current AppPage, which will be updated when navigating to the registration page
//   - updateWindow: A function to call when the app page changes to update the window content
//
// Returns:
//   - A configured widget.Form ready to be displayed
func MakeLoginForm(ap *AppPage, updateWindow func()) *widget.Form {
	form := widget.NewForm()
	emailInput := widget.NewEntry()
	emailInput.SetPlaceHolder("Enter your email address")
	form.AppendItem(widget.NewFormItem("Email", emailInput))
	passwordInput := widget.NewPasswordEntry()
	passwordInput.SetPlaceHolder("Enter your password")
	form.AppendItem(widget.NewFormItem("Password", passwordInput))
	registerButton := widget.NewButton("Register", func() {
		*ap = Register
		updateWindow()
	})
	form.AppendItem(widget.NewFormItem("Don't have an account yet?", registerButton))
	form.SubmitText = "Login"
	form.OnSubmit = func() {
		log.Printf("Login user with email: %s, password: %s", emailInput.Text, passwordInput.Text)
	}
	return form
}

// MakeRegisterForm creates and returns a registration form widget.
// It includes fields for email, password, and password confirmation with validation.
// The password must contain at least one uppercase letter, one lowercase letter,
// one digit, one special character, and be between 8 and 255 characters long.
// The form is disabled until all validation requirements are met.
//
// Returns:
//   - A configured widget.Form ready to be displayed
func MakeRegisterForm() *widget.Form {
	form := widget.NewForm()
	emailInput := widget.NewEntry()
	emailInput.SetPlaceHolder("Enter your email address")
	form.AppendItem(widget.NewFormItem("Email", emailInput))
	passwordInput := widget.NewPasswordEntry()
	passwordInput.SetPlaceHolder("Enter your password")
	form.AppendItem(widget.NewFormItem("Password", passwordInput))
	passwordConfirmInput := widget.NewPasswordEntry()
	passwordInput.Validator = func(s string) error {
		log.Printf("validation of pass input: %s", s)
		matched := upperRe.MatchString(s) && lowerRe.MatchString(s) && digitRe.MatchString(s) && specialRe.MatchString(s)
		length := len(s)
		if !matched || length < 8 || length > 255 {
			log.Printf("password doesn't match requirements")
			return errors.New("invalid password")
		}
		return nil
	}
	passwordConfirmInput.SetPlaceHolder("Confirm your password")
	passwordConfirmInput.Validator = func(s string) error {
		currentPass := passwordInput.Text
		log.Printf("pass: %s, confirm: %s", currentPass, s)
		if s != currentPass {
			log.Printf("password doesn't match")
			return errors.New("passwords do not match")
		}
		return nil
	}
	onChanged := func(owner string) func(s string) {
		return func(s string) {
			var mirror string
			switch owner {
			case "passwordInput":
				mirror = passwordConfirmInput.Text
			case "passwordConfirmInput":
				mirror = passwordInput.Text
			default:
				panic("Unknown owner")
			}
			log.Printf("current: %s. mirror: %s", s, mirror)
			if s != mirror {
				log.Printf("passwords don't match")
				passwordConfirmInput.SetValidationError(errors.New("passwords do not match"))
				form.Disable()
			} else {
				passwordConfirmInput.SetValidationError(nil)
				form.Enable()
			}
		}
	}
	passwordConfirmInput.OnChanged = func(s string) {
		log.Printf("confirm onChanged: %s", s)
		onChanged("passwordConfirmInput")(s)
	}
	passwordInput.OnChanged = func(s string) {
		log.Printf("pass onChanged: %s", s)
		if err := passwordConfirmInput.Validate(); err != nil {
			passwordConfirmInput.SetValidationError(err)
			form.Disable()
			return
		}
		if err := passwordInput.Validator(s); err != nil {
			passwordInput.SetValidationError(err)
			form.Disable()
		} else {
			form.Enable()
		}
	}

	form.AppendItem(widget.NewFormItem("Confirm Password", passwordConfirmInput))
	form.SubmitText = "Register"
	form.OnSubmit = func() {
		log.Printf("Registering user with email: %s, password: %s", emailInput.Text, passwordInput.Text)
	}
	return form
}
