package db

import (
	"database/sql"
	"log"
	"os"
	"path"
)

type DB struct {
	database *sql.DB
}

type User struct {
	Username string
	Password string
	Faction  string
	Token    string
}

func NewDB() (DB, error) {
	newDB := DB{}
	cwd, err := os.Getwd()
	if err != nil {
		return newDB, err
	}
	dbConn, err := sql.Open("sqlite3", path.Join(cwd, "spacetraders.db"))
	if err != nil {
		return newDB, err
	}
	newDB.database = dbConn
	return newDB, nil
}

func (db *DB) CreateUser(user User) error {
	insertSQL := "INSERT INTO users(username, password, faction, accessToken) VALUES(?, ?, ?)"
	statement, err := db.database.Prepare(insertSQL)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	_, err = statement.Exec(user.Username, user.Password, user.Faction, user.Token)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}
