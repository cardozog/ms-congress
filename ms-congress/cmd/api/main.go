package main

import (
	"context"
	"log"
	"ms-congress/internal/config/app"
	"ms-congress/internal/config/db"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	dbCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	gormDB, err := db.InitDB(dbCtx)
	if err != nil {
		log.Fatalf("connect to PostgreSQL: %v", err)
	}
	sqlDatabase, err := gormDB.DB()
	if err != nil {
		log.Fatalf("get PostgreSQL connection: %v", err)
	}
	defer sqlDatabase.Close()

	if err := db.RunMigrations(dbCtx, gormDB); err != nil {
		log.Fatalf("run PostgreSQL migrations: %v", err)
	}

	router, err := app.InitApp(ctx, gormDB)
	if err != nil {
		log.Fatalf("init app: %v", err)
	}
	router.Run(":8080")

}
