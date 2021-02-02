package analysis

import (
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/labstack/echo/v4"

	"github.com/edersonbrilhante/vilicus"
	"github.com/edersonbrilhante/vilicus/pkg/api/analysis/platform/pgsql"
	"github.com/edersonbrilhante/vilicus/pkg/util/config"
	"github.com/edersonbrilhante/vilicus/analyzer"
)

// Service represents analysis application interface
type Service interface {
	Create(echo.Context, vilicus.Analysis) (vilicus.Analysis, error)
	View(echo.Context, string) (vilicus.Analysis, error)
	Update(echo.Context, vilicus.Analysis) (vilicus.Analysis, error)
}

// Repository represents analysis repository interface
type Repository interface {
	Create(orm.DB, vilicus.Analysis) (vilicus.Analysis, error)
	View(orm.DB, string) (vilicus.Analysis, error)
	Update(orm.DB, vilicus.Analysis) (vilicus.Analysis, error)
}

// Analysis represents analysis application service
type Analysis struct {
	db         *pg.DB
	repository Repository
	vendors    *config.Vendors
}

// Create creates a new Analysis
func (a Analysis) Create(c echo.Context, req vilicus.Analysis) (vilicus.Analysis, error) {
	return a.repository.Create(a.db, req)
}

// View returns single Analysis
func (a Analysis) View(c echo.Context, id string) (vilicus.Analysis, error) {
	return a.repository.View(a.db, id)
}

// Update updates a Analysis
func (a Analysis) Update(c echo.Context, req vilicus.Analysis) (vilicus.Analysis, error) {
	req.Status = "started"
	req, err := a.repository.Update(a.db, req)
	if err != nil {
		return req, err
	}
	analyzer.StartAnalysis(a.vendors, &req)
	req.Status = "finished"
	return a.repository.Update(a.db, req)
}

// New creates new analysis application service
func New(db *pg.DB, vendors *config.Vendors, repository Repository) *Analysis {
	return &Analysis{db: db, vendors: vendors, repository: repository}
}

// Initialize initalizes Analysis application service with defaults
func Initialize(db *pg.DB, vendors *config.Vendors) *Analysis {
	return New(db, vendors, pgsql.Analysis{})
}
