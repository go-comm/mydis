package mysql

import (
	"github.com/go-comm/mydis"
)

func init() {
	mydis.Register("mysql", &MysqlDriver{})
}
