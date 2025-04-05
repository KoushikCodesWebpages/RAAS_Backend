package models

import (
	"fmt"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"RAAS/config"
)

// DB is the global database variable
var DB *gorm.DB

// InitDB initializes the database connection and sets up the models
func InitDB(cfg *config.Config) *gorm.DB {
	log.Println("Starting database initialization...")

	var err error
	var dbType = cfg.DBType // Get the DB type from the config

	// Determine which database to use
	if dbType == "sqlite" {
		log.Println("Using SQLite for development")
		DB, err = gorm.Open(sqlite.Open(cfg.DBName), &gorm.Config{})
	} else {
		log.Println("Using Azure SQL Database")
		dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s",
			cfg.DBUser,
			cfg.DBPassword,
			cfg.DBServer,
			cfg.DBPort,
			cfg.DBName,
		)
		log.Println("DSN Built:", dsn)
		DB, err = gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	}

	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	log.Println("Database connection successful.")
	ResetDB(DB, dbType, cfg.DBName) // Pass the dbName as parameter to ResetDB

	log.Println("Starting AutoMigrate...")
	AutoMigrate()
	log.Println("AutoMigrate completed.")

	return DB
}

// AutoMigrate will automatically migrate all models to the database
func AutoMigrate() {
	err := DB.AutoMigrate(
		// Add all your model structs here
		&AuthUser{},
		&Seeker{},
		&Admin{},
		&LinkedinJobMetadata{}, // Example model
		// Add more models as needed
	)
	if err != nil {
		log.Fatalf("Error automigrating models: %v", err)
	}
}

// ResetDB drops the entire database and recreates it for SQL Server
func ResetDB(DB *gorm.DB, dbType string, dbName string) {
	

	// If SQL Server, proceed with normal reset
	if dbType == "sqlserver" {
		log.Println("Detected SQL Server, proceeding with full reset...")

		// Step 1: Drop foreign key constraints
		log.Println("Dropping foreign key constraints...")
		var constraints []struct {
			TableName      string
			ConstraintName string
		}

		DB.Raw(`
			SELECT TABLE_NAME, CONSTRAINT_NAME 
			FROM INFORMATION_SCHEMA.TABLE_CONSTRAINTS 
			WHERE CONSTRAINT_TYPE = 'FOREIGN KEY'
		`).Scan(&constraints)

		for _, c := range constraints {
			if c.TableName == "" || c.ConstraintName == "" {
				log.Println("Skipping invalid foreign key constraint with empty values")
				continue
			}

			dropFKQuery := fmt.Sprintf("ALTER TABLE [%s] DROP CONSTRAINT [%s];", c.TableName, c.ConstraintName)
			if err := DB.Exec(dropFKQuery).Error; err != nil {
				log.Printf("Error dropping foreign key %s on table %s: %v", c.ConstraintName, c.TableName, err)
			} else {
				log.Printf("Dropped foreign key %s on table %s", c.ConstraintName, c.TableName)
			}
		}

		// Step 2: Get all table names dynamically
		var tables []string
		DB.Raw("SELECT TABLE_NAME FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_TYPE = 'BASE TABLE'").Scan(&tables)

		// Step 3: Drop tables
		log.Println("Dropping tables...")
		for _, table := range tables {
			dropTableQuery := fmt.Sprintf("DROP TABLE IF EXISTS [%s];", table)
			if err := DB.Exec(dropTableQuery).Error; err != nil {
				log.Printf("Error dropping table %s: %v", table, err)
			} else {
				log.Printf("Dropped table %s successfully", table)
			}
		}

		log.Println("SQL Server database reset completed.")
	}
}
