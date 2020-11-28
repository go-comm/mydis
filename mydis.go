package mydis

import "context"

type KV interface {
	Get(ctx context.Context, k []byte) (interface{}, error)
	GetSet(ctx context.Context, k []byte, v interface{}) (interface{}, error)
	Set(ctx context.Context, k []byte, v interface{}) error
	SetEx(ctx context.Context, k []byte, v interface{}, sec int64) error
	SetNx(ctx context.Context, k []byte, v interface{}) error
	Delete(ctx context.Context, k []byte) error
	Expire(ctx context.Context, k []byte, sec int64) error
	PExpire(ctx context.Context, k []byte, mill int64) error
}
