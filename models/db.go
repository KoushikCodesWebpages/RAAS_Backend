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
	ResetDB(DB, dbType, cfg.DBName,[]string{
	}) // Pass the dbName as parameter to ResetDB

	log.Println("Starting AutoMigrate...")
	AutoMigrate()
	log.Println("AutoMigrate completed.")

	SeedJobs(DB)
	

	return DB
}

// AutoMigrate will automatically migrate all models to the database
func AutoMigrate() {
	err := DB.AutoMigrate(
		// Add all your model structs here
		&AuthUser{},
		&Seeker{},
		&Admin{},
		&PreferredJobTitle{},
		&LinkedInJobMetaData{},
		&XingJobMetaData{},
		&LinkedInFailedJob{},
		&XingFailedJob{},
		&LinkedInJobApplicationLink{},
		&XingJobApplicationLink{},
		// Add more models as needed
	)
	if err != nil {
		log.Fatalf("Error automigrating models: %v", err)
	}
}

// ResetDB drops the entire database and recreates it for SQL Server
func ResetDB(DB *gorm.DB, dbType string, dbName string, tablesToDrop []string) {
	if dbType == "sqlserver" {
		log.Println("Detected SQL Server, proceeding with selective reset...")

		// Step 1: Drop foreign key constraints only for specified tables
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
			// Only drop constraints from tables we want to reset
			if contains(tablesToDrop, c.TableName) {
				dropFKQuery := fmt.Sprintf("ALTER TABLE [%s] DROP CONSTRAINT [%s];", c.TableName, c.ConstraintName)
				if err := DB.Exec(dropFKQuery).Error; err != nil {
					log.Printf("Error dropping foreign key %s on table %s: %v", c.ConstraintName, c.TableName, err)
				} else {
					log.Printf("Dropped foreign key %s on table %s", c.ConstraintName, c.TableName)
				}
			}
		}

		// Step 2: Drop only selected tables
		log.Println("Dropping selected tables...")
		for _, table := range tablesToDrop {
			dropTableQuery := fmt.Sprintf("DROP TABLE IF EXISTS [%s];", table)
			if err := DB.Exec(dropTableQuery).Error; err != nil {
				log.Printf("Error dropping table %s: %v", table, err)
			} else {
				log.Printf("Dropped table %s successfully", table)
			}
		}

		log.Println("Selected table reset completed.")
	}
}

// Helper to check if a string exists in a list
func contains(list []string, val string) bool {
	for _, item := range list {
		if item == val {
			return true
		}
	}
	return false
}


	//DEVELOPMENT
	func SeedJobs(db *gorm.DB) {
		// Check if jobs already seeded
		var count int64
		db.Model(&LinkedInJobMetaData{}).Count(&count)
		if count > 0 {
			return // Already seeded
		}
	
		// Dummy LinkedIn jobs
		linkedinJobs := []LinkedInJobMetaData{
			{ID: "lnk1", JobID: "L001", Title: "Software Engineer", Company: "LinkedIn", Location: "Berlin", PostedDate: "2024-04-01", Link: "https://linkedin.com/jobs/1", Processed: true},
			{ID: "lnk2", JobID: "L002", Title: "Backend Developer", Company: "Google", Location: "Munich", PostedDate: "2024-04-02", Link: "https://linkedin.com/jobs/2", Processed: true},
			// Add 3-4 more...
		}
	
		xingJobs := []XingJobMetaData{
			{ID: "xg1", JobID: "X001", Title: "Data Engineer", Company: "Xing", Location: "Hamburg", PostedDate: "2024-04-01", Link: "https://xing.com/jobs/1", Processed: true},
			{ID: "xg2", JobID: "X002", Title: "DevOps Engineer", Company: "SAP", Location: "Berlin", PostedDate: "2024-04-02", Link: "https://xing.com/jobs/2", Processed: true},
			// Add 3-4 more...
		}
	
		// Insert into DB
		db.Create(&linkedinJobs)
		db.Create(&xingJobs)
	}