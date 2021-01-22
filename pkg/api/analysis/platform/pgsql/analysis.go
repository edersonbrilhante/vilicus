package pgsql

import (
	"github.com/edersonbrilhante/ccvs"
	"github.com/go-pg/pg/v10/orm"
)

// Analysis represents the client for analysis table
type Analysis struct{}

// Create creates a new analysis on database
func (a Analysis) Create(db orm.DB, al ccvs.Analysis) (ccvs.Analysis, error) {
	_, err := db.Model(&al).Insert()

	return al, err
}

// View returns single analysis by ID
func (a Analysis) View(db orm.DB, id string) (ccvs.Analysis, error) {
	al := ccvs.Analysis{ID: id}
	err := db.Model(&al).WherePK().Select()

	return al, err
}

// Update updates analysis's info
func (a Analysis) Update(db orm.DB, al ccvs.Analysis) error {
	return db.Update(&al)
}


