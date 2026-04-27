# migrate

A fork of [golang-migrate/migrate](https://github.com/golang-migrate/migrate) — Database migrations written in Go. Use as CLI or import as library.

[![Go Reference](https://pkg.go.dev/badge/github.com/your-org/migrate.svg)](https://pkg.go.dev/github.com/your-org/migrate)
[![CI](https://github.com/your-org/migrate/actions/workflows/ci.yaml/badge.svg)](https://github.com/your-org/migrate/actions/workflows/ci.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/your-org/migrate)](https://goreportcard.com/report/github.com/your-org/migrate)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Features

- **Database drivers**: PostgreSQL, MySQL, SQLite, MongoDB, CockroachDB, and more
- **Source drivers**: File system, Go embed, GitHub, S3, GCS
- **CLI tool**: Easy to use command-line interface
- **Library**: Import and use directly in your Go application
- **Graceful error handling**: Dirty state detection and recovery

## Installation

### CLI

```bash
# Using Go install
go install github.com/your-org/migrate/cmd/migrate@latest

# Using Homebrew
brew install migrate
```

### Library

```bash
go get github.com/your-org/migrate/v4
```

## Quick Start

### CLI Usage

```bash
# Run all up migrations
migrate -path ./migrations -database "postgres://localhost:5432/mydb?sslmode=disable" up

# Rollback last migration
migrate -path ./migrations -database "postgres://localhost:5432/mydb?sslmode=disable" down 1

# Check current version
migrate -path ./migrations -database "postgres://localhost:5432/mydb?sslmode=disable" version

# Force set version (use with caution)
migrate -path ./migrations -database "postgres://localhost:5432/mydb?sslmode=disable" force 1
```

### Library Usage

```go
package main

import (
    "log"

    "github.com/your-org/migrate/v4"
    _ "github.com/your-org/migrate/v4/database/postgres"
    _ "github.com/your-org/migrate/v4/source/file"
)

func main() {
    m, err := migrate.New(
        "file://./migrations",
        "postgres://localhost:5432/mydb?sslmode=disable",
    )
    if err != nil {
        log.Fatal(err)
    }
    defer m.Close()

    // Set a lock timeout to avoid hanging indefinitely on busy databases
    m.LockTimeout = 15 * time.Second

    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        log.Fatal(err)
    }
}
```

## Migration Files

Migration files follow the naming convention:

```
{version}_{title}.up.{extension}
{version}_{title}.down.{extension}
```

Example:
```
1_create_users.up.sql
1_create_users.down.sql
2_add_email_index.up.sql
2_add_email_index.down.sql
```

## Supported Databases

| Database | Driver Import |
|----------|---------------|
| PostgreSQL | `github.com/your-org/migrate/v4/database/postgres` |
| MySQL | `github.com/your-org/migrate/v4/database/mysql` |
| SQLite | `github.com/your-org/migrate/v4/database/sqlite3` |
| MongoDB | `github.com/your-org/migrate/v4/database/mongodb` |
| CockroachDB | `github.com/your-org/migrate/v4/database/cockroachdb` |

## Development

```bash
# Clone the repository
git clone https://github.com/your-org/migrate.git
cd migrate

# Run tests
go test ./...
```
