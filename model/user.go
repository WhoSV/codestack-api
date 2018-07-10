package model

// User roles
const (
	RoleUndefined = "UNDEFINED"

	RoleAdmin   = "ADMIN"
	RoleTeacher = "TEACHER"
	RoleStudent = "STUDENT"
)

// Person Type
type Person struct {
	ID       uint   `json:"id,omitempty" gorm:"primary_key"`
	FullName string `json:"full_name,omitempty"`
	Email    string `json:"email,omitempty" gorm:"unique,not null"`
	Role     string `json:"role,omitempty"`
	Password string `json:"password,omitempty"`
	Token    string `json:"token"`
}
