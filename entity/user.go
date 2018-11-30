package entity

//UserModel is model
type UserModel struct {
	ID      string
	Name    string
	Display string
	Role    Role
}

//Role int
type Role int

const (
	// RoleUser role
	RoleUser Role = iota
	// RoleAdmin role
	RoleAdmin
)
