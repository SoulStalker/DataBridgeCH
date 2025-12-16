package repo

import (
	"context"
	"database/sql"
	"fmt"
	"sync/atomic"

	_ "github.com/microsoft/go-mssqldb"

	"github.com/SoulStalker/data_bridge_ch/internal/config"
	"golang.org/x/time/rate"
)

type MSSQLRepo struct {
	db         *sql.DB
	limiter    *rate.Limiter
	queryCount atomic.Int64
}

func NewMSSQLRepo(cfg config.MSSQLConfig) (*MSSQLRepo, error) {
	db, err := sql.Open("sqlserver", cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("open mssql: %w", err)
	}
	db.SetMaxOpenConns(cfg.MaxConns)
	db.SetMaxIdleConns(cfg.MaxConns / 2)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping mssql: %w", err)
	}

	return &MSSQLRepo{
		db:      db,
		limiter: rate.NewLimiter(rate.Limit(cfg.RateLimit), cfg.RateLimit),
	}, nil
}

type TableInfo struct {
	SchemaName string
	Name       string
	RowCount   int64
}

func (r *MSSQLRepo) Tables(ctx context.Context) ([]TableInfo, error) {
	const query = `
		SELECT
			SCH.name AS SchemaName,
			T.name AS TableName,
			SUM(P.rows) AS [RowCount] 
		FROM sys.tables T
		INNER JOIN sys.schemas SCH ON T.schema_id = SCH.schema_id
		INNER JOIN sys.partitions P ON T.object_id = P.object_id
		WHERE P.index_id IN (0, 1)
		GROUP BY SCH.name, T.name
		ORDER BY SCH.name, T.name;
		`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query list tables: %w", err)
	}
	defer rows.Close()

	result := make([]TableInfo, 0)
	for rows.Next() {
		var tbl TableInfo
		if err := rows.Scan(&tbl.SchemaName, &tbl.Name, &tbl.RowCount); err != nil {
			return nil, fmt.Errorf("scan table info: %w", err)
		}
		result = append(result, tbl)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows err: %w", err)
	}

	return result, err
}
