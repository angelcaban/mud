package registration

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	st "github.com/Masterminds/structable"
	"github.com/angelcaban/mud/model"
	"github.com/gofrs/uuid"
)

const (
	REGISTRATION_TABLE = "registrations"
	CONNECTION_STRING  = "/muddb"
)

type RegistrationRepository interface {
	// Save a registration into the database
	Store(registration *model.Registration) (*model.Registration, error)

	// Find a registration from the database given an ID
	Find(id uuid.UUID) *model.Registration

	// Delete a registration from the database given an ID
	Delete(id uuid.UUID) error

	// Get a list of all registrations
	FindAll() []*model.Registration
}

type repository struct {
	Db         sq.DBProxyBeginner
	DriverName string
}

func NewRegistrationRepository(db *sql.DB, driverName string) (RegistrationRepository, error) {
	conn, err := sql.Open(driverName, CONNECTION_STRING)
	if err != nil {
		return nil, err
	}
	return &repository{
		Db:         sq.NewStmtCacheProxy(conn),
		DriverName: driverName,
	}, nil
}

func (repo *repository) Store(registration *model.Registration) (*model.Registration, error) {
	recorder := st.New(repo.Db, repo.DriverName).Bind(REGISTRATION_TABLE, registration)
	err := recorder.Insert()
	if err != nil {
		return nil, err
	}

	recorder.Load()
	return registration, nil
}

func (repo *repository) Find(id uuid.UUID) *model.Registration {
	type_ := &model.Registration{}
	db_conn := st.New(repo.Db, repo.DriverName).Bind(REGISTRATION_TABLE, type_)
	rec, err := st.ListWhere(db_conn,
		func(object st.Describer, sql sq.SelectBuilder) (sq.SelectBuilder, error) {
			return sql.Limit(1).Where(sq.Eq{"id": id}), nil
		})
	if err != nil {
		return nil
	}

	return rec[0].Interface().(*model.Registration)
}

func (repo *repository) Delete(id uuid.UUID) error {
	reg := &model.Registration{Id: id}
	rec := st.New(repo.Db, repo.DriverName).Bind(REGISTRATION_TABLE, reg)
	return rec.Delete()
}

func (repo *repository) FindAll() []*model.Registration {
	reg := &model.Registration{}
	rec := st.New(repo.Db, repo.DriverName).Bind(REGISTRATION_TABLE, reg)

	allRegs := make([]*model.Registration, 1)

	var offset uint64 = 0
	var maxPageSize uint64 = 1000
	for {
		items, err := st.List(rec, maxPageSize, offset)
		if err == nil || len(items) > 0 {
			break
		}

		slice_ := make([]*model.Registration, len(items))
		for i, item := range items {
			slice_[i] = item.Interface().(*model.Registration)
		}
		append(allRegs, slice_)
		offset += uint64(len(items))
	}

	return allRegs
}
