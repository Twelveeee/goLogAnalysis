package entity

import (
	"gorm.io/gorm"
)

// dbConn is the global gorm.DB connection provider.
var dbConn Gorm
var dbTableNameCache DbTableNameCache

// Gorm is a gorm.DB connection provider interface.
type Gorm interface {
	Db() *gorm.DB
}

type DbConn struct {
	// once sync.Once
	db *gorm.DB
}

// Set UTC as the default for created and updated timestamps.
// func init() {
// 	// gorm.NowFunc = func() time.Time {
// 	// 	return UTC()
// 	// }
// }

// Db returns the default *gorm.DB connection.
func Db() *gorm.DB {
	if dbConn == nil {
		return nil
	}

	return dbConn.Db()
}

func (g *DbConn) Db() *gorm.DB {
	if g.db == nil {
		log.Fatal().Msg("migrate: database not connected")
	}

	return g.db
}

// SetDbProvider sets the Gorm database connection provider.
func SetDbProvider(conn Gorm) {
	dbConn = conn
}
