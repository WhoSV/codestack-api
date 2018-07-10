package model

// Survey Type
type Survey struct {
	ID       uint `json:"id,omitempty" gorm:"primary_key"`
	CourseID int  `json:"course_id,omitempty"`
	First    int  `json:"first,omitempty"`
	Second   int  `json:"second,omitempty"`
	Third    int  `json:"third,omitempty"`
	Fourth   int  `json:"fourth,omitempty"`
	Fifth    int  `json:"fifth,omitempty"`
}
