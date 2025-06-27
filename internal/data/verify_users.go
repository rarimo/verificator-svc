package data

import (
	"database/sql"
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

	BirthDateLowerBound       sql.NullString `db:"birth_date_lower_bound"`
	BirthDateUpperBound       sql.NullString `db:"birth_date_upper_bound"`
	EventData                 sql.NullString `db:"event_data"`
	ExpirationDateUpperBound  sql.NullString `db:"expiration_date_upper_bound"`
	IdentityCounter           int32          `db:"identity_counter"`
	IdentityCounterLowerBound int32          `db:"identity_counter_lower_bound"`
	IdentityCounterUpperBound int32          `db:"identity_counter_upper_bound"`
	Selector                  int32          `db:"selector"`
	TimestampLowerBound       sql.NullTime   `db:"timestamp_lower_bound"`
	TimestampUpperBound       sql.NullTime   `db:"timestamp_upper_bound"`
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
	FilterByEventData(eventData string) VerifyUsersQ
}
