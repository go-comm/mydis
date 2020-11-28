package mydis

import (
	"context"
)

type Conn interface {
	Begin(ctx context.Context, readOnly bool) (Tx, error)
	Get(ctx context.Context, k []byte) (interface{}, error)
	Set(ctx context.Context, k []byte, v interface{}) error
}
