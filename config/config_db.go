package config

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/twelveeee/log_analysis/service/entity"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	MySQL    = "mysql"
	MariaDB  = "mariadb"
	Postgres = "postgres"
	SQLite3  = "sqlite3"
)

var lockDb = sync.Mutex{}

func (c *Config) initDb() error {
	lockDb.Lock()
	defer lockDb.Unlock()

	dbDriver := c.DatabaseDriver()
	dbDsn := c.DatabaseDsn()

	if dbDriver == "" {
		return errors.New("config: database driver not specified")
	}

	if dbDsn == "" {
		return errors.New("config: database DSN not specified")
	}

	// Open database connection.
	var db *gorm.DB

	if dbDriver == "mysql" {
		mysqlDb, err := getMysqlDb(dbDsn)
		if err != nil {
			return err
		}

		db = mysqlDb
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	// Set database connection parameters.
	sqlDB.SetMaxOpenConns(c.DatabaseConns())
	sqlDB.SetMaxIdleConns(c.DatabaseConnsIdle())
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Check database server version.
	if err = c.checkDb(db); err != nil {
		return err
	}

	// Ok.
	c.db = db

	return nil
}

func (c *Config) checkDb(db *gorm.DB) error {
	switch c.DatabaseDriver() {
	case MySQL:
		type Res struct {
			Value string `gorm:"column:Value;"`
		}
		var res Res
		if err := db.Raw("SHOW VARIABLES LIKE 'innodb_version'").Scan(&res).Error; err != nil {
			return err
		}
	}

	return nil
}

func (c *Config) DatabaseDriver() string {
	switch strings.ToLower(c.options.Database.Driver) {
	case MySQL, MariaDB:
		c.options.Database.Driver = MySQL
	}

	return c.options.Database.Driver
}

func (c *Config) DatabaseDsn() string {
	if c.options.Database.Dsn == "" {
		switch c.DatabaseDriver() {
		case MySQL, MariaDB:
			address := c.DatabaseServer()

			// Connect via TCP or Unix Domain Socket?
			if strings.HasPrefix(address, "/") {
				log.Debug().Msg("mariadb: connecting via Unix domain socket")
				address = fmt.Sprintf("unix(%s)", address)
			} else {
				address = fmt.Sprintf("tcp(%s)", address)
			}

			return fmt.Sprintf(
				"%s:%s@%s/%s?charset=utf8mb4,utf8&collation=utf8mb4_unicode_ci&parseTime=true",
				c.DatabaseUser(),
				c.DatabasePassword(),
				address,
				c.DatabaseDbName(),
			)
		default:
			log.Error().Msg("config: empty database dsn")
			return ""
		}
	}

	return c.options.Database.Dsn
}

func (c *Config) DatabaseServer() string {
	return c.options.Database.Server
}

func (c *Config) DatabaseUser() string {
	return c.options.Database.User
}

func (c *Config) DatabasePassword() string {
	return c.options.Database.Password
}

func (c *Config) DatabaseDbName() string {
	return c.options.Database.DbName
}

func (c *Config) DatabasePort() int {
	return c.options.Database.Port
}

func (c *Config) DatabaseConns() int {
	return c.options.Database.Conns
}

func (c *Config) DatabaseConnsIdle() int {
	return c.options.Database.ConnsIdle
}

func (c *Config) Db() *gorm.DB {
	if c.db == nil {
		log.Fatal().Msg("config: database not connected")
	}

	return c.db
}

// RegisterDb sets the database options and connection provider.
func (c *Config) RegisterDb() {
	c.SetDbOptions()
	entity.SetDbProvider(c)
	entity.ClearHasTableCache()
}

// SetDbOptions sets the database collation to unicode if supported.
func (c *Config) SetDbOptions() {
	switch c.DatabaseDriver() {
	case MySQL, MariaDB:
		c.Db().Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci")
	case Postgres:
		// Ignore for now.
	case SQLite3:
		// Not required as unicode is default.
	}
}

func getMysqlDb(dbDsn string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dbDsn))
	if err != nil || db == nil {
		for i := 1; i <= 3; i++ {
			db, err = gorm.Open(mysql.Open(dbDsn))
			if db != nil && err == nil {
				break
			}

			time.Sleep(5 * time.Second)
		}

		if err != nil || db == nil {
			return nil, errors.New("config: db connect fail")
		}
	}

	return db, nil
}
