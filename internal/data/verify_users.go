package data

import (
	"time"
)

type VerifyUsers struct {
	UserID               string    `db:"user_id"`
	UserIDHash           string    `db:"user_id_hash"`
	AgeLowerBound        int       `db:"age_lower_bound"`
	Nationality          string    `db:"nationality"`
	CreatedAt            time.Time `db:"created_at"`
	Uniqueness           bool      `db:"uniqueness"`
	EventID              string    `db:"event_id"`
	Status               string    `db:"status"`
	Proof                []byte    `db:"proof"`
	Sex                  string    `db:"sex"`
	SexEnable            bool      `db:"sex_enable"`
	NationalityEnable    bool      `db:"nationality_enable"`
	AnonymousID          string    `db:"anonymous_id"`
	Nullifier            string    `db:"nullifier"`
	ExpirationLowerBound string    `db:"expiration_lower_bound"`
}

type VerifyUsersQ interface {
	New() VerifyUsersQ

	Get() (*VerifyUsers, error)
	Select() ([]VerifyUsers, error)
	Update(*VerifyUsers) error
	Upsert(*VerifyUsers) (VerifyUsers, error)
	Delete() error

	DeleteByID(*VerifyUsers) error
	WhereID(userId string) VerifyUsersQ
	WhereHashID(userId string) VerifyUsersQ
	WhereCreatedAtLt(createdAt time.Time) VerifyUsersQ
	FilterByInternalAID(aid string) VerifyUsersQ
	FilterByNullifier(nullifier string) VerifyUsersQ
}
