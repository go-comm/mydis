package mydis

import (
	"errors"
)

var (
	ErrNoKey = errors.New("mydis: key not found")

	ErrMaxSizeExceed = errors.New("mydis: max size exceed")
)
