// Copyright 2020 Talhuang<talhuang1231@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package options

import (
	"time"

	"github.com/skeleton1231/gotal/pkg/db"
	"github.com/spf13/pflag"
	"gorm.io/gorm"
)

// MySQLOptions defines options for mysql database.
type MySQLOptions struct {
	Host                  string        `json:"host,omitempty"                     mapstructure:"host"`
	Username              string        `json:"username,omitempty"                 mapstructure:"username"`
	Password              string        `json:"-"                                  mapstructure:"password"`
	Database              string        `json:"database"                           mapstructure:"database"`
	MaxIdleConnections    int           `json:"max-idle-connections,omitempty"     mapstructure:"max-idle-connections"`
	MaxOpenConnections    int           `json:"max-open-connections,omitempty"     mapstructure:"max-open-connections"`
	MaxConnectionLifeTime time.Duration `json:"max-connection-life-time,omitempty" mapstructure:"max-connection-life-time"`
	LogLevel              int           `json:"log-level"                          mapstructure:"log-level"`
}

// NewMySQLOptions create a `zero` value instance.
func NewMySQLOptions() *MySQLOptions {
	return &MySQLOptions{
		Host:                  "127.0.0.1:3306",
		Username:              "",
		Password:              "",
		Database:              "",
		MaxIdleConnections:    100,
		MaxOpenConnections:    100,
		MaxConnectionLifeTime: time.Duration(10) * time.Second,
		LogLevel:              1, // Silent
	}
}

// Validate verifies flags passed to MySQLOptions.
func (o *MySQLOptions) Validate() []error {
	errs := []error{}

	return errs
}

// AddFlags adds flags related to mysql storage for a specific APIServer to the specified FlagSet.
func (o *MySQLOptions) AddFlags(fs *pflag.FlagSet) {

	fs.StringVar(&o.Host, "mysql.host", o.Host, "MySQL service host address.")

	fs.StringVar(&o.Username, "mysql.username", o.Username, "Username for accessing mysql service.")

	fs.StringVar(&o.Password, "mysql.password", o.Password, "Password for accessing mysql service.")

	fs.StringVar(&o.Database, "mysql.database", o.Database, "Name of the database for the server to use.")

	fs.IntVar(&o.MaxIdleConnections, "mysql.max-idle-connections", o.MaxOpenConnections, "Max idle connections allowed for mysql.")

	fs.IntVar(&o.MaxOpenConnections, "mysql.max-open-connections", o.MaxOpenConnections, "Max open connections allowed for mysql.")

	fs.DurationVar(&o.MaxConnectionLifeTime, "mysql.max-connection-life-time", o.MaxConnectionLifeTime, "Max connection life time for mysql.")

	fs.IntVar(&o.LogLevel, "mysql.log-mode", o.LogLevel, "Specify gorm log level.")
}

// NewClient create mysql store with the given config.
func (o *MySQLOptions) NewClient() (*gorm.DB, error) {
	opts := &db.Options{
		Host:                  o.Host,
		Username:              o.Username,
		Password:              o.Password,
		Database:              o.Database,
		MaxIdleConnections:    o.MaxIdleConnections,
		MaxOpenConnections:    o.MaxOpenConnections,
		MaxConnectionLifeTime: o.MaxConnectionLifeTime,
		LogLevel:              o.LogLevel,
	}

	return db.New(opts)
}
