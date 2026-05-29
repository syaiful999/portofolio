package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/lib/pq"
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

type PoolConfig struct {
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

func DefaultPoolConfig() PoolConfig {
	return PoolConfig{
		MaxOpenConns:    100,
		MaxIdleConns:    10,
		ConnMaxLifetime: 30 * time.Minute,
		ConnMaxIdleTime: 5 * time.Minute,
	}
}

type Postgre struct {
	DB
	SslMode string
	Tz      string
}

func (m Postgre) Connect() (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		m.DB.Address, m.DB.Port, m.DB.User, m.DB.Password, m.DB.Name,
	)
	return sql.Open("postgres", dsn)
}

func InitPostgre(db *DB) *Postgre {
	return &Postgre{DB: *db}
}

type SqlServer struct {
	DB
}

func (m SqlServer) Connect() (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"server=%s;user id=%s;password=%s;port=%d;database=%s;",
		m.DB.Host, m.DB.User, m.DB.Password, m.DB.Port, m.DB.Name,
	)
	return sql.Open("sqlserver", dsn)
}

func InitSqlServer(db *DB) *SqlServer {
	return &SqlServer{DB: *db}
}

func Connection(db *DB) (*sql.DB, error) {
	return ConnectionWithPool(db, DefaultPoolConfig())
}

func ConnectionWithPool(db *DB, pool PoolConfig) (*sql.DB, error) {
	var conn IConnect

	switch strings.ToLower(db.Driver) {
	case "postgre", "postgres", "postgresql":
		conn = InitPostgre(db)
	case "sqlserver", "mssql":
		conn = InitSqlServer(db)
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", db.Driver)
	}

	dbConn, err := conn.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to open connection: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = dbConn.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	dbConn.SetMaxOpenConns(pool.MaxOpenConns)
	dbConn.SetMaxIdleConns(pool.MaxIdleConns)
	dbConn.SetConnMaxLifetime(pool.ConnMaxLifetime)
	dbConn.SetConnMaxIdleTime(pool.ConnMaxIdleTime)

	return dbConn, nil
}
