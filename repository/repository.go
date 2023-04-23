package repository

import (
	"context"
	"database/sql"
	"github.com/getground/tech-tasks/backend/errors"
	"github.com/getground/tech-tasks/backend/types"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"net/http"
)

var (
	insertTableQuery           = "INSERT INTO tables (capacity, allowed_extras, created_at) VALUES (:capacity, :allowed_extras, :created_at)"
	insertGuestQuery           = "INSERT INTO guests (name, table_id, accompanying_guests, created_at) VALUES (:name, :table_id, :accompanying_guests, :created_at)"
	selectTableQuery           = "SELECT * FROM tables WHERE id=?"
	selectGuestsByTableIdQuery = "SELECT * FROM guests WHERE table_id=? AND deleted_at IS NULL"
	selectGuestsByNameQuery    = "SELECT * FROM guests WHERE name=? AND deleted_at IS NULL"
	updateGuestQuery           = "UPDATE guests SET accompanying_guests = ? WHERE name=?"
	removeGuestQuery           = "UPDATE guests SET deleted_at=now() WHERE name=?"
	queryGuests                = "SELECT * FROM guests WHERE deleted_at IS NULL"
	queryTables                = "SELECT * FROM tables"
)

type Repository interface {
	AddTable(capacity, allowedExtras int) (*types.Table, error)
	AddGuest(name string, tableId, accompanyingGuests int) (*types.Guest, error)
	GetTable(id int) (*types.Table, error)
	GetGuestByTableId(tableId int) (*types.Guest, error)
	GetGuestByName(name string) (*types.Guest, error)
	UpdateGuest(name string, accompanyingGuests int) error
	RemoveGuests(name string) error
	GetAllTables() ([]*types.Table, error)
	GetAllGuests() ([]*types.Guest, error)
}

type repository struct {
	db     *sqlx.DB
	logger *logrus.Logger
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db, logger: logrus.StandardLogger()}
}

func (r *repository) AddTable(capacity, allowedExtras int) (*types.Table, error) {
	table := types.NewTable(capacity, allowedExtras)
	result, err := r.db.NamedExecContext(context.Background(), insertTableQuery, table)
	if err != nil {
		r.logger.WithError(err).
			Error("failed to insert table")
		return nil, errors.New(http.StatusInternalServerError, "unable to create table at this time. please retry later")
	}

	newId, _ := result.LastInsertId()
	r.logger.WithField("create_table_result", newId).
		Info("table created")

	return r.GetTable(int(newId))
}

func (r *repository) AddGuest(name string, tableId, accompanyingGuests int) (*types.Guest, error) {
	guest := types.NewGuest(name, tableId, accompanyingGuests)
	result, err := r.db.NamedExecContext(context.Background(), insertGuestQuery, guest)
	if err != nil {
		r.logger.WithError(err).
			Error("failed to insert guest")
		return nil, errors.New(http.StatusInternalServerError, "unable to add guest at this time. please retry later")
	}

	newId, _ := result.LastInsertId()
	r.logger.WithField("add_guest_result", newId).
		Info("guest_added")
	return guest, nil
}

func (r *repository) GetTable(id int) (*types.Table, error) {
	table := &types.Table{}
	err := r.db.Get(table, selectTableQuery, id)
	if err == sql.ErrNoRows {
		return nil, err
	}

	if err != nil {
		r.logger.WithField("table_id", id).WithError(err).
			Error("failed to select a table")
		return nil, errors.New(http.StatusInternalServerError, "unable to select table")
	}
	return table, err
}

func (r *repository) GetGuestByTableId(tableId int) (*types.Guest, error) {
	guest := &types.Guest{}
	err := r.db.Get(guest, selectGuestsByTableIdQuery, tableId)
	if err != nil {
		return nil, err
	}

	return guest, err
}

func (r *repository) GetGuestByName(name string) (*types.Guest, error) {
	guest := &types.Guest{}
	err := r.db.Get(guest, selectGuestsByNameQuery, name)
	if err != nil {
		return nil, err
	}

	return guest, err
}

func (r *repository) UpdateGuest(name string, accompanyingGuests int) error {
	result, err := r.db.Exec(updateGuestQuery, accompanyingGuests, name)
	if err != nil {
		r.logger.WithError(err).
			Error("failed to update guest")
		return errors.New(http.StatusInternalServerError, "unable to onboard guests at this time, please retry later")
	}

	r.logger.WithField("update_guest_result", result).
		Info("guest updated")
	return err
}

func (r *repository) RemoveGuests(name string) error {
	result, err := r.db.Exec(removeGuestQuery, name)
	if err != nil {
		r.logger.WithError(err).Error("failed to complete update")
		return errors.New(http.StatusBadRequest, "failed to complete update at this time. please retry later")
	}

	r.logger.WithField("remove_guest_result", result).Info("guest removed")
	return nil
}

func (r *repository) GetAllTables() ([]*types.Table, error) {
	tables := make([]*types.Table, 0)
	err := r.db.Select(&tables, queryTables)
	if err == sql.ErrNoRows {
		return nil, errors.New(http.StatusOK, "no table found")
	}

	if err != nil {
		r.logger.WithError(err).
			Error("error querying tables")
		return nil, errors.New(http.StatusInternalServerError, "error searching tables")
	}

	return tables, err
}

func (r *repository) GetAllGuests() ([]*types.Guest, error) {
	guests := make([]*types.Guest, 0)
	err := r.db.Select(&guests, queryGuests)
	if err == sql.ErrNoRows {
		return nil, errors.New(http.StatusOK, "no guest found")
	}

	if err != nil {
		r.logger.WithError(err).
			Error("error querying guests")
		return nil, errors.New(http.StatusInternalServerError, "error searching guests")
	}

	r.logger.Infof("Guests=%v", guests)
	return guests, err
}
