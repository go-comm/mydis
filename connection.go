package mydis

type DriverConnection interface {
	Close() error
}

type driverConnection struct {
}

type pooledDriverConnection struct {
}
