package data

type MasterQ interface {
	New() MasterQ

	VerifyUsersQ() VerifyUsersQ

	Transaction(func() error) error
}
