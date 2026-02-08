package config

import "os"

type sqliteConfig struct {
	SQLiteDBLocation string
	SqliteVecPath          string
}

func newSqliteConfig() sqliteConfig {
	betterMemDataPath := getString(
		"XDG_DATA_HOME", getString(
			"HOME",
			"",
		)+"/better-mem",
	)
	if err := os.Mkdir(betterMemDataPath, os.ModeDir); err != nil {
	}
	sqliteDBLocation := betterMemDataPath + "/better-mem.db"
	vssPath := betterMemDataPath + "/sqlite_extensions/vec.so"
	return sqliteConfig{
		SQLiteDBLocation: sqliteDBLocation,
		SqliteVecPath:          vssPath,
	}
}

var SQLite = newSqliteConfig()
