package pg

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/rarimo/verificator-svc/internal/data"
	"gitlab.com/distributed_lab/kit/pgdb"
)

type VerifyUsersQ struct {
	db  *pgdb.DB
	sel sq.SelectBuilder
	del sq.DeleteBuilder
}

func NewVerifyUsersQ(db *pgdb.DB) data.VerifyUsersQ {
	return &VerifyUsersQ{
		db:  db,
		sel: sq.Select("*").From(verifyUsersTableName),
		del: sq.Delete(verifyUsersTableName),
	}
}

func (q *VerifyUsersQ) New() data.VerifyUsersQ {
	return NewVerifyUsersQ(q.db.Clone())
}

func (q *VerifyUsersQ) Select() ([]data.VerifyUsers, error) {
	var result []data.VerifyUsers

	err := q.db.Select(&result, q.sel)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to select rows")
	}

	return result, nil
}

func (q *VerifyUsersQ) Get() (*data.VerifyUsers, error) {
	var result data.VerifyUsers

	err := q.db.Get(&result, q.sel)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to get row")
	}

	return &result, nil
}

func (q *VerifyUsersQ) Upsert(VerifyUsers *data.VerifyUsers) (data.VerifyUsers, error) {
	var response data.VerifyUsers
	proofJSON, err := json.Marshal(VerifyUsers.Proof)
	if err != nil {
		return response, fmt.Errorf("failed to marshal proof for user %s: %w", VerifyUsers.UserID, err)
	}

	newData := map[string]interface{}{
		"user_id_hash":                 VerifyUsers.UserIDHash,
		"age_lower_bound":              VerifyUsers.AgeLowerBound,
		"nationality":                  VerifyUsers.Nationality,
		"uniqueness":                   VerifyUsers.Uniqueness,
		"event_id":                     VerifyUsers.EventID,
		"status":                       VerifyUsers.Status,
		"proof":                        proofJSON,
		"sex":                          VerifyUsers.Sex,
		"sex_enable":                   VerifyUsers.SexEnable,
		"nationality_enable":           VerifyUsers.NationalityEnable,
		"anonymous_id":                 VerifyUsers.AnonymousID,
		"nullifier":                    VerifyUsers.Nullifier,
		"expiration_lower_bound":       VerifyUsers.ExpirationLowerBound,
		"identity_counter":             VerifyUsers.IdentityCounter,
		"identity_counter_lower_bound": VerifyUsers.IdentityCounterLowerBound,
		"identity_counter_upper_bound": VerifyUsers.IdentityCounterUpperBound,
		"selector":                     VerifyUsers.Selector,
		"birth_date_lower_bound":       VerifyUsers.BirthDateLowerBound,
		"birth_date_upper_bound":       VerifyUsers.BirthDateUpperBound,
		"event_data":                   VerifyUsers.EventData,
		"expiration_date_upper_bound":  VerifyUsers.ExpirationDateUpperBound,
		"timestamp_lower_bound":        VerifyUsers.TimestampLowerBound,
		"timestamp_upper_bound":        VerifyUsers.TimestampUpperBound,
	}

	updateStmt, args, _ := sq.Update(" ").SetMap(newData).ToSql()

	newData["user_id"] = VerifyUsers.UserID

	query := sq.Insert(verifyUsersTableName).SetMap(newData).
		Suffix("ON CONFLICT (user_id) DO "+updateStmt, args...).
		Suffix("RETURNING *")

	if err = q.db.Get(&response, query); err != nil {
		return response, errors.Wrap(err, "failed to upsert new row")
	}

	return response, nil
}

func (q *VerifyUsersQ) Update(VerifyUsers *data.VerifyUsers) error {
	err := q.db.Exec(
		sq.Update(verifyUsersTableName).
			SetMap(map[string]interface{}{
				"status":       VerifyUsers.Status,
				"proof":        VerifyUsers.Proof,
				"sex":          VerifyUsers.Sex,
				"nationality":  VerifyUsers.Nationality,
				"anonymous_id": VerifyUsers.AnonymousID,
				"nullifier":    VerifyUsers.Nullifier,
			}).
			Where(sq.Eq{userIdColumnName: VerifyUsers.UserID}),
	)
	if err != nil {
		return errors.Wrap(err, "failed to update rows")
	}

	return nil
}

func (q *VerifyUsersQ) Delete() error {
	err := q.db.Exec(q.del)
	if err != nil {
		return errors.Wrap(err, "failed to delete rows")
	}

	return nil
}

func (q *VerifyUsersQ) DeleteByID(VerifyUsers *data.VerifyUsers) error {
	err := q.db.Exec(
		sq.Delete(verifyUsersTableName).Where(sq.Eq{userIdColumnName: VerifyUsers.UserID}),
	)
	if err != nil {
		return errors.Wrap(err, "failed to delete user by userID")
	}
	return nil

}

func (q *VerifyUsersQ) WhereID(userId string) data.VerifyUsersQ {
	q.sel = q.sel.Where(sq.Eq{userIdColumnName: userId})
	q.del = q.del.Where(sq.Eq{userIdColumnName: userId})
	return q
}

func (q *VerifyUsersQ) WhereHashID(userHashId string) data.VerifyUsersQ {
	q.sel = q.sel.Where(sq.Eq{userIdHashColumnName: userHashId})
	q.del = q.del.Where(sq.Eq{userIdHashColumnName: userHashId})
	return q
}

func (q *VerifyUsersQ) WhereCreatedAtLt(createdAt time.Time) data.VerifyUsersQ {
	q.sel = q.sel.Where(sq.Lt{createdAtColumnName: &createdAt})
	q.del = q.del.Where(sq.Lt{createdAtColumnName: &createdAt})
	return q
}

func (q *VerifyUsersQ) FilterByInternalAID(aid string) data.VerifyUsersQ {
	q.sel = q.sel.Where(sq.Eq{"anonymous_id": aid})
	return q
}

func (q *VerifyUsersQ) FilterByNullifier(nullifier string) data.VerifyUsersQ {
	q.sel = q.sel.Where(sq.Eq{"nullifier": nullifier})
	return q
}
