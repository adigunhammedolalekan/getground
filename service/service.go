package service

import (
	"database/sql"
	"fmt"
	"github.com/getground/tech-tasks/backend/errors"
	"github.com/getground/tech-tasks/backend/repository"
	"github.com/getground/tech-tasks/backend/types"
	"github.com/sirupsen/logrus"
	"net/http"
)

type (
	Service interface {
		AddTable(capacity, allowedExtras int) (*types.Table, error)
		AddGuest(name string, tableId, accompanyingGuests int) (*types.Guest, error)
		GetArrivedGuests() ([]*types.Guest, error)
		GetGuestsList() ([]*types.Guest, error)
		GuestArrives(name string, accompanyingGuests int) (*types.Guest, error)
		GuestsLeave(name string) error
		AvailableSeats() (int, error)
	}

	service struct {
		r      repository.Repository
		logger *logrus.Logger
	}
)

func NewService(r repository.Repository) Service {
	return &service{r: r, logger: logrus.StandardLogger()}
}

func (s *service) AddTable(capacity, allowedExtras int) (*types.Table, error) {
	if capacity <= 0 || allowedExtras < 0 {
		return nil, errors.New(http.StatusBadRequest, "invalid capacity or allow extras")
	}
	return s.r.AddTable(capacity, allowedExtras)
}

func (s *service) AddGuest(name string, tableId, accompanyingGuests int) (*types.Guest, error) {
	if accompanyingGuests <= 0 {
		return nil, errors.New(http.StatusBadRequest, "Accompanying guests size must not be 0")
	}
	table, err := s.r.GetTable(tableId)
	if err == sql.ErrNoRows {
		return nil, errors.New(http.StatusNotFound, "table does not exists")
	}

	if err != nil {
		s.logger.WithError(err).
			Error("failed to fetch table")
		return nil, errors.New(http.StatusInternalServerError, "failed to find table at this time, please retry later")
	}

	if accompanyingGuests > (table.Capacity + table.AllowedExtras) {
		return nil, errors.New(http.StatusBadRequest, fmt.Sprintf("no seats available for table: %d", table.Id))
	}

	return s.r.AddGuest(name, table.Id, accompanyingGuests)
}

func (s *service) GuestArrives(name string, accompanyingGuests int) (*types.Guest, error) {
	guest, err := s.r.GetGuestByName(name)
	if err == sql.ErrNoRows {
		return nil, errors.New(http.StatusNotFound, fmt.Sprintf("guest with name %s not found", name))
	}

	if err != nil {
		s.logger.WithError(err).
			Error("failed to fetch guest")
		return nil, errors.New(http.StatusInternalServerError, "failed to find guest at this time, please retry later")
	}

	table, err := s.r.GetTable(guest.TableId)
	if err != nil {
		return nil, err
	}

	newTotalGuests := guest.AccompanyingGuests + accompanyingGuests
	if newTotalGuests > (table.Capacity + table.AllowedExtras) {
		return nil, errors.New(http.StatusOK, "this table is has reached the maximum allowed guests")
	}

	if err := s.r.UpdateGuest(name, newTotalGuests); err != nil {
		return nil, err
	}

	return guest, nil
}

func (s *service) GetArrivedGuests() ([]*types.Guest, error) {
	return s.r.GetAllGuests()
}

func (s *service) GetGuestsList() ([]*types.Guest, error) {
	return s.r.GetAllGuests()
}

func (s *service) GuestsLeave(name string) error {
	guest, err := s.r.GetGuestByName(name)
	if err != nil {
		return err
	}

	return s.r.RemoveGuests(guest.Name)
}

func (s *service) AvailableSeats() (int, error) {
	guests, err := s.r.GetAllGuests()
	if err != nil {
		return 0, err
	}

	tables, err := s.r.GetAllTables()
	if err != nil {
		return 0, err
	}

	totalGuests, totalTableCapacity := 0, 0
	for _, table := range tables {
		totalTableCapacity += table.Capacity + table.AllowedExtras
	}
	for _, guest := range guests {
		totalGuests += guest.AccompanyingGuests
	}

	return totalTableCapacity - totalGuests, nil
}
