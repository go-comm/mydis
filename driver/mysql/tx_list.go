package mysql

import (
	"context"
	"strconv"

	"github.com/go-comm/mydis"
)

type ListMeta struct {
	Size int `json:"size"`
	L    int `json:"l"`
	R    int `json:"r"`
}

func (tx *mysqlTx) makeListMetaKey(k []byte) []byte {
	var d = make([]byte, 0, len(k)+5)
	d = append(d, k...)
	d = append(d, []byte("#LIST")...)
	return d
}

func (tx *mysqlTx) makeListIdxKey(k []byte, idx int) []byte {
	var d = make([]byte, 0, len(k)+6)
	d = append(d, k...)
	d = append(d, '@')
	strconv.AppendInt(d, int64(idx), 0)
	return d
}

func (tx *mysqlTx) LPush(ctx context.Context, k []byte, v interface{}) error {
	var metaKey = tx.makeListMetaKey(k)

	var meta ListMeta
	err := tx.Scan(ctx, &meta, metaKey)
	if err != nil && err != mydis.ErrNoKey {
		return err
	}

	var idxKey = tx.makeListIdxKey(k, meta.L)
	if err = tx.Set(ctx, idxKey, v); err != nil {
		return err
	}
	meta.L--
	return tx.Set(ctx, metaKey, meta)
}

func (tx *mysqlTx) RPush(ctx context.Context, k []byte, v interface{}) error {
	var metaKey = tx.makeListMetaKey(k)

	var meta ListMeta
	err := tx.Scan(ctx, &meta, metaKey)
	if err != nil && err != mydis.ErrNoKey {
		return err
	}

	var idxKey = tx.makeListIdxKey(k, meta.R)
	if err = tx.Set(ctx, idxKey, v); err != nil {
		return err
	}
	meta.R++
	return tx.Set(ctx, metaKey, meta)
}

func (tx *mysqlTx) LPop(ctx context.Context, k []byte) (interface{}, error) {
	var metaKey = tx.makeListMetaKey(k)

	var meta ListMeta
	err := tx.Scan(ctx, &meta, metaKey)
	if err != nil {
		return nil, err
	}
	meta.L++
	var idxKey = tx.makeListIdxKey(k, meta.L)
	v, err := tx.Get(ctx, idxKey)
	if err != nil {
		return nil, err
	}
	if err = tx.Del(ctx, idxKey); err != nil {
		return nil, err
	}
	if meta.L >= meta.R {
		err = tx.Del(ctx, metaKey)
	} else {
		err = tx.Set(ctx, metaKey, meta)
	}
	return v, err
}

func (tx *mysqlTx) RPop(ctx context.Context, k []byte) (interface{}, error) {
	var metaKey = tx.makeListMetaKey(k)

	var meta ListMeta
	err := tx.Scan(ctx, &meta, metaKey)
	if err != nil {
		return nil, err
	}
	meta.R--
	var idxKey = tx.makeListIdxKey(k, meta.R)
	v, err := tx.Get(ctx, idxKey)
	if err != nil {
		return nil, err
	}
	if err = tx.Del(ctx, idxKey); err != nil {
		return nil, err
	}
	if meta.L >= meta.R {
		err = tx.Del(ctx, metaKey)
	} else {
		err = tx.Set(ctx, metaKey, meta)
	}
	return v, err
}
