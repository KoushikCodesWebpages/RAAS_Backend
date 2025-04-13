package models

import (
	"fmt"
	"log"
	"RAAS/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
)

// DB is the global database variable
var DB *gorm.DB

func InitDB(cfg *config.Config) *gorm.DB {
	var err error

	log.Println("Using MySQL database")

	// Set the GORM logger to silent or info level as needed
	gormLogger := logger.New(
		log.New(os.Stdout, "", log.LstdFlags), // Output, Prefix, and Flags
		logger.Config{
			LogLevel: logger.Silent, // Change to logger.Info to show logs, logger.Silent to hide
			Colorful: true,
		},
	)

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBServer,
		cfg.DBPort,
		cfg.DBName,
	)

	// Open the DB with the custom logger
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		log.Fatalf("❌ Error connecting to MySQL: %v", err)
	}

	log.Println("✅ MySQL connection established")

	ResetDB(DB, []string{
	// "auth_users",
	// "seekers",
	// "admins",
	// "preferred_job_titles",
	// "personal_infos",
	// "professional_summaries",
	// "work_experiences",
	// "educations",
	// "languages",
	// "certificates",
	// "user_entry_timelines",
	// "linked_in_job_meta_data",
	// "xing_job_meta_data",
	// "linked_in_job_application_links",
	// "xing_job_application_links",
	// "linked_in_job_descriptions",
	// "xing_job_descriptions",
	// "job_match_scores",
	// "selected_job_applications",
	})

	AutoMigrate()
	SeedJobs(DB)
	PrintAllTables(DB, cfg.DBName)
	return DB
}

// AutoMigrate will automatically migrate all models to the database
func AutoMigrate() {
	// Phase 1: Create Major Tables
	err := DB.AutoMigrate(
		// Major Tables (without foreign key dependencies)
		&AuthUser{},
		&Seeker{},
		&Admin{},
		&PersonalInfo{},
		&ProfessionalSummary{},
		&WorkExperience{},
		&Education{},
		&Certificate{},
		&Language{},
		&PreferredJobTitle{},

	)

	if err != nil {
		log.Fatalf("Error creating major tables: %v", err)
	}
	log.Println("Major tables migration completed successfully")

	// Phase 2: Create Foreign Key Related Tables
	err = DB.AutoMigrate(
		// Job-related tables
		&UserEntryTimeline{},
		&LinkedInJobMetaData{},
		&XingJobMetaData{},
		&JobMatchScore{},

		// Foreign Key Dependent Tables
		&LinkedInJobApplicationLink{},
		&XingJobApplicationLink{},
		&LinkedInJobDescription{},
		&XingJobDescription{},

		&SelectedJobApplication{},
	)

	if err != nil {
		log.Fatalf("Error creating foreign key related tables: %v", err)
	}
	log.Println("Foreign key related tables migration completed successfully")
}

func ResetDB(DB *gorm.DB, tablesToDrop []string) {
	log.Println("Resetting selected tables...")

	// Get the dialect (MySQL, SQLServer, SQLite)
	dialect := DB.Dialector.Name()

	// Fetch foreign key constraints for tables in the drop list
	var constraints []struct {
		TableName      string
		ConstraintName string
	}

	// Get foreign key constraints from the information schema
	DB.Raw(`
		SELECT TABLE_NAME, CONSTRAINT_NAME
		FROM INFORMATION_SCHEMA.TABLE_CONSTRAINTS
		WHERE CONSTRAINT_TYPE = 'FOREIGN KEY'
	`).Scan(&constraints)

	// Drop foreign key constraints before dropping tables
	for _, c := range constraints {
		if contains(tablesToDrop, c.TableName) {
			var query string
			switch dialect {
			case "mysql", "sqlite":
				// MySQL/SQLite: Drop foreign key constraint
				query = fmt.Sprintf("ALTER TABLE `%s` DROP FOREIGN KEY `%s`;", c.TableName, c.ConstraintName)
			case "sqlserver":
				// SQLServer: Drop foreign key constraint
				query = fmt.Sprintf("ALTER TABLE %s DROP CONSTRAINT %s;", c.TableName, c.ConstraintName)
			default:
				log.Printf("⚠️ Unknown dialect: %s, skipping foreign key removal.", dialect)
				continue
			}

			if err := DB.Exec(query).Error; err != nil {
				log.Printf("⚠️ Error dropping FK %s on table %s: %v", c.ConstraintName, c.TableName, err)
			}
		}
	}

	// Now, drop the selected tables
	for _, table := range tablesToDrop {
		var query string
		switch dialect {
		case "mysql", "sqlite":
			// MySQL/SQLite: Drop tables
			query = fmt.Sprintf("DROP TABLE IF EXISTS `%s`;", table)
		case "sqlserver":
			// SQLServer: Drop tables (without IF EXISTS)
			query = fmt.Sprintf("DROP TABLE %s;", table)
		default:
			log.Printf("⚠️ Unknown dialect: %s, skipping table drop.", dialect)
			continue
		}

		if err := DB.Exec(query).Error; err != nil {
			log.Printf("⚠️ Error dropping table %s: %v", table, err)
		}
	}

	log.Println("✅ Tables reset")
}

// contains checks if a string exists in a slice.
func contains(list []string, val string) bool {
	for _, item := range list {
		if item == val {
			return true
		}
	}
	return false
}

func PrintAllTables(db *gorm.DB, dbName string) {
	var tables []string
	err := db.Raw("SELECT table_name FROM information_schema.tables WHERE table_schema = ?", dbName).Scan(&tables).Error
	if err != nil {
		log.Fatalf("❌ Error fetching table names: %v", err)
	}
	log.Println("📦 Tables in the database:")
	for _, table := range tables {
		log.Println(" -", table)
	}
}
