package models

// LoginCredentials defines the structure for user sign-in
type LoginCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
