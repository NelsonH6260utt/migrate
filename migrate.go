// Package migrate provides database migration functionality.
// It is a fork of golang-migrate/migrate with additional features and fixes.
package migrate

import (
	"errors"
	"fmt"
	"os"
	"sync"
)

// ErrNoChange is returned when no migration is needed.
var ErrNoChange = errors.New("no change")

// ErrNilVersion is returned when the version is nil.
var ErrNilVersion = errors.New("no migration version found")

// ErrLocked is returned when the database is already locked.
var ErrLocked = errors.New("database locked")

// ErrLockTimeout is returned when the lock timeout is exceeded.
var ErrLockTimeout = errors.New("lock timeout")

// DefaultPrefetchMigrations is the default number of migrations to prefetch.
const DefaultPrefetchMigrations = 10

// DefaultLockTimeout is the default timeout for acquiring a database lock.
const DefaultLockTimeout = 15

// Migrate is the main struct for managing database migrations.
type Migrate struct {
	// sourceName is the registered source driver name.
	sourceName string
	// sourceDrv is the source driver instance.
	sourceDrv Source

	// databaseName is the registered database driver name.
	databaseName string
	// databaseDrv is the database driver instance.
	databaseDrv Database

	// Log is an optional logger. If nil, no logging is performed.
	Log Logger

	// GracefulStop is a channel to signal a graceful stop.
	GracefulStop chan bool
	stop         bool

	// PrefetchMigrations controls how many migrations are pre-read from source.
	PrefetchMigrations uint

	// LockTimeout controls how long to wait for a database lock (in seconds).
	LockTimeout uint

	stateMu sync.Mutex
	isLocked bool
}

// Logger is the interface for logging migration activity.
type Logger interface {
	Printf(format string, v ...interface{})
	Verbose() bool
}

// New returns a new Migrate instance from the provided source and database URLs.
func New(sourceURL, databaseURL string) (*Migrate, error) {
	m := &Migrate{
		GracefulStop:       make(chan bool, 1),
		PrefetchMigrations: DefaultPrefetchMigrations,
		LockTimeout:        DefaultLockTimeout,
	}

	sourceDrv, err := Open(sourceURL)
	if err != nil {
		return nil, fmt.Errorf("source: %w", err)
	}
	m.sourceDrv = sourceDrv

	databaseDrv, err := OpenDatabase(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("database: %w", err)
	}
	m.databaseDrv = databaseDrv

	return m, nil
}

// Close closes the source and database connections.
func (m *Migrate) Close() (source error, database error) {
	databaseSrvClose := make(chan error)
	sourceSrvClose := make(chan error)

	go func() {
		databaseSrvClose <- m.databaseDrv.Close()
	}()

	go func() {
		sourceSrvClose <- m.sourceDrv.Close()
	}()

	return <-sourceSrvClose, <-databaseSrvClose
}

// logPrintf logs a message if a logger is set.
func (m *Migrate) logPrintf(format string, v ...interface{}) {
	if m.Log != nil {
		m.Log.Printf(format, v...)
	}
}

// logVerbosePrintf logs a verbose message if a logger is set and verbose mode is enabled.
func (m *Migrate) logVerbosePrintf(format string, v ...interface{}) {
	if m.Log != nil && m.Log.Verbose() {
		m.Log.Printf(format, v...)
	}
}

// isGracefulStop returns true if a graceful stop has been requested.
func (m *Migrate) isGracefulStop() bool {
	select {
	case <-m.GracefulStop:
		return true
	default:
		return false
	}
}

// stderr writes a message to stderr.
func stderr(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format+"\n", args...)
}
