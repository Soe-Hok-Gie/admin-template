package fondation

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"

)

func NewDB(config DatabaseConfig) *sql.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Name,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	// if err := db.Ping(); err != nil {
	// 	panic(err)
	// }

	return db
}

