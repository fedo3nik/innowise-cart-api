package main

import (
	_ "cart-api/docs"
	"cart-api/internal/pkg/common/db"
	"cart-api/internal/pkg/config"
	"cart-api/internal/server"
)

//	@title		Trainee cart-Api
//	@version	1.0

// host localhost:3000
// BasePath /
func main() {
	cfg, err := config.InitConfig()

	if err != nil {
		panic(err)
	}

	dbPool, err := db.InitDb()

	if err != nil {
		panic(err)
	}

	defer dbPool.Close()

	app := server.NewServer(&cfg, dbPool)

	if err := app.Run(); err != nil {
		panic(err)
	}
}
