package db

import (
	"database/sql"
	"log"
	"path"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	database *sql.DB
	mu       sync.RWMutex
}

type User struct {
	Username string
	Password []byte
	Faction  string
	Token    string
}

func NewDB(DB_PATH string) (*DB, error) {
	newDB := DB{}
	dbConn, err := sql.Open("sqlite3", path.Join(DB_PATH, "spacetraders.db"))
	if err != nil {
		return &newDB, err
	}
	newDB.database = dbConn
	return &newDB, nil
}

func (db *DB) CreateUser(user User) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	insertSQL := "INSERT INTO users(username, password, faction, accessToken) VALUES(?, ?, ?, ?)"
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

func (db *DB) ReadUserFromDB(username string) (User, error) {
	var user User
	selectSQL := `SELECT * FROM users WHERE username = ?;`
	row := db.database.QueryRow(selectSQL, username)
	row.Scan(&user.Username, &user.Password, &user.Faction, &user.Token)
	return user, nil
}
