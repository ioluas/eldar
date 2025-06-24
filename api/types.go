package api

import "net/http"

type DatabaseHTTPClient struct {
	URL         string
	APIKey      string
	Client      *http.Client
	Headers     map[string]string
	AccessToken string
	UserID      string
	UserEmail   string
}

type SignUpRequest struct {
	Email    string                 `json:"email"`
	Password string                 `json:"password"`
	Options  map[string]interface{} `json:"options"`
}

type SignUpOption func(sr SignUpRequest) SignUpRequest

type SignUpResponse struct {
	Data  *signUpData `json:"data"`
	Error *AuthError  `json:"error"`
}

type signUpData struct {
	User    *User    `json:"user"`
	Session *Session `json:"session"`
}

type User struct {
	Id               string                 `json:"id"`
	Email            *string                `json:"email"`
	Phone            *string                `json:"phone"`
	EmailConfirmedAt *string                `json:"email_confirmed_at"`
	PhoneConfirmedAt *string                `json:"phone_confirmed_at"`
	CreatedAt        string                 `json:"created_at"`
	UpdatedAt        string                 `json:"updated_at"`
	UserMetadata     map[string]interface{} `json:"user_metadata"`
	AppMetadata      map[string]interface{} `json:"app_metadata"`
	Aud              string                 `json:"aud"`
	Role             *string                `json:"role"`
}

type Session struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	ExpiresAt    int    `json:"expires_at"`
	User         User   `json:"user"`
}

type AuthError struct {
	Message string  `json:"message"`
	Status  *int    `json:"status"`
	Code    *string `json:"code"`
}
