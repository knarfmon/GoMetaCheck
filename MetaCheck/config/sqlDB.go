package config

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"

	"os"
)
var SQLdb *sql.DB

func init(){
	//var (
	//	connectionName = mustGetenv("CLOUDSQL_CONNECTION_NAME")   //=====web
	//	user           = mustGetenv("CLOUDSQL_USER")
	//	password       = os.Getenv("CLOUDSQL_PASSWORD")
	//)

	var err error
	SQLdb, err = sql.Open("mysql","knarfmon:Great4me@/getmetacheck")
	//sqlDB, err = sql.Open("mysql", fmt.Sprintf("%s:%s@cloudsql(%s)/getmetacheck", user, password, connectionName)) //===web
	if err != nil {
		log.Fatalf("Could not open db: %v", err)
	}
}


func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Panicf("%s environment variable not set.", k)
	}
	return v
}
