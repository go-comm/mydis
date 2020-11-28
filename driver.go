package mydis

import (
	"sync"
)

var (
	driversMu sync.RWMutex
	drivers   = make(map[string]Driver)
)

type Driver interface {
	Open(dsn string) (Conn, error)
}

func Register(name string, driver Driver) {
	driversMu.Lock()
	defer driversMu.Unlock()
	if driver == nil {
		panic("mydis: Register driver is nil")
	}
	if _, dup := drivers[name]; dup {
		panic("mydis: Register called twice for driver " + name)
	}
	drivers[name] = driver
}

func Open(driverName, dsn string) (*DB, error) {
	db := &DB{}
	return db, nil
}
