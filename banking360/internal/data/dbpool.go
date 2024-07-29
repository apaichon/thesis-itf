package data
import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/ClickHouse/clickhouse-go/v2"
)

type DBPool struct {
	db *sql.DB
}

var (
	pools     = make(map[string]*DBPool)
	poolsLock = &sync.Mutex{}
	poolsOnce = make(map[string]*sync.Once)
)

func GetDBPool(dsn string, maxOpenConns, maxIdleConns int, connMaxLifetime time.Duration) (*DBPool, error) {
	poolsLock.Lock()
	if pool, exists := pools[dsn]; exists {
		poolsLock.Unlock()
		return pool, nil
	}

	if _, exists := poolsOnce[dsn]; !exists {
		poolsOnce[dsn] = &sync.Once{}
	}
	once := poolsOnce[dsn]
	poolsLock.Unlock()

	var pool *DBPool
	var err error

	once.Do(func() {
		pool, err = newDBPool(dsn, maxOpenConns, maxIdleConns, connMaxLifetime)
		if err == nil {
			poolsLock.Lock()
			pools[dsn] = pool
			poolsLock.Unlock()
		}
	})

	return pool, err
}

func newDBPool(dsn string, maxOpenConns, maxIdleConns int, connMaxLifetime time.Duration) (*DBPool, error) {
	db, err := sql.Open("clickhouse", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(connMaxLifetime)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DBPool{db: db}, nil
}

func (p *DBPool) Close() error {
	return p.db.Close()
}

func (p *DBPool) Insert(query string, args ...interface{}) (int64, error) {
	result, err := p.db.Exec(query, args...)
	if err != nil {
		return -1,fmt.Errorf("failed to insert: %w", err)
	}
	
	return result.RowsAffected()
}

func (p *DBPool) Update(query string, args ...interface{}) error {
	_, err := p.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update: %w", err)
	}
	return nil
}


func (p *DBPool) Delete(query string, args ...interface{}) error {
	_, err := p.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete: %w", err)
	}
	return nil
}

func (p *DBPool) Query(query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := p.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}
	return rows, nil
}