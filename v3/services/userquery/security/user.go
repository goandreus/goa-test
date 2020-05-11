package security

// User represents an authenticated user
type User struct {
	ID        string   `json:"id"`
	Login     string   `json:"login"`
	Firstname string   `json:"fistname"`
	Lastname  string   `json:"lastname"`
	Email     string   `json:"email"`
	Language  string   `json:"language"`
	Roles     []string `json:"roles"`
}
