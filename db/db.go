package db

import (
	"log"
	"os"

	walletmodal "example.com/smart-contract/repository/modals"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/joho/godotenv"
)

func setupTables(db *pg.DB) {
	walletModalErr := db.Model(&walletmodal.Wallet{}).CreateTable(&orm.CreateTableOptions{
		IfNotExists: true,
	})

	if walletModalErr != nil {
		panic(walletModalErr)
	}

	defer db.Close()
}

func Connection() *pg.DB {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading configuration from .env file. Please check your configuration.")
	}

	options := &pg.Options{
		User:     os.Getenv("DATABASE_USER"),
		Password: os.Getenv("DATABASE_PASSWORD"),
		Database: os.Getenv("DATABASE_NAME"),
	}

	db := pg.Connect(options)

	if db == nil {
		log.Fatal("Database Connection Error. Please check your configuration.")
	}

	return db
}

func InitializeDatabase() *pg.DB {
	dbConnection := Connection()

	setupTables(dbConnection)

	return dbConnection
}
