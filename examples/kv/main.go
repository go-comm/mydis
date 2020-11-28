package main

import (
	"context"
	"database/sql"
	"log"

	"github.com/go-comm/mydis"
	"github.com/go-comm/mydis/driver/mysql"
)

func main() {

	db, err := mydis.OpenDB(mysql.NewConn(sql.Open("mysql", "")))
	// db, err := mydis.Open("mysql", "")

	if err != nil {
		log.Println(err)
		return
	}

	db.Set(context.TODO(), []byte("foo"), "fff")

	db.Get(context.TODO(), []byte("foo"))

}
