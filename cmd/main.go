package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/freeloginname/otusGoBasicProject/pkg/notes"
	"github.com/freeloginname/otusGoBasicProject/pkg/pgdb"
	"github.com/freeloginname/otusGoBasicProject/pkg/ui"
	"github.com/freeloginname/otusGoBasicProject/pkg/users"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/spf13/viper"
)

var DSN string

// var embedMigrations embed.FS

type connectionFunc func(context.Context, string, int32) (*pgxpool.Pool, error)

func retry(
	ctx context.Context,
	attempts int,
	sleep time.Duration,
	f connectionFunc,
	dbDSN string,
	maxOpenConns int32,
) (h *pgxpool.Pool, err error) {
	for i := 0; i < attempts; i++ {
		if i > 0 {
			log.Println("retrying after error:", err)
			time.Sleep(sleep)
			sleep *= 2
		}
		h, err = f(ctx, dbDSN, maxOpenConns)
		if err == nil {
			return h, nil
		}
	}
	return nil, fmt.Errorf("after %d attempts, last error: %w", attempts, err)
}

func migrate(h *pgxpool.Pool, migrationDir string) (err error) {
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	gooseDB := stdlib.OpenDBFromPool(h)
	if err := goose.Up(gooseDB, migrationDir); err != nil {
		return err
	}
	return nil
}

func main() {
	curPath, _ := os.Getwd()
	fmt.Printf("Current path: %s\n", curPath)
	envPath := filepath.Join(curPath, ".env")

	viper.SetConfigFile(".env")
	viper.SetConfigFile(".")
	viper.SetConfigFile(envPath)
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	port := viper.Get("APP_HTTP_PORT").(string)
	dbURL := viper.Get("DB_DSN").(string)
	envMigrationDir := viper.Get("MIGRATION_DIR_DOTLESS").(string)
	migrationDir := filepath.Join(curPath, envMigrationDir)

	ctx := context.Background()
	h, err := retry(ctx, 5, 5*time.Second, pgdb.New, dbURL, 20)
	if err != nil {
		panic(fmt.Errorf("failed to connect to DB: %w", err))
	}
	// h, err := pgdb.New(ctx, dbURL, 20)
	// if err != nil {
	// 	panic(fmt.Errorf("failed to connect to DB: %v", err))
	// }

	// migration
	// goose.SetBaseFS(embedMigrations)
	// if err := goose.SetDialect("postgres"); err != nil {
	// 	panic(err)
	// }
	// gooseDB := stdlib.OpenDBFromPool(h)
	// if err := goose.Up(gooseDB, migrationDir); err != nil {
	// 	panic(err)
	// }

	err = migrate(h, migrationDir)
	if err != nil {
		panic(err)
	}

	// routes
	r := gin.Default()
	users.RegisterRoutes(r, h, []byte(viper.Get("SECRET").(string)))
	notes.RegisterRoutes(r, h, []byte(viper.Get("SECRET").(string)))
	ui.RegisterRoutes(r, h, []byte(viper.Get("SECRET").(string)))
	r.LoadHTMLGlob("templates/*.tmpl")
	r.Static("/css", "./static/css")
	r.Static("/assets", "./static/assets")
	r.Static("/js", "./static/js")
	r.Run(port)
}
