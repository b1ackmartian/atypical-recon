# my_recon_go

loosely based off what's happening at [this article](https://dhiyaneshgeek.github.io/bug/bounty/2020/02/06/recon-with-me/)

## setup

Set configuration in `.config.yaml`:

```yaml
MY_RECON_DB: .project.db
MY_RECON_DOMAINS:
  - example.gov
  - example.org
```

## run

```sh
go run main.go
```
