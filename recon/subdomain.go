package recon

import (
	"time"

	"github.com/pkg/errors"
)

type Subdomain struct {
	Host      string    `json:"host"`
	Source    string    `json:"source"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SubdomainModel struct {
	Database SQLDatabase
}

const createSubdomainsTableSQL = `
	CREATE TABLE IF NOT EXISTS subdomains (
		host TEXT PRIMARY KEY NOT NULL,
		source TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
`

func (m SubdomainModel) InitDB() error {
	_, err := m.Database.Exec(createSubdomainsTableSQL)
	if err != nil {
		return errors.Wrap(err, "creating table")
	}

	return nil
}

const selectAllSubdomainsSQL = "SELECT * FROM subdomains"

func (m SubdomainModel) All() ([]Subdomain, error) {
	rows, err := m.Database.Query(selectAllSubdomainsSQL)
	if err != nil {
		return nil, errors.Wrap(err, "querying subdomains")
	}
	defer rows.Close()

	var subs []Subdomain

	for rows.Next() {
		var sub Subdomain

		err := rows.Scan(&sub.Host, &sub.Source, &sub.CreatedAt, &sub.UpdatedAt)
		if err != nil {
			return nil, errors.Wrap(err, "scanning subdomain")
		}

		subs = append(subs, sub)
	}
	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "processing rows")
	}

	return subs, nil
}

const insertSubdomainSQL = "INSERT INTO subdomains (host, source) VALUES ($1, $2);"

func (m SubdomainModel) Create(sub *Subdomain) error {
	_, err := m.Database.Exec(insertSubdomainSQL, sub.Host, sub.Source)
	if err != nil {
		return errors.Wrap(err, "inserting subdomain")
	}

	return nil
}
