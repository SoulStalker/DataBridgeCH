package main

import (
	"context"
	"fmt"

	"github.com/SoulStalker/data_bridge_ch/internal/config"
	repo "github.com/SoulStalker/data_bridge_ch/internal/mssql"
)

func main() {
	cfg := config.MustLoad("./config/config.yaml")
	db, err := repo.NewMSSQLRepo(cfg.MSSQL)
	if err != nil {
		fmt.Println(err)
	}

	ctx := context.Background()
	rows, err := db.Tables(ctx)
	if err != nil {
		fmt.Println(err)
	}
	for _, row := range rows {
		fmt.Println(row.Name, row.RowCount)
	}
}
