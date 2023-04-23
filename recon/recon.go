package recon

import (
	"fmt"
	"io"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/projectdiscovery/subfinder/v2/pkg/resolve"
	"github.com/projectdiscovery/subfinder/v2/pkg/runner"
)

type Recon struct {
	DomainConfig []string
	Subdomains   interface {
		InitDB() error
		All() ([]Subdomain, error)
		Create(book *Subdomain) error
	}
}

func (r *Recon) Init() error {
	if r.DomainConfig == nil {
		return fmt.Errorf("domain config is not set")
	}

	err := r.Subdomains.InitDB()
	if err != nil {
		return err
	}

	return nil
}

func (r *Recon) ScapeSubdomains() error {
	runnerInstance, err := runner.NewRunner(&runner.Options{
		Threads:            10,                       // Thread controls the number of threads to use for active enumerations
		Timeout:            30,                       // Timeout is the seconds to wait for sources to respond
		MaxEnumerationTime: 10,                       // MaxEnumerationTime is the maximum amount of time in mins to wait for enumeration
		Resolvers:          resolve.DefaultResolvers, // Use the default list of resolvers by marshaling it to the config
		ResultCallback: func(s *resolve.HostEntry) {
			err := r.Subdomains.Create(&Subdomain{
				Host:   s.Host,
				Source: s.Source,
			})
			if err != nil {
				log.Printf("problem adding subdomain [%s] from source [%s]: %s", s.Host, s.Source, err)
			}
		},
		Domain: r.DomainConfig,
		Output: io.Discard,
	})
	if err != nil {
		return err
	}

	err = runnerInstance.RunEnumeration()
	if err != nil {
		return err
	}

	return nil
}
