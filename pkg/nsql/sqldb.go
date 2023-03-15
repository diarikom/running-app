package nsql

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

type SqlDatabase struct {
	Conn *sqlx.DB
}

// Prepare prepare sql statements or exit api if fails or error
func (s *SqlDatabase) Prepare(query string) *sqlx.Stmt {
	stmt, err := s.Conn.Preparex(query)
	if err != nil {
		panic(fmt.Errorf("nsql: error while preparing statment [%s] (%s)", query, err))
	}
	return stmt
}

// PrepareNamed prepare sql statements with named bindvars or exit api if fails or error
func (s *SqlDatabase) PrepareNamed(query string) *sqlx.NamedStmt {
	stmt, err := s.Conn.PrepareNamed(query)
	if err != nil {
		panic(fmt.Errorf("nsql: error while preparing statment [%s] (%s)", query, err))
	}
	return stmt
}

func NewSqlDatabase(conf Config) (*SqlDatabase, error) {
	// Set default connection values
	conf.setDefault()

	// Get DSN
	dsn, err := conf.getDSN()
	if err != nil {
		return nil, err
	}

	// Create connection
	db, err := sqlx.Connect(conf.Driver, dsn)
	if err != nil {
		return nil, err
	}

	// Set connection settings
	db.SetConnMaxLifetime(time.Duration(*conf.MaxConnLifetime) * time.Second)
	db.SetMaxOpenConns(*conf.MaxOpenConn)
	db.SetMaxIdleConns(*conf.MaxIdleConn)

	// Return
	return &SqlDatabase{Conn: db}, nil
}
