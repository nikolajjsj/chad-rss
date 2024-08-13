# Project chad-rss

One Paragraph of project description goes here

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

## MakeFile

run all make commands with clean tests
```bash
make all build
```

build the application
```bash
make build
```

run the application
```bash
make run
```

Create DB container
```bash
make docker-run
```

Shutdown DB container
```bash
make docker-down
```

live reload the application
```bash
make watch
```

run the test suite
```bash
make test
```

clean up binary from the last build
```bash
make clean
```

# chad-rss

## 👀 Usage

```
#### 3. You have to migrate the database.
> ##### 🎯 It is a "database-first" ORM as opposed to "code-first" (like gorm/gorp). That means you must first create your database schema.
> ##### 🎯 I used [golang-migrate](https://github.com/golang-migrate/migrate) to proceed with the migrate.
###### 1. Make Migration files
```bash
$ migrate create -ext sql -dir ./database/migrations -seq create_initial_table
```
```console
sqlc/database/migrations/000001_create_initial_table.up.sql
sqlc/database/migrations/000001_create_initial_table.up.sql
```
###### 2. Migrate
```bash
$ migrate -path database/migrations -database "postgresql://user:password@localhost:5432/fiber_demo?sslmode=disable" -verbose up
```
```console
2023/09/28 20:00:00 Start buffering 1/u create_initial_table
2023/09/28 20:00:00 Read and execute 1/u create_initial_table
2023/09/28 20:00:00 Finished 1/u create_initial_table (read 24.693541ms, ran 68.30925ms)
2023/09/28 20:00:00 Finished after 100.661625ms
2023/09/28 20:00:00 Closing source and database
```
###### 3. Rollback Migrate
```bash
$ migrate -path database/migrations -database "postgresql://user:password@localhost:5432/fiber_demo?sslmode=disable" -verbose down
```
```console
2023/09/28 20:00:00 Are you sure you want to apply all down migrations? [y/N]
y
2023/09/28 20:00:00 Applying all down migrations
2023/09/28 20:00:00 Start buffering 1/d create_initial_table
2023/09/28 20:00:00 Read and execute 1/d create_initial_table
2023/09/28 20:00:00 Finished 1/d create_initial_table (read 39.681125ms, ran 66.220125ms)
2023/09/28 20:00:00 Finished after 1.83152475s
```
#### 4. Use sqlc
###### 1. Install
```bash
# Go 1.17 and above:
$ go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```
###### 2. Create a configuration file
###### Example
###### sqlc.yaml
```yaml
version: "2"
sql:
  - engine: "postgresql"
    queries: "database/query"
    schema: "database/migrations"
    gen:
      go:
        package: "sqlc"
        out: "database/sqlc"
```
###### post.sql
```sql
-- name: GetPosts :many
SELECT * FROM post;

-- name: GetPost :one
SELECT * FROM post WHERE id = $1;

-- name: NewPost :one
INSERT INTO post (title, content, author) VALUES ($1, $2, $3) RETURNING *;

-- name: UpdatePost :one
UPDATE post SET title = $1, content = $2, author = $3 WHERE id = $4 RETURNING *;

-- name: DeletePost :exec
DELETE FROM post WHERE id = $1;

```
###### 3. Generate
```bash
$ sqlc generate
```
```text
sqlc/
├── db.go
├── models.go
├── post.sql.go
```
#### 5. Reference
[sqlc document](https://docs.sqlc.dev/en/stable/)
