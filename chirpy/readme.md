#### Install Postgresql tooling
1. `sudo apt update`
2. `sudo apt install postgresql postgresql-contrib`
3. `psql --version`
4. Setup password for system user: `sudo passwd postgres`

#### Start Postgresql Server
1. `sudo service postgresql start`
2. Connect to server: `sudo -u postgres psql`
3. Create a new database: `CREATE DATABASE <database_name>;`
4. Setup password for db user: `ALTER USER postgres WITH PASSWORD '<password>';`

#### Setup Goose (database migration scripts)
1. `go install github.com/pressly/goose/v3/cmd/goose@latest`
2. Setup folder structure: `mkdir sql && mkdir ./sql/schema`
3. Each migration should be prefixed with `00x_` where `x = migration number`
4. To run a goose migration
    1. Navigate to the schema dir: `cd ./sql/schema`
    2. Get the postgres connect url, which is of the form: `protocol://username:password@host:port/database`
    3. Normally this is in `./env`
    4. Verify the url by running: `psql "postgres://<username>:<password>@localhost:5432/<database_name>"`
    5. Run the goose migration: `goose postgres <connection_string> up`

### Setup SQLC (Go-sql helper functions)
1. `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`
2. Create the sqlc.yaml file, should look something like:
```
version: "2"
sql:
  - schema: "sql/schema"
    queries: "sql/queries"
    engine: "postgresql"
    gen:
      go:
        out: "internal/database"
```

3. Create queries folder (per above config): `mkdir ./sql/queries`
4. Create the output folder (where the generated go code will live, per above config): `mkdir internal && mkdir ./internal/database`
5. Run generate whenever you add a new query: `sqlc generate`


### Setup Supporting Packages
1. Get the UUI package: `go get github.com/google/uuid`
2. Get the postgresql package: `go get github.com/lib/pq`
3. Dump the pq package at the top of `main.go`
4. Get the Go .env package: `go get github.com/joho/godotenv`
    1. Now we can reference the .env DB _URL variable with: `dbURL := os.Getenv("DB_URL")`
    2. Subsequently, we can open a db connection with `db, err := sql.Open("postgres", dbURL)`
    3. Finally, we can setup queries: `dbQueries := database.New(db)`
5. Get Argon2 (used for hashing passwords): `go get github.com/alexedwards/argon2id`
