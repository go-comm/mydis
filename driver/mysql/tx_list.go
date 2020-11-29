package mysql

import (
	"bytes"
	"context"
	"database/sql/driver"
	"fmt"
	"strconv"

	"github.com/go-comm/mydis"
)

const (
	defaultMaxListSize = 1 << 30
)

type ListMeta struct {
	Size  int `json:"s"`
	Left  int `json:"l"`
	Right int `json:"r"`
}

func (m ListMeta) marshal() ([]byte, error) {
	// avoid use  fmt.Sprintf("%d,%d,%d",...) or json.Marshal
	var buf = make([]byte, 0, 20)
	buf = strconv.AppendInt(buf, int64(m.Size), 0)
	buf = append(buf, ',')
	buf = strconv.AppendInt(buf, int64(m.Left), 0)
	buf = append(buf, ',')
	buf = strconv.AppendInt(buf, int64(m.Right), 0)
	return buf, nil
}

func (m *ListMeta) unmarshal(b []byte) error {
	bn := bytes.Split(b, []byte(","))
	if len(bn) != 3 {
		return nil
	}
	var err error
	var tmp int64
	if tmp, err = strconv.ParseInt(string(bn[0]), 10, 0); err != nil {
		return err
	}
	m.Size = int(tmp)
	if tmp, err = strconv.ParseInt(string(bn[1]), 10, 0); err != nil {
		return err
	}
	m.Left = int(tmp)
	if tmp, err = strconv.ParseInt(string(bn[2]), 10, 0); err != nil {
		return err
	}
	m.Right = int(tmp)
	return err
}

func (m ListMeta) Value() (driver.Value, error) {
	b, err := m.marshal()
	return b, err
}

func (m *ListMeta) Scan(v interface{}) error {
	switch b := v.(type) {
	case []byte:
		return m.unmarshal(b)
	case string:
		return m.unmarshal([]byte(b))
	case nil:
		return nil
	default:
		return fmt.Errorf("mydis: scanning unsupported type: %T", b)
	}
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

func (tx *mysqlTx) LPush(ctx context.Context, k []byte, v0 interface{}, v ...interface{}) error {
	var metaKey = tx.makeListMetaKey(k)
	var meta ListMeta
	err := tx.Scan(ctx, &meta, metaKey)
	if err != nil && err != mydis.ErrNoKey {
		return err
	}
	if meta.Size+1+len(v) > defaultMaxListSize {
		return mydis.ErrMaxSizeExceed
	}
	var idxKey = tx.makeListIdxKey(k, meta.Left)
	if err = tx.Set(ctx, idxKey, v0); err != nil {
		return err
	}
	meta.Left--
	meta.Size++
	for _, vv := range v {
		idxKey = tx.makeListIdxKey(k, meta.Left)
		if err = tx.Set(ctx, idxKey, vv); err != nil {
			return err
		}
		meta.Left--
		meta.Size++
	}
	return tx.Set(ctx, metaKey, meta)
}

func (tx *mysqlTx) RPush(ctx context.Context, k []byte, v0 interface{}, v ...interface{}) error {
	var metaKey = tx.makeListMetaKey(k)
	var meta ListMeta
	err := tx.Scan(ctx, &meta, metaKey)
	if err != nil && err != mydis.ErrNoKey {
		return err
	}
	if meta.Size+1+len(v) > defaultMaxListSize {
		return mydis.ErrMaxSizeExceed
	}
	var idxKey = tx.makeListIdxKey(k, meta.Right)
	if err = tx.Set(ctx, idxKey, v); err != nil {
		return err
	}
	meta.Right++
	meta.Size++
	for _, vv := range v {
		idxKey = tx.makeListIdxKey(k, meta.Right)
		if err = tx.Set(ctx, idxKey, vv); err != nil {
			return err
		}
		meta.Right++
		meta.Size++
	}
	return tx.Set(ctx, metaKey, meta)
}

func (tx *mysqlTx) LPop(ctx context.Context, k []byte) (interface{}, error) {
	var metaKey = tx.makeListMetaKey(k)
	var meta ListMeta
	err := tx.Scan(ctx, &meta, metaKey)
	if err != nil {
		return nil, err
	}
	meta.Left++
	meta.Size--
	var idxKey = tx.makeListIdxKey(k, meta.Left)
	v, err := tx.Get(ctx, idxKey)
	if err != nil {
		return nil, err
	}
	if err = tx.Del(ctx, idxKey); err != nil {
		return nil, err
	}
	if meta.Left >= meta.Right {
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
	meta.Right--
	meta.Size--
	var idxKey = tx.makeListIdxKey(k, meta.Right)
	v, err := tx.Get(ctx, idxKey)
	if err != nil {
		return nil, err
	}
	if err = tx.Del(ctx, idxKey); err != nil {
		return nil, err
	}
	if meta.Left >= meta.Right {
		err = tx.Del(ctx, metaKey)
	} else {
		err = tx.Set(ctx, metaKey, meta)
	}
	return v, err
}

func (tx *mysqlTx) LLen(ctx context.Context, k []byte) (interface{}, error) {
	var metaKey = tx.makeListMetaKey(k)
	var meta ListMeta
	err := tx.Scan(ctx, &meta, metaKey)
	if err != nil {
		return nil, err
	}
	return meta.Size, nil
}

func (tx *mysqlTx) LIndex(ctx context.Context, k []byte, i int) (interface{}, error) {
	var metaKey = tx.makeListMetaKey(k)
	var meta ListMeta
	err := tx.Scan(ctx, &meta, metaKey)
	if err != nil {
		return nil, err
	}
	idx := meta.Left + i
	var idxKey = tx.makeListIdxKey(k, idx)
	return tx.Get(ctx, idxKey)
}

func (tx *mysqlTx) LSet(ctx context.Context, k []byte, i int, v interface{}) error {
	var metaKey = tx.makeListMetaKey(k)
	var meta ListMeta
	err := tx.Scan(ctx, &meta, metaKey)
	if err != nil {
		return err
	}
	idx := meta.Left + i
	var idxKey = tx.makeListIdxKey(k, idx)
	return tx.Set(ctx, idxKey, v)
}

func (tx *mysqlTx) LRange(ctx context.Context, k []byte, start int, stop int) (interface{}, error) {
	var metaKey = tx.makeListMetaKey(k)
	var meta ListMeta
	err := tx.Scan(ctx, &meta, metaKey)
	if err != nil {
		return nil, err
	}
	var ls []interface{}
	for i := start; i < stop; i++ {
		idx := meta.Left + i
		var idxKey = tx.makeListIdxKey(k, idx)
		v, err := tx.Get(ctx, idxKey)
		if err != nil && err != mydis.ErrNoKey {
			return nil, err
		}
		ls = append(ls, v)
	}
	return ls, nil
}

func (tx *mysqlTx) LTrim(ctx context.Context, k []byte, start int, stop int) error {
	var metaKey = tx.makeListMetaKey(k)
	var meta ListMeta
	err := tx.Scan(ctx, &meta, metaKey)
	if err != nil {
		return err
	}
	// TODO
	return nil
}
