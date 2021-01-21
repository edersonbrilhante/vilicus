package analysis

import (
	"github.com/labstack/echo/v4"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/edersonbrilhante/ccvs/pkg/api/analysis/platform/pgsql"
	"github.com/edersonbrilhante/ccvs"
)

// Service represents analysis application interface
type Service interface {
	Create(echo.Context, ccvs.Analysis) (ccvs.Analysis, error)
	View(echo.Context, string) (ccvs.Analysis, error)
}

// Repository represents analysis repository interface
type Repository interface {
	Create(orm.DB, ccvs.Analysis) (ccvs.Analysis, error)
	View(orm.DB, string) (ccvs.Analysis, error)
}

// Analysis represents analysis application service
type Analysis struct {
	db  *pg.DB
	repository Repository
}

// Create creates a new Analysis
func (a Analysis) Create(c echo.Context, req ccvs.Analysis) (ccvs.Analysis, error) {
	return a.repository.Create(a.db, req)
}

// View returns single Analysis
func (a Analysis) View(c echo.Context, id string) (ccvs.Analysis, error) {
	return a.repository.View(a.db, id)
}

// New creates new analysis application service
func New(db *pg.DB, repository Repository) *Analysis {
	return &Analysis{db: db, repository: repository}
}

// Initialize initalizes Analysis application service with defaults
func Initialize(db *pg.DB) *Analysis {
	return New(db, pgsql.Analysis{})
}

