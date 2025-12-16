package main

import (
	"fmt"

	"github.com/SoulStalker/data_bridge_ch/internal/config"
	"github.com/SoulStalker/data_bridge_ch/internal/migrator/repo"
)

func main() {
	cfg := config.MustLoad("./config/config.yaml")
	_, err := repo.NewMSSQLRepo(cfg.MSSQL)
	if err != nil {
		fmt.Println(err)
	}

}
