package models

// Request
type LoggingUser struct {
	UserName string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
}

type RegUser struct {
	FirstName string `json:"first_name" db:"first_name"`
	LastName  string `json:"last_name" db:"last_name"`
	UserName  string `json:"username" db:"username"`
	Password  string `json:"password" db:"password"`
}

// Response
type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
