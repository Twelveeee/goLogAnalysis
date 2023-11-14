package entity

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"gopkg.in/yaml.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var initDbOnce sync.Once

type testDb struct {
	db       *gorm.DB
	Database struct {
		Driver    string `yaml:"Driver"`
		Dsn       string `yaml:"Dsn"`
		DbName    string `yaml:"DbName"`
		Server    string `yaml:"Server"`
		User      string `yaml:"User"`
		Password  string `yaml:"Password"`
		Port      int    `yaml:"Port"`
		ConnsIdle int    `yaml:"ConnsIdle"`
		Conns     int    `yaml:"Conns"`
	} `yaml:"Database"`
}

func (c *testDb) Db() *gorm.DB {
	if c.db == nil {
		log.Fatal().Msg("config: database not connected")
	}

	return c.db
}

func initDb() {
	config := &testDb{}
	yamlConfig, err := os.ReadFile("../../config/config.yaml")
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(yamlConfig, &config)
	if err != nil {
		panic(err)
	}

	dbDsn := config.Database.Dsn

	if dbDsn == "" {
		address := config.Database.Server

		// Connect via TCP or Unix Domain Socket?
		if strings.HasPrefix(address, "/") {
			log.Debug().Msg("mariadb: connecting via Unix domain socket")
			address = fmt.Sprintf("unix(%s)", address)
		} else {
			address = fmt.Sprintf("tcp(%s)", address)
		}

		dbDsn = fmt.Sprintf(
			"%s:%s@%s/%s?charset=utf8mb4,utf8&collation=utf8mb4_unicode_ci&parseTime=true",
			config.Database.User,
			config.Database.Password,
			address,
			config.Database.DbName,
		)
	}

	db, err := gorm.Open(mysql.Open(dbDsn))
	if err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	sqlDB.SetMaxOpenConns(config.Database.Conns)
	sqlDB.SetMaxIdleConns(config.Database.ConnsIdle)
	sqlDB.SetConnMaxLifetime(time.Hour)

	db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci")

	config.db = db
	SetDbProvider(config)

}

func testInitDB() {
	initDbOnce.Do(initDb)
}
