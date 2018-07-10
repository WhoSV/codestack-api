package model

// Course status
const (
	StatusUndefined = "UNDEFINED"

	StatusActive  = "ACTIVE"
	StatusBlocked = "BLOCKED"
)

// Course Type
type Course struct {
	ID          uint   `json:"id,omitempty" gorm:"primary_key"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Teacher     string `json:"teacher,omitempty"`
	TeacherID   string `json:"teacher_id,omitempty"`
	Status      string `json:"status,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
	FileName    string `json:"file_name,omitempty"`
}
