package api

import (
	"os"

	"github.com/edersonbrilhante/vilicus/pkg/api/analysis"
	al "github.com/edersonbrilhante/vilicus/pkg/api/analysis/logging"
	at "github.com/edersonbrilhante/vilicus/pkg/api/analysis/transport"
	ht "github.com/edersonbrilhante/vilicus/pkg/api/healthcheck/transport"
	"github.com/edersonbrilhante/vilicus/pkg/util/config"
	"github.com/edersonbrilhante/vilicus/pkg/util/postgres"
	"github.com/edersonbrilhante/vilicus/pkg/util/server"
	"github.com/edersonbrilhante/vilicus/pkg/util/zlog"
)

// Start starts the API service
func Start(cfg *config.Configuration) error {

	// db connection
	db, err := postgres.New(os.Getenv("DATABASE_URL"), cfg.DB.Timeout, cfg.DB.LogQueries)
	if err != nil {
		return err
	}

	log := zlog.New()

	e := server.New()
	v1 := e.Group("")

	at.NewHTTP(al.New(analysis.Initialize(db, cfg.Vendors), log), v1)
	ht.NewHTTP(v1)

	server.Start(e, &server.Config{
		Port:                cfg.Server.Port,
		ReadTimeoutSeconds:  cfg.Server.ReadTimeout,
		WriteTimeoutSeconds: cfg.Server.WriteTimeout,
		Debug:               cfg.Server.Debug,
	})

	return nil
}
