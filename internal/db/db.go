package db

import (
	"database/sql"
	"errors"
	"log"
	"path"
	"sync"
	"time"

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
	db.mu.RLock()
	defer db.mu.RUnlock()
	var user User
	selectSQL := `SELECT * FROM users WHERE username = ?;`
	row := db.database.QueryRow(selectSQL, username)
	err := row.Scan(&user.Username, &user.Password, &user.Faction, &user.Token)
	if err != nil {
		return user, errors.New("invalid user")
	}
	return user, nil
}

func (db *DB) WriteRefreshToken(refreshToken, accessToken string, expiryTime int64) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	insertSQL := "INSERT INTO refreshTokens(refreshToken, accessToken, expiryTime) VALUES(?, ?, ?)"
	statement, err := db.database.Prepare(insertSQL)
	if err != nil {
		return err
	}
	_, err = statement.Exec(refreshToken, accessToken, expiryTime)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) ReadRefreshToken(refreshToken string) (string, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	var accessToken string
	selectSQL := `SELECT accessToken FROM refreshTokens WHERE refreshToken = ? and expiryTime > ?;`
	row := db.database.QueryRow(selectSQL, refreshToken, time.Now().Unix())
	err := row.Scan(&accessToken)
	if err != nil {
		return "", errors.New("invalid token")
	}
	return accessToken, nil
}
