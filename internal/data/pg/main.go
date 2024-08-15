package pg

import (
	"github.com/rarimo/verificator-svc/internal/data"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const (
	verifyUsersTableName = "verify_users"
)

const (
	userIdColumnName     = "user_id"
	userIdHashColumnName = "user_id_hash"
	createdAtColumnName  = "created_at"
)

type masterQ struct {
	db *pgdb.DB
}

func NewMasterQ(db *pgdb.DB) data.MasterQ {
	return &masterQ{
		db: db,
	}
}

func (q *masterQ) New() data.MasterQ {
	return NewMasterQ(q.db.Clone())
}

func (q *masterQ) VerifyUsersQ() data.VerifyUsersQ {
	return NewVerifyUsersQ(q.db)
}

func (q *masterQ) Transaction(fn func() error) error {
	return q.db.Transaction(fn)
}
