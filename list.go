package mydis

import (
	"context"
)

func (db *DB) LPush(ctx context.Context, k []byte, v0 interface{}, v ...interface{}) error {
	return Update(ctx, db, func(tx Tx) error {
		return tx.LPush(ctx, k, v0, v...)
	})
}

func (db *DB) RPush(ctx context.Context, k []byte, v0 interface{}, v ...interface{}) error {
	return Update(ctx, db, func(tx Tx) error {
		return tx.RPush(ctx, k, v0, v...)
	})
}

func (db *DB) LPop(ctx context.Context, k []byte) (v interface{}, err error) {
	err = Update(ctx, db, func(tx Tx) error {
		v, err = tx.LPop(ctx, k)
		return err
	})
	return v, err
}

func (db *DB) RPop(ctx context.Context, k []byte) (v interface{}, err error) {
	err = Update(ctx, db, func(tx Tx) error {
		v, err = tx.RPop(ctx, k)
		return err
	})
	return v, err
}
