package models

import (
	"fmt"

	"github.com/codemaestro64/skeleton/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gocraft/dbr"
)

type Database struct {
	connection *dbr.Connection
	session    *dbr.Session
}

func Connect(cfg config.DatabaseConfig) (*Database, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=%t", cfg.Username, cfg.Password, cfg.Host, cfg.Name, cfg.ParseTime)

	conn, err := dbr.Open(cfg.Driver, dsn, nil)
	if err != nil {
		return nil, fmt.Errorf("error opening database connection: %s", err.Error())
	}

	err = conn.Ping()
	if err != nil {
		return nil, fmt.Errorf("error pinging database: %s", err.Error())
	}

	return &Database{
		connection: conn,
	}, nil
}

func (d *Database) NewSession() *Database {
	d.session = d.connection.NewSession(nil)
	return d
}

func (d *Database) Disconnect() error {
	if d.connection != nil {
		err := d.connection.Close()
		if err != nil {
			return fmt.Errorf("error closing database connection: %s", err.Error())
		}
	}

	return nil
}
