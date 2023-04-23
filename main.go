package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"

	"github.com/atypicaltech/my_recon_go/config"
	"github.com/atypicaltech/my_recon_go/recon"
)

func main() {
	cfg := config.Get()
	dbFile := cfg.GetString(config.DB_FILE)
	domainConfig := cfg.GetStringSlice(config.DOMAIN_CONFIG)

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	recon := &recon.Recon{
		DomainConfig: domainConfig,
		Subdomains:   recon.SubdomainModel{DB: db},
	}

	err = recon.Init()
	if err != nil {
		log.Fatal(err)
	}

	err = recon.ScapeSubdomains()
	if err != nil {
		log.Fatal(err)
	}

	subs, err := recon.Subdomains.All()
	if err != nil {
		log.Fatal(err)
	}

	for _, sub := range subs {
		fmt.Println(sub.Host)
	}
}
