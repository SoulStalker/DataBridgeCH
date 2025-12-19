package main

import (
	"context"
	"fmt"
	"log"
	"slices"

	"github.com/SoulStalker/data_bridge_ch/internal/clickhouse"
	"github.com/SoulStalker/data_bridge_ch/internal/config"
	"github.com/SoulStalker/data_bridge_ch/internal/model"
	"github.com/SoulStalker/data_bridge_ch/internal/mssql"
	"github.com/google/uuid"
)

func main() {
	ROWS_1C := []string{"_InfoRg5404", "_InfoRg5415", "_InfoRg4432", "_InfoRg6453", "_AccumRg10970", "_AccumRg5628"}

	cfg := config.MustLoad("./config/config.yaml")
	db, err := mssql.NewMSSQLRepo(cfg.MSSQL)
	if err != nil {
		fmt.Println(err)
	}

	ctx := context.Background()
	rows, err := db.Tables(ctx)
	if err != nil {
		fmt.Println(err)
	}
	for _, row := range rows {
		if slices.Contains(ROWS_1C, row.Name) {
			fmt.Printf("%s.%s\n", row.SchemaName, row.Name)
		}
	}

	goods := make(chan []model.Reference379, 10)
	go func() {
		if err := db.FetchProducts(ctx, 100, goods); err != nil {
			log.Printf("fetch goods err: %v", err)
		}
	}()

	for batch := range goods {
		for _, p := range batch {
			fmt.Printf("ID: %s, Code: %s, Description: %s\n", guidToString(p.ID), p.Code, p.Description)
		}
	}

	err = clickhouse.InitDB(cfg.ClickHouse)
	if err != nil {
		fmt.Println(err)
	}
}

// TODO: move this func
func guidToString(b []byte) string {
	if len(b) != 16 {
		return ""
	}
	u, err := uuid.FromBytes(b)
	if err != nil {
		return ""
	}
	return u.String()
}
