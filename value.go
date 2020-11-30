package mydis

import (
	"fmt"
	"strconv"
)

const (
	TNone uint8 = iota
	TString
	TList
	TSet
	TZSet
	THash
)

func Int64(x interface{}, err error) (int64, error) {
	switch d := x.(type) {
	case []byte:
		return strconv.ParseInt(string(d), 10, 0)
	case string:
		return strconv.ParseInt(d, 10, 0)
	}
	return 0, fmt.Errorf("mydis: couldn't convert %v (%T) into type int64", x, x)
}

func Uint64(x interface{}, err error) (uint64, error) {
	switch d := x.(type) {
	case []byte:
		return strconv.ParseUint(string(d), 10, 0)
	case string:
		return strconv.ParseUint(d, 10, 0)
	}
	return 0, fmt.Errorf("mydis: couldn't convert %v (%T) into type uint64", x, x)
}

func String(x interface{}, err error) (string, error) {
	switch d := x.(type) {
	case []byte:
		return BytesToString(d), nil
	case string:
		return d, nil
	}
	return "", fmt.Errorf("mydis: couldn't convert %v (%T) into type string", x, x)
}

func Bytes(x interface{}, err error) ([]byte, error) {
	switch d := x.(type) {
	case []byte:
		return d, nil
	case string:
		return []byte(d), nil
	}
	return nil, fmt.Errorf("mydis: couldn't convert %v (%T) into type string", x, x)
}

type Value interface{}

func Convert(dst interface{}, src interface{}) {

}
