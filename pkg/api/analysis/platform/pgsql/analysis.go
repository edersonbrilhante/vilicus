package pgsql

import (
	"time"

	"github.com/go-pg/pg/v10/orm"

	"github.com/edersonbrilhante/vilicus"
)

// Analysis represents the client for analysis table
type Analysis struct{}

// Create creates a new analysis on database
func (a Analysis) Create(db orm.DB, al vilicus.Analysis) (vilicus.Analysis, error) {
	_, err := db.Model(&al).Insert()

	return al, err
}

// View returns single analysis by ID
func (a Analysis) View(db orm.DB, id string) (vilicus.Analysis, error) {
	al := vilicus.Analysis{ID: id}
	err := db.Model(&al).WherePK().Select()

	return al, err
}

// Update updates analysis's info
func (a Analysis) Update(db orm.DB, al vilicus.Analysis) (vilicus.Analysis, error) {
	al.UpdatedAt = time.Now()
	_, err := db.Model(&al).WherePK().Update()

	return al, err
}
