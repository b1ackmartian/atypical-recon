package recon

import "database/sql"

type Subdomain struct {
	Host   string `json:"host"`
	Source string `json:"source"`
}

type SubdomainModel struct {
	DB *sql.DB
}

func (m SubdomainModel) InitDB() error {
	createTableSQL := `
		CREATE TABLE IF NOT EXISTS subdomains (
			host TEXT PRIMARY KEY NOT NULL,
			source TEXT NOT NULL
		);
	`

	_, err := m.DB.Exec(createTableSQL)
	if err != nil {
		return err
	}

	return nil
}

func (m SubdomainModel) All() ([]Subdomain, error) {
	stmt, err := m.DB.Prepare("SELECT * FROM subdomains")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []Subdomain

	for rows.Next() {
		var sub Subdomain

		err := rows.Scan(&sub.Host, &sub.Source)
		if err != nil {
			return nil, err
		}

		subs = append(subs, sub)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return subs, nil
}

func (m SubdomainModel) Create(sub *Subdomain) error {
	stmt, err := m.DB.Prepare("INSERT INTO subdomains (host, source) VALUES ($1, $2);")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(sub.Host, sub.Source)
	if err != nil {
		return err
	}

	return nil
}
