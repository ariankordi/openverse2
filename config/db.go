package config
import (
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
  "github.com/go-gorp/gorp"

  "log"
)

func initDb() *gorp.DbMap {
	db, err := sql.Open("mysql", "USERNAME:PASSWORD@/DATABASE?collation=utf8mb4_general_ci")
  if err != nil {
    log.Fatal(err)
  }
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
	return dbmap
}

var Dbmap = initDb()
