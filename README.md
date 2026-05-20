# URL Shortener

Its a project I'm doing for fun, this is to shorten the URLs

## Setup

Running
```bash
docker compose build && docker compose up --wait
```

## Migration
Creattin new migration

```bash
  migrate create -ext sql -dir services/migration/db/migrations -seq <migration name>
```
Forcing specific version
    
```bash
  migrate -path services/migration/db/migrations -database "postgres://<user>:<password>@localhost:5430/shortner sslmode=disable" force <version>
```
