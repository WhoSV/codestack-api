package model

// Favorite Type
type Favorite struct {
	ID       uint `json:"id,omitempty" gorm:"primary_key"`
	UserID   int  `json:"user_id,omitempty"`
	CourseID int  `json:"course_id,omitempty"`
}
