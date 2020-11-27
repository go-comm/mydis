package mydis

type Driver interface {
	Open(params interface{}) (DriverConnection, error)
	Close(DriverConnection) error
}
