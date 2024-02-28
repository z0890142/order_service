package database

import (
	"database/sql"
	"fmt"

	"order_service/c"
	"order_service/config"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func OpenPostgreSQLDatabase(option *config.DatabaseOption) (db *sql.DB, err error) {
	connection, err := GetConnectionString(option)
	if err != nil {
		return nil, fmt.Errorf("openPostgreSQLDatabase: %v", err)
	}

	if db, err = sql.Open(c.DriverPostgreSQL, connection); err != nil {
		return nil, fmt.Errorf("openPostgreSQLDatabase: %v", err)
	} else {
		err = db.Ping()
		if err != nil {
			return nil, fmt.Errorf("openPostgreSQLDatabase: %v", err)
		}
	}

	// Set connection pool
	if option.PoolSize > 0 {
		db.SetMaxIdleConns(option.PoolSize)
		db.SetMaxOpenConns(option.PoolSize)
	}

	return
}

func GetConnectionString(option *config.DatabaseOption) (string, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		option.Host, option.Port, option.Username, option.Password, option.DBName)
	return connStr, nil
}

func InitGormClient(db *sql.DB) (*gorm.DB, error) {
	gormClient, err := gorm.Open(postgres.New(postgres.Config{
		PreferSimpleProtocol: true,
		Conn:                 db,
	}), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		return nil, fmt.Errorf("InitGormClientHook: %s", err.Error())
	}
	return gormClient, nil
}
