package config

import (
	"database/sql"
	"fmt"
	"log"
)

func InitDb(cfg *Config) (*sql.DB, error) {
	param := "parseTime=true"
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?%s", cfg.Db.User, cfg.Db.Password, cfg.Db.Host, cfg.Db.DbSchema, param)
	db, err := sql.Open("mysql", dsn)
	log.Println("database info>>>", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mysql:%s", err.Error())
	}
	return db, nil
}
