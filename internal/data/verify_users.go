package data

import (
	"time"
)

type VerifyUsers struct {
	UserID        string    `db:"user_id"`
	UserIdHash    string    `db:"user_id_hash"`
	AgeLowerBound int       `db:"age_lower_bound"`
	Nationality   string    `db:"nationality"`
	CreatedAt     time.Time `db:"created_at"`
	Uniqueness    bool      `db:"uniqueness"`
	Status        string    `db:"status"`
}

type VerifyUsersQ interface {
	New() VerifyUsersQ

	Get() (*VerifyUsers, error)
	Select() ([]VerifyUsers, error)
	Update(*VerifyUsers) error
	Insert(*VerifyUsers) error
	Delete() error

	WhereID(userId string) VerifyUsersQ
	WhereHashID(userId string) VerifyUsersQ
	WhereCreatedAtLt(createdAt time.Time) VerifyUsersQ
}
