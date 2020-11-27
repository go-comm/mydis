package mydis

import "context"

func (dc *driverConnection) Get(ctx context.Context, key []byte) (v interface{}, err error) {
	return
}

func (dc *driverConnection) GetSet(ctx context.Context, key []byte, v interface{}) (oldv interface{}, err error) {
	return
}

func (dc *driverConnection) Set(ctx context.Context, key []byte, v interface{}) (err error) {
	return
}

func (dc *driverConnection) SetEx(ctx context.Context, key []byte, v interface{}, sec int64) (err error) {
	return
}

func (dc *driverConnection) SetNx(ctx context.Context, key []byte, v interface{}) (err error) {
	return
}

func (dc *driverConnection) Delete(ctx context.Context, key []byte) (err error) {
	return
}

func (dc *driverConnection) Expire(ctx context.Context, key []byte, sec int64) (err error) {
	return
}

func (dc *driverConnection) PExpire(ctx context.Context, key []byte, mill int64) (err error) {
	return
}
