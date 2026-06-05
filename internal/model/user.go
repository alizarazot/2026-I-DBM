package model

import "time"

type User struct {
	Email string

	Role UserRole
	Info UserInfo
}

type UserRole uint8

//go:generate go tool stringer -trimprefix UserRole -type UserRole

const (
	UserRoleInvalid UserRole = iota
	UserRoleManager
	UserRoleTeacher
	UserRoleStudent
)

type UserInfo struct {
	FirstName, MiddleName, FirstSurname, SecondSurname string

	Birthdate time.Time
	Genre     UserGenre
}

type UserGenre uint8

//go:generate go tool stringer -trimprefix UserGenre -type UserGenre

const (
	UserGenreInvalid UserGenre = iota
	UserGenreMale
	UserGenreFemale
	UserGenreOther
)
