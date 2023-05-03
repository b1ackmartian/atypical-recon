package recon

import (
	"database/sql"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type ActiveHost struct {
	Subdomain    string    `json:"subdomain"`
	Method       string    `json:"method"`
	URL          string    `json:"url"`
	StatusCode   int       `json:"status_code"`
	Title        string    `json:"title"`
	Technologies []string  `json:"technologies"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ActiveHostModel struct {
	Database SQLDatabase
}

const createActiveHostsTableSQL = `
	CREATE TABLE IF NOT EXISTS active_hosts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		subdomain TEXT NOT NULL,
		method TEXT NOT NULL,
		url TEXT NOT NULL,
		status_code INTEGER NOT NULL,
		title TEXT,
		technologies TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		last_nuclei_scan TIMESTAMP,
		UNIQUE(subdomain, method, url)
	);
`

func (m ActiveHostModel) InitDB() error {
	_, err := m.Database.Exec(createActiveHostsTableSQL)
	if err != nil {
		return errors.Wrap(err, "creating table")
	}

	return nil
}

const selectAllActiveHostsSQL = "SELECT * FROM active_hosts"

func (m ActiveHostModel) All() ([]ActiveHost, error) {
	rows, err := m.Database.Query(selectAllActiveHostsSQL)
	if err != nil {
		return nil, errors.Wrap(err, "querying active hosts")
	}
	defer rows.Close()

	var hosts []ActiveHost

	for rows.Next() {
		var host ActiveHost
		var technologies string

		err := rows.Scan(&host.Subdomain, &host.Method, &host.URL, &host.StatusCode, &host.Title, &technologies, &host.CreatedAt, &host.UpdatedAt)
		if err != nil {
			return nil, errors.Wrap(err, "scanning active host")
		}

		host.Technologies = parseTechnologies(technologies)
		hosts = append(hosts, host)
	}
	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "processing rows")
	}

	return hosts, nil
}

const insertActiveHostSQL = "INSERT INTO active_hosts (subdomain, method, url, status_code, title, technologies) VALUES (?, ?, ?, ?, ?, ?);"

func (m ActiveHostModel) Create(host *ActiveHost) error {
	technologies := formatTechnologies(host.Technologies)
	_, err := m.Database.Exec(insertActiveHostSQL, host.Subdomain, host.Method, host.URL, host.StatusCode, host.Title, technologies)
	if err != nil {
		return errors.Wrap(err, "inserting active host")
	}

	return nil
}

const selectActiveHostSQL = "SELECT status_code, technologies FROM active_hosts WHERE subdomain = ? AND method = ? AND url = ?"

func (m ActiveHostModel) Upsert(host *ActiveHost) error {
	row := m.Database.QueryRow(selectActiveHostSQL, host.Subdomain, host.Method, host.URL)

	var statusCode int
	var technologiesStr string
	err := row.Scan(&statusCode, &technologiesStr)
	if err != nil {
		if err == sql.ErrNoRows {
			return m.Create(host)
		}
		return errors.Wrap(err, "querying active host")
	}

	technologies := parseTechnologies(technologiesStr)
	if statusCode != host.StatusCode || !equalSlices(technologies, host.Technologies) {
		host.UpdatedAt = time.Now()
		return m.Update(host)
	}

	return nil
}

const updateActiveHostSQL = "UPDATE active_hosts SET status_code = ?, title = ?, technologies = ?, updated_at = ? WHERE subdomain = ? AND method = ? AND url = ?"

func (m ActiveHostModel) Update(host *ActiveHost) error {
	technologies := formatTechnologies(host.Technologies)
	_, err := m.Database.Exec(updateActiveHostSQL, host.StatusCode, host.Title, technologies, host.UpdatedAt, host.Subdomain, host.Method, host.URL)
	if err != nil {
		return errors.Wrap(err, "updating active host")
	}

	return nil
}

func parseTechnologies(technologies string) []string {
	if technologies == "" {
		return nil
	}
	return strings.Split(technologies, ",")
}

func formatTechnologies(technologies []string) string {
	return strings.Join(technologies, ",")
}

func equalSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
