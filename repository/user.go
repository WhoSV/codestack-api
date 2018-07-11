package repository

import (
	"fmt"

	"github.com/WhoSV/codestack-api/database"
	"github.com/WhoSV/codestack-api/model"
)

// CreateNewPerson ...
func CreateNewPerson(email string, password string, role string, fullName string) (model.Person, error) {
	var person model.Person

	if email == "" || password == "" || role == "" || fullName == "" {
		return person, &UserError{What: "EmailOrPasswordOrRoleOrFullName", Type: "Empty", Arg: ""}
	}

	person = model.Person{
		Email:    email,
		FullName: fullName,
		Role:     role,
		Password: password,
	}

	var db = database.DB()

	q := db.Create(&person)
	if q.Error != nil {
		return model.Person{}, UserError{
			What: "User",
			Type: "Can-Not-Create",
			Arg:  email,
		}
	}

	return person, nil
}

// UserError ...
type UserError struct {
	What string
	Type string
	Arg  string
}

func (e UserError) Error() string {
	return fmt.Sprintf("%s: <%s> %s", e.Type, e.What, e.Arg)
}
