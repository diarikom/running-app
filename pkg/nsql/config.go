package nsql

import (
	"fmt"
)

type Config struct {
	Driver          string `yaml:"driver"`
	Host            string `yaml:"host"`
	Port            string `yaml:"port"`
	Username        string `yaml:"username"`
	Password        string `yaml:"password"`
	Database        string `yaml:"database"`
	MaxIdleConn     *int   `yaml:"max_idle_connection"`
	MaxOpenConn     *int   `yaml:"max_open_connection"`
	MaxConnLifetime *int   `yaml:"max_connection_lifetime"`
}

func (c *Config) setDefault() {
	// Check for optional values, set values if unset
	// If max idle connection is unset, set to 10
	if c.MaxIdleConn == nil {
		c.MaxIdleConn = newInt(10)
	}
	// If max open connection is unset, set to 10
	if c.MaxOpenConn == nil {
		c.MaxOpenConn = newInt(10)
	}
	// If max idle connection is unset, set to 1 second
	if c.MaxConnLifetime == nil {
		c.MaxConnLifetime = newInt(1)
	}
}

func (c *Config) getDSN() (dsn string, err error) {
	switch c.Driver {
	case DriverMySQL:
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", c.Username, c.Password, c.Host, c.Port,
			c.Database)
	case DriverPostgreSQL:
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", c.Host, c.Port,
			c.Username, c.Password, c.Database)
	default:
		err = fmt.Errorf("nsql: unsupported database driver %s", c.Driver)
	}
	return
}
