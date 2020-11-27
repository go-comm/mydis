package mydis

import "context"

func (dc *driverConnection) HGet(ctx context.Context, key []byte, field []string) (v interface{}, err error) {
	return
}

func (dc *driverConnection) HSet(ctx context.Context, key []byte, field []string, v interface{}) (err error) {
	return
}
