package types

import (
	"database/sql"
	"time"
)

type (
	Table struct {
		Id            int          `json:"id" db:"id"`
		Capacity      int          `json:"capacity" db:"capacity"`
		AllowedExtras int          `json:"-" db:"allowed_extras"`
		CreatedAt     time.Time    `json:"-" db:"created_at"`
		DeletedAt     sql.NullTime `json:"-" db:"deleted_at"`
	}

	Guest struct {
		Id                 int          `json:"-" db:"id"`
		Name               string       `json:"name" db:"name"`
		TableId            int          `json:"table" db:"table_id"`
		AccompanyingGuests int          `json:"accompanying_guests" db:"accompanying_guests"`
		CreatedAt          time.Time    `json:"-" db:"created_at"`
		DeletedAt          sql.NullTime `json:"-" db:"deleted_at"`
	}
)

func NewTable(capacity, allowedExtras int) *Table {
	return &Table{
		Capacity:      capacity,
		AllowedExtras: allowedExtras,
		CreatedAt:     time.Now().UTC(),
	}
}

func NewGuest(name string, tableId, accompanyingGuests int) *Guest {
	return &Guest{
		Name:               name,
		TableId:            tableId,
		AccompanyingGuests: accompanyingGuests,
		CreatedAt:          time.Now().UTC(),
	}
}
