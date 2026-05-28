package database

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
)

var (
	dbConn *sql.DB
	err    error

	Conn IConnect
)

type IConnect interface {
	Connect() (*sql.DB, error)
}

type DB struct {
	Name     string
	User     string
	Password string
	Host     string
	Port     int
	Driver   string
	Address  string
}

type Postgre struct {
	DB
	SslMode string
	Tz      string
}

func (m Postgre) Connect() (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", m.DB.Address, m.DB.Port, m.DB.User, m.DB.Password, m.DB.Name)
	return sql.Open("postgres", dsn)
}

func InitPostgre(db *DB) *Postgre {
	init := &Postgre{
		DB: *db,
	}
	return init
}

type SqlServer struct {
	DB
}

func (m SqlServer) Connect() (*sql.DB, error) {
	dsn := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;", m.DB.Host, m.DB.User, m.DB.Password, m.DB.Port, m.DB.Name)
	return sql.Open("sqlserver", dsn)
}

func InitSqlServer(db *DB) *SqlServer {
	init := &SqlServer{
		DB: *db,
	}
	return init
}

func Connection(db *DB) (*sql.DB, error) {
	switch s := db.Driver; {
	case strings.ToLower(s) == "postgre":
		Conn = InitPostgre(db)
	case strings.ToLower(s) == "sqlserver":
		Conn = InitSqlServer(db)
	}
	dbConn, err = Conn.Connect()

	if err != nil {
		return nil, err
	}

	if err = dbConn.Ping(); err != nil {
		return nil, err
	}

	dbConn.SetMaxIdleConns(10)
	dbConn.SetMaxOpenConns(500)

	return dbConn, nil
}
