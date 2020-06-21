package storage

import (
	"database/sql/driver"
	"github.com/nrocco/qb"
)

// Tags is a slice of string values
type Tags []string

// Value implements the Valuer interface
func (t Tags) Value() (driver.Value, error) {
	return qb.JSONValue(t)
}

// Scan implements the Scanner interface
func (t *Tags) Scan(value interface{}) error {
	return qb.JSONScan(t, value)
}
