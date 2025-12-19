package clickhouse

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/ClickHouse/clickhouse-go/v2"
	_ "github.com/ClickHouse/clickhouse-go/v2/lib/driver"

	"github.com/SoulStalker/data_bridge_ch/internal/config"
)

func InitDB(cfg config.CHConfig) error {
	db, err := sql.Open("clickhouse", cfg.DSN)
	if err != nil {
		log.Fatalf("open clickhouse: %v", err)
	}

	// TODO: remove test conn
	var version string
	if err := db.QueryRow("SELECT version()").Scan(&version); err != nil {
		log.Fatalf("query version: %v", err)
	}
	fmt.Println("CH version:", version)

	// TODO: creating goods table
	return ensureDimGoodsTable(db)
	// return db
}

func ensureDimGoodsTable(db *sql.DB) error {
	ddl := `
		CREATE TABLE IF NOT EXISTS dim_goods
		(
			id        UUID,
			parent_id UUID,
			code      String,
			description String
		) ENGINE = MergeTree()
		ORDER BY id;
		`
	_, err := db.Exec(ddl)
	if err != nil {
		return fmt.Errorf("create goods: %w", err)
	}
	return nil
}
