package pg

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/fatih/structs"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/kit/pgdb"

	sq "github.com/Masterminds/squirrel"

	"github.com/rarimo/verificator-svc/internal/data"
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

func (q *VerifyUsersQ) Insert(VerifyUsers *data.VerifyUsers) error {
	stmt := sq.Insert(verifyUsersTableName).SetMap(map[string]interface{}{
		"user_id":      VerifyUsers.UserID,
		"user_id_hash": VerifyUsers.UserIdHash,
		"created_at":   VerifyUsers.CreatedAt,
		"status":       VerifyUsers.Status,
	})

	if err := q.db.Exec(stmt); err != nil {
		return fmt.Errorf("insert balance %+v: %w", VerifyUsers, err)
	}

	return nil
}

func (q *VerifyUsersQ) Update(VerifyUsers *data.VerifyUsers) error {
	err := q.db.Exec(
		sq.Update(verifyUsersTableName).
			SetMap(structs.Map(VerifyUsers)).
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
