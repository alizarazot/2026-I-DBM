package model

import (
	"encoding/json"
	"strings"
	"time"
)

type User struct {
	Email string `json:"email"`

	Role UserRole `json:"role"`
	Info UserInfo `json:"info"`
}

type UserRole uint8

//go:generate go tool stringer -trimprefix UserRole -type UserRole

func (ur UserRole) MarshalJSON() ([]byte, error) {
	return json.Marshal(strings.ToLower(ur.String()))
}

const (
	UserRoleInvalid UserRole = iota
	UserRoleManager
	UserRoleTeacher
	UserRoleStudent
)

type UserInfo struct {
	FirstName     string `json:"firstName"`
	MiddleName    string `json:"middleName"`
	FirstSurname  string `json:"firstSurname"`
	SecondSurname string `json:"secondSurname"`

	Birthdate time.Time `json:"birthdate"`
	Genre     UserGenre `json:"genre"`
}

type UserGenre uint8

//go:generate go tool stringer -trimprefix UserGenre -type UserGenre

func (ug UserGenre) MarshalJSON() ([]byte, error) {
	return json.Marshal(strings.ToLower(ug.String()))
}

const (
	UserGenreInvalid UserGenre = iota
	UserGenreMale
	UserGenreFemale
	UserGenreOther
)
