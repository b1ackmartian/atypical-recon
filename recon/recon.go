package recon

import (
	"errors"
	"io"

	_ "github.com/mattn/go-sqlite3"
	"github.com/projectdiscovery/goflags"
	httpxRunner "github.com/projectdiscovery/httpx/runner"
	"github.com/projectdiscovery/subfinder/v2/pkg/resolve"
	subfinderRunner "github.com/projectdiscovery/subfinder/v2/pkg/runner"
	"go.uber.org/zap"
)

type Recon struct {
	DomainConfig []string
	Subdomains   SubdomainRepository
	ActiveHosts  ActiveHostRepository
	Logger       *zap.Logger
}

type SubdomainRepository interface {
	InitDB() error
	All() ([]Subdomain, error)
	Create(subdomain *Subdomain) error
}

type ActiveHostRepository interface {
	InitDB() error
	All() ([]ActiveHost, error)
	Create(activeHost *ActiveHost) error
}

func (r *Recon) Init() error {
	if r.DomainConfig == nil || len(r.DomainConfig) == 0 {
		return errors.New("domain config is not set")
	}

	if err := r.Subdomains.InitDB(); err != nil {
		return err
	}

	if err := r.ActiveHosts.InitDB(); err != nil {
		return err
	}

	return nil
}

func (r *Recon) ScapeSubdomains() error {
	runnerInstance, err := subfinderRunner.NewRunner(&subfinderRunner.Options{
		ResultCallback: r.handleSubdomainResult,
		Domain:         r.DomainConfig,
		Output:         io.Discard,
	})
	if err != nil {
		return err
	}

	return runnerInstance.RunEnumeration()
}

func (r *Recon) handleSubdomainResult(s *resolve.HostEntry) {
	if err := r.Subdomains.Create(&Subdomain{
		Host:   s.Host,
		Source: s.Source,
	}); err != nil {
		r.Logger.Error("problem adding subdomain",
			zap.String("host", s.Host),
			zap.String("source", s.Source),
			zap.Error(err),
		)
	}
}

func (r *Recon) HTTPProbe() error {
	hosts, err := r.getHostsFromSubdomains()
	if err != nil {
		return err
	}

	options := httpxRunner.Options{
		Silent:              true,
		InputTargetHost:     goflags.StringSlice(hosts),
		Methods:             "GET,HEAD",
		StatusCode:          true,
		ExtractTitle:        true,
		TechDetect:          true,
		FollowRedirects:     true,
		FollowHostRedirects: true,
		OnResult:            r.handleHTTPXResult,
	}

	if err := options.ValidateOptions(); err != nil {
		return err
	}

	httpxRunnerInstance, err := httpxRunner.New(&options)
	if err != nil {
		return err
	}
	defer httpxRunnerInstance.Close()

	httpxRunnerInstance.RunEnumeration()

	return nil
}

func (r *Recon) getHostsFromSubdomains() ([]string, error) {
	subs, err := r.Subdomains.All()
	if err != nil {
		return nil, err
	}

	hosts := make([]string, len(subs))
	for i, sub := range subs {
		hosts[i] = sub.Host
	}

	return hosts, nil
}

func (r *Recon) handleHTTPXResult(result httpxRunner.Result) {
	if err := r.ActiveHosts.Create(&ActiveHost{
		Subdomain:    result.Input,
		Method:       result.Method,
		URL:          result.URL,
		StatusCode:   result.StatusCode,
		Title:        result.Title,
		Technologies: result.Technologies,
	}); err != nil {
		r.Logger.Error("problem adding active host",
			zap.String("input", result.Input),
			zap.String("method", result.Method),
			zap.String("url", result.URL),
			zap.Error(err),
		)
	}
}
