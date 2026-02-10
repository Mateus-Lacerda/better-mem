package sqlite

import (
	"github.com/Mateus-Lacerda/better-mem/internal/config"
	"database/sql"
	"fmt"
	"log"

	sqlite_vec "github.com/asg017/sqlite-vec-go-bindings/cgo"
	"github.com/khepin/liteq"
	_ "github.com/mattn/go-sqlite3"
	gorm_sqlite "gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDb() *sql.DB {
	// Setup sqlite_vec for vector stuff
	sqlite_vec.Auto()

	// Setup liteq for local queue stuff
	setupLiteQ()

	db, err := sql.Open("sqlite3", config.SQLite.SQLiteDBLocation)
	if err != nil {
		log.Fatal(err)
	}
	var sqliteVersion string
	var vecVersion string
	err = db.QueryRow("select sqlite_version(), vec_version()").Scan(&sqliteVersion, &vecVersion)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("sqlite_version=%s, vec_version=%s\n", sqliteVersion, vecVersion)

	return db
}

// Bootstraps liteq's tables if needed
// TODO: Catch errors except for table already exists
func setupLiteQ() {
	liteqDb, err := sql.Open("sqlite3", config.SQLite.SQLiteDBLocation)
	if err != nil {
		// log.Fatal(err)
	}
	if err := liteq.Setup(liteqDb); err != nil {
		// log.Fatal(err)
	}
	liteqDb.Close()
}

// Returns the SQLite connection
func GetDb() *gorm.DB {
	db, err := gorm.Open(
		gorm_sqlite.Dialector{
			DriverName: "sqlite3",
			DSN:        config.SQLite.SQLiteDBLocation,
		},
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Info), // Logging verboso

		},
	)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

// Returns the sqlite connection with sqlite_vec applied
func GetDbVec() *sql.DB {
	sqlite_vec.Auto()

	db, err := sql.Open("sqlite3", config.SQLite.SQLiteDBLocation)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
