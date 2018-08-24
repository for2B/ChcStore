package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func NewDB(host, port, user, pwd, dbName, driverName string) *sql.DB {
	//构建连接字符串
	dateSource := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pwd, dbName)
	fmt.Println(fmt.Sprintf("current db detail:host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host,
		port, user, pwd,	dbName))
	db, err := sql.Open(driverName, dateSource)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		log.Println("init model failed at ping: " + err.Error())
		panic(err)
	}
	return db
}
