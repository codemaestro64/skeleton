package models

import (
	"fmt"
	"time"

	"github.com/codemaestro64/skeleton/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gocraft/dbr"
)

type Model interface {
	/**First()
	Find()
	Create()
	Update(m Model)
	Delete()**/
}

type session struct {
	*dbr.Session
}

type Database struct {
	connection *dbr.Connection
	session    *session
	models     map[string]Model
}

const Models = "MODELS"

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

	// https://github.com/go-sql-driver/mysql/issues/461
	conn.SetConnMaxLifetime(time.Minute * cfg.ConnMaxLifetime)
	conn.SetMaxIdleConns(cfg.MaxIdleConns)
	conn.SetMaxOpenConns(cfg.MaxOpenConns)

	d := &Database{connection: conn}
	d.registerModels()

	return d, nil
}

func (d *Database) NewSession() *Database {
	d.session = &session{d.connection.NewSession(nil)}
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

func (d *Database) GetModel(name string) (Model, error) {
	if model, ok := d.models[name]; ok {
		return model, nil
	}

	return nil, fmt.Errorf("%s: %s is not registered", Models, name)
}

func (d *Database) registerModelOnce(name string, model Model) {
	if _, ok := d.models[name]; !ok {
		d.registerModel(name, model)
	}
}

func (d *Database) registerModel(name string, model Model) {
	if d.models == nil {
		d.models = make(map[string]Model)
	}
	d.models[name] = model
}

func (d *Database) registerModels() {
	d.registerUserModel()
}
