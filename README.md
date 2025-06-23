
# PostgreSQL
## Install PostgreSQL
```
sudo apt update
sudo apt install postgresql postgresql-contrib
```

## Update postgres password
`sudo passwd postgres`

## Set user password
```
sudo -u postgres psql
ALTER USER postgres PASSWORD 'postgres';
```
Exit with `\q` or `exit`.

## Create database 'gator'
```
CREATE DATABASE gator;
```

# Install Goose
Goose https://github.com/pressly/goose
```
go install github.com/pressly/goose/v3/cmd/goose@latest
```

# Connection string
`postgres://postgres:postgres@localhost:5432/gator`

Run migrations to create the tables:
```
cd sql/schema
goose postgres postgres://postgres:postgres@localhost:5432/gator up
cd ../..
```

# Config file
Create `~/.gatorconfig.json` with the connection string plus `?sslmode=disable`:
```
{
  "db_url": "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable"
}
```

The program will save config information in this JSON formatted file.


# Other
SQLC https://sqlc.dev/
```
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

Aixo no cal? S'instal·la sol quan fem `go run .`
```
go get github.com/google/uuid
go get github.com/lib/pq
```