package models

import (
	"fmt"
	"log"
	"github.com/glebarez/sqlite"
	//"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"RAAS/config"
)

// DB is the global database variable
var DB *gorm.DB
func InitDB(cfg *config.Config) *gorm.DB {
	//log.Println("Starting database initialization...")

	var err error

	// Using SQLite only (no conditional)
	log.Println("Using SQLite for development")

	dbPath := cfg.DBName
	if dbPath == "" {
		dbPath = "RAASDATABASE" // fallback to file-based DB
	}
	//log.Printf("SQLite DB Path: %s", dbPath)

	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})

	// === Commented out: SQL Server logic ===
	// var dbType = cfg.DBType
	// if dbType == "sqlserver" {
	// 	log.Println("Using Azure SQL Database")
	// 	dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s",
	// 		cfg.DBUser,
	// 		cfg.DBPassword,
	// 		cfg.DBServer,
	// 		cfg.DBPort,
	// 		cfg.DBName,
	// 	)
	// 	log.Printf("DSN Built: %s", dsn)
	// 	DB, err = gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	// }

	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	//log.Println("Database connection successful.")
	ResetDB(DB, "sqlite", cfg.DBName, []string{
		// "auth_users",
		// "seekers",
		// "admins",
		//"preferred_job_titles",
		"linked_in_job_meta_data",
		"xing_job_meta_data",
		"linked_in_failed_jobs",
		"xing_failed_jobs",
		"linked_in_job_application_links",
		"xing_job_application_links",
		"linked_in_job_descriptions",
		"xing_job_descriptions",
		"job_match_scores",
	})

	//log.Println("Starting AutoMigrate...")
	AutoMigrate()
	//log.Println("AutoMigrate completed seeding starts.")

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

		&PersonalInfo{},
		&ProfessionalSummary{},
		&PreferredJobTitle{},

		&JobMatchScore{},


		&LinkedInJobMetaData{},
		&XingJobMetaData{},
		&LinkedInFailedJob{},
		&XingFailedJob{},
		&LinkedInJobApplicationLink{},
		&XingJobApplicationLink{},
		&LinkedInJobDescription{},
		&XingJobDescription{},
		
		// Add more models as needed
	)
	if err != nil {
		log.Fatalf("Error automigrating models: %v", err)
	}
}

// ResetDB drops selected tables for SQL Server and SQLite
func ResetDB(DB *gorm.DB, dbType string, dbName string, tablesToDrop []string) {
	//log.Println("ResetDB started...")

	switch dbType {
	case "sqlserver":
		//log.Println("Detected SQL Server, proceeding with selective reset...")

		// Print all existing tables before reset
		//log.Println("Listing all tables before reset...")

		var tableNames []string
		DB.Raw(`
			SELECT TABLE_NAME 
			FROM INFORMATION_SCHEMA.TABLES 
			WHERE TABLE_TYPE = 'BASE TABLE'
		`).Scan(&tableNames)
		// for _, table := range tableNames {
		// 	log.Println("Existing table:", table)
		// }

		// Drop foreign key constraints
		//log.Println("Dropping foreign key constraints...")
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
			if contains(tablesToDrop, c.TableName) {
				dropFKQuery := fmt.Sprintf("ALTER TABLE [%s] DROP CONSTRAINT [%s];", c.TableName, c.ConstraintName)
				if err := DB.Exec(dropFKQuery).Error; err != nil {
					log.Printf("Error dropping foreign key %s on table %s: %v", c.ConstraintName, c.TableName, err)
				} else {
					//log.Printf("Dropped foreign key %s on table %s", c.ConstraintName, c.TableName)
				}
			}
		}

		// Drop selected tables
		//log.Println("Dropping selected tables...")
		for _, table := range tablesToDrop {
			dropTableQuery := fmt.Sprintf("DROP TABLE IF EXISTS [%s];", table)
			if err := DB.Exec(dropTableQuery).Error; err != nil {
				log.Printf("Error dropping table %s: %v", table, err)
			} else {
				log.Printf("Dropped table %s successfully", table)
			}
		}

	case "sqlite":
		//log.Println("Detected SQLite, proceeding with selective reset...")

		// Print all existing tables before reset
		//log.Println("Listing all tables before reset...")
		var tableNames []string
		DB.Raw(`SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%'`).Scan(&tableNames)
		for _, table := range tableNames {
			log.Println("Existing table:", table)
		}

		// Drop selected tables
		for _, table := range tablesToDrop {
			dropTableQuery := fmt.Sprintf("DROP TABLE IF EXISTS `%s`;", table)
			if err := DB.Exec(dropTableQuery).Error; err != nil {
				log.Printf("Error dropping table %s: %v", table, err)
			} else {
				//log.Printf("Dropped table %s successfully", table)
			}
		}

	default:
		log.Printf("ResetDB not supported for dbType: %s", dbType)
	}

	//log.Println("ResetDB completed.")
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

