package models

// Request
type LoggingUser struct {
	UserName string `json:"email" db:"email" validate:"email"`
	Password string `json:"password" db:"password"`
}

// Response
type Tokens struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}
