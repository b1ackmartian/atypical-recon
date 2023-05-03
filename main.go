package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"

	"github.com/atypicaltech/my_recon_go/config"
	"github.com/atypicaltech/my_recon_go/recon"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	conf, err := config.New()
	if err != nil {
		logger.Fatal("Failed to create new config", zap.Error(err))
	}

	db, err := openDatabase(conf.GetString(config.DbFile))
	if err != nil {
		logger.Fatal("Failed to open database", zap.Error(err))
	}
	defer db.Close()

	err = runRecon(db, conf.GetStringSlice(config.DomainConfig), logger)
	if err != nil {
		logger.Fatal("Failed to run recon", zap.Error(err))
	}
}

func openDatabase(dbFile string) (*sql.DB, error) {
	return sql.Open("sqlite3", dbFile)
}

func runRecon(db recon.SQLDatabase, domainConfig []string, logger *zap.Logger) error {
	r := &recon.Recon{
		DomainConfig: domainConfig,
		Subdomains:   recon.SubdomainModel{Database: db},
		ActiveHosts:  recon.ActiveHostModel{Database: db},
		Logger:       logger,
	}

	if err := r.Init(); err != nil {
		return err
	}
	if err := r.ScapeSubdomains(); err != nil {
		return err
	}
	if err := r.HTTPProbe(); err != nil {
		return err
	}
	return nil
}
