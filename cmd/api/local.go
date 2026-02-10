//go:build local

package main

import "github.com/Mateus-Lacerda/better-mem/internal/database/sqlite"

func setup() error {
	sqlite.InitDb()
	sqlite.Migrate(sqlite.GetDb())
	return nil
}
