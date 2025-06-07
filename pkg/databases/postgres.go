package databases

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/lib/pq"
	"github.com/saufiroja/fin-ai/config"
	logging "github.com/saufiroja/fin-ai/pkg/loggings"
)

type PostgresManager interface {
	Connection() *sql.DB
	StartTransaction() (*sql.Tx, error)
	CommitTransaction(tx *sql.Tx) error
	RollbackTransaction(tx *sql.Tx) error
	CloseConnection() error
}

type Postgres struct {
	db *sql.DB
}

var (
	postgresInstance *Postgres
	once             sync.Once
)

func NewPostgres(conf *config.AppConfig, logger logging.Logger) PostgresManager {
	once.Do(func() {
		user := conf.Postgres.User
		password := conf.Postgres.Pass
		dbHost := conf.Postgres.Host
		dbPort := conf.Postgres.Port
		dbName := conf.Postgres.Name
		dbSslMode := conf.Postgres.SSL

		dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			dbHost, dbPort, user, password, dbName, dbSslMode)

		db, err := sql.Open("postgres", dsn)
		if err != nil {
			logger.LogPanic(fmt.Sprintf("Error opening databases: %v", err))
		}

		// Set connection pool settings
		db.SetMaxOpenConns(20)                 // Maksimum total koneksi terbuka ke DB
		db.SetMaxIdleConns(10)                 // Maksimum koneksi idle di pool
		db.SetConnMaxIdleTime(5 * time.Minute) // Idle lebih dari ini akan ditutup
		db.SetConnMaxLifetime(1 * time.Hour)   // Umur maksimal koneksi

		if err := db.Ping(); err != nil {
			logger.LogPanic(fmt.Sprintf("Error connecting to databases: %v", err))
		}

		logger.LogInfo("Database connected")
		postgresInstance = &Postgres{db: db}
	})

	return postgresInstance
}

func (p *Postgres) Connection() *sql.DB {
	return p.db
}

func (p *Postgres) StartTransaction() (*sql.Tx, error) {
	return p.db.Begin()
}

func (p *Postgres) CommitTransaction(tx *sql.Tx) error {
	return tx.Commit()
}

func (p *Postgres) RollbackTransaction(tx *sql.Tx) error {
	return tx.Rollback()
}

func (p *Postgres) CloseConnection() error {
	return p.db.Close()
}
