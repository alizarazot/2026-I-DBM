package model

import (
	"encoding/json"
	"time"
)

type User struct {
	Email string `json:"email"`

	Role UserRole `json:"role"`
	Info UserInfo `json:"info"`
}

type UserRole uint8

//go:generate go tool stringer -trimprefix UserRole -type UserRole

const (
	UserRoleInvalid UserRole = iota
	UserRoleManager
	UserRoleTeacher
	UserRoleStudent
)

func NewUserRole(canonical string) UserRole {
	switch canonical {
	case "manager":
		return UserRoleManager
	case "teacher":
		return UserRoleTeacher
	case "student":
		return UserRoleStudent
	default:
		return UserRoleInvalid
	}
}

func (ur UserRole) CanonicalString() string {
	switch ur {
	case UserRoleManager:
		return "manager"
	case UserRoleTeacher:
		return "teacher"
	case UserRoleStudent:
		return "student"
	default:
		return "invalid"
	}
}

func (ur UserRole) MarshalJSON() ([]byte, error) {
	return json.Marshal(ur.CanonicalString())
}

func (ur *UserRole) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	*ur = NewUserRole(s)

	return nil
}

type UserInfo struct {
	FirstName     string `json:"firstName"`
	MiddleName    string `json:"middleName"`
	FirstSurname  string `json:"firstSurname"`
	SecondSurname string `json:"secondSurname"`

	Birthdate time.Time `json:"birthdate"`
	Genre     UserGenre `json:"genre"`
}

type UserGenre uint8

//go:generate go tool stringer -type UserGenre

const (
	UserGenreInvalid UserGenre = iota
	UserGenreMale
	UserGenreFemale
	UserGenreOther
)

func NewUserGenre(canonical string) UserGenre {
	switch canonical {
	case "male":
		return UserGenreMale
	case "female":
		return UserGenreFemale
	case "other":
		return UserGenreOther
	default:
		return UserGenreInvalid
	}
}

func (ug UserGenre) CanonicalString() string {
	switch ug {
	case UserGenreMale:
		return "male"
	case UserGenreFemale:
		return "female"
	case UserGenreOther:
		return "other"
	default:
		return "invalid"
	}
}

func (ug UserGenre) MarshalJSON() ([]byte, error) {
	return json.Marshal(ug.CanonicalString())
}

func (ug *UserGenre) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	*ug = NewUserGenre(s)

	return nil
}
