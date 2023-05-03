# my_recon_go

`my_recon_go` is a Golang tool for subdomain enumeration and HTTP probing, loosely based on the workflow described in [this article](https://dhiyaneshgeek.github.io/bug/bounty/2020/02/06/recon-with-me/). The tool uses SubFinder for subdomain enumeration and HTTPX for HTTP probing. It stores the discovered subdomains and active hosts in a SQLite database.

## Database Schema

The SQLite database contains two tables:

1. **subdomains**: This table stores the discovered subdomains.

| Column     | Type    | Description                      |
|------------|---------|----------------------------------|
| host       | TEXT    | Subdomain host (PK)              |
| source     | TEXT    | Source of the subdomain          |
| created_at | INTEGER | Timestamp of creation (UTC)      |
| updated_at | INTEGER | Timestamp of last update (UTC)   |

2. **active_hosts**: This table stores the active hosts found during HTTP probing.

| Column          | Type    | Description                     |
|-----------------|---------|---------------------------------|
| id              | INTEGER | Unique identifier (PK)          |
| subdomain       | TEXT    | Subdomain host                  |
| method          | TEXT    | HTTP method used                |
| url             | TEXT    | URL of the active host          |
| status_code     | INTEGER | HTTP status code                |
| title           | TEXT    | Webpage title                   |
| technologies    | TEXT    | Detected technologies (JSON)    |
| created_at      | INTEGER | Timestamp of creation (UTC)      |
| updated_at      | INTEGER | Timestamp of last update (UTC)  |
| last_nuclei_scan| INTEGER | Timestamp of last Nuclei scan (UTC) |

## Setup

1. Set up the configuration in `.config.yaml`:

```yaml
MY_RECON_DB: .project.db
MY_RECON_DOMAINS:
  - example.gov
  - example.org
```

This configuration file specifies the SQLite database file and the list of domains to run the recon process on.

2. Since there is a `go.mod` file in the root of the package/module, the required dependencies will be automatically downloaded and installed when you build or run the project.

## Run

To run the tool, execute the following command:

```sh
go run main.go
```

This command will run the recon process, which includes initializing the database, running subdomain enumeration using SubFinder, and performing HTTP probing using HTTPX. The discovered subdomains and active hosts will be stored in the SQLite database specified in the `.config.yaml` file.