package db

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"

	"github.com/jmoiron/sqlx"
)

func Migrate(db *sqlx.DB) {
	path := "internal/db/migrations"

	files, err := filepath.Glob(filepath.Join(path, "*.sql"))
	if err != nil {
		log.Fatalf("failed to read migrations: %v", err)
	}

	//文件按名称排序
	sort.Strings(files)

	for _, f := range files {
		sqlBytes, err := ioutil.ReadFile(f)
		if err != nil {
			log.Fatalf("failed to read migration file %s: %v", f, err)
		}

		sql := string(sqlBytes)
		fmt.Printf("Applying migration: %s\n", f)

		if _, err := db.Exec(sql); err != nil {
			log.Fatalf("failed to execute migration %s: %v", f, err)
		}
	}

	fmt.Println("All migrations applied.")
}
