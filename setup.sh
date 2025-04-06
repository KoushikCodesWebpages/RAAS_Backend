#!/bin/bash

# Colors for output
GREEN="\033[1;32m"
YELLOW="\033[1;33m"
BLUE="\033[1;34m"
RESET="\033[0m"

# Set the project name
PROJECT_NAME="RAAS"

echo -e "${YELLOW}🚀 Setting up your Gin-based backend...${RESET}"
echo -e "${BLUE}Project Name: ${GREEN}$PROJECT_NAME${RESET}\n"

######################################
# 1️⃣ Creating project structure
######################################

echo -e "${BLUE}📂 Creating project directories and files...${RESET}"

mkdir -p $PROJECT_NAME && cd $PROJECT_NAME
go mod init $PROJECT_NAME

mkdir -p config controllers data middlewares models public repositories responses routes security 

touch go.mod go.sum README.md air.toml
echo "# $PROJECT_NAME" > README.md

######################################
# 2️⃣ Adding working code to files
######################################

echo -e "${BLUE}📝 Writing working code into files...${RESET}"

# Populate .env file
cat <<EOF > .env
# Security Settings
SECRET_KEY=ebnIoNWjZ-YYpqAVPymzIyhHSzy5VrAP2ccxySB_Z9w

# Custom User Model (for authentication, if needed)
AUTH_USER_MODEL=accounts.AuthUser

# CORS Settings (list of allowed origins)
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://127.0.0.1:3000

# JWT Token Settings (authentication)
ACCESS_TOKEN_LIFETIME=60         # Lifetime of access tokens in minutes
REFRESH_TOKEN_LIFETIME=1440      # Lifetime of refresh tokens in minutes
ROTATE_REFRESH_TOKENS=true       # Enable token rotation after refresh
BLACKLIST_AFTER_ROTATION=true    # Blacklist the refresh token after rotation
AUTH_HEADER_TYPES=Bearer         # Token type used in Authorization header

# Email Backend Configuration (SMTP for sending emails)
EMAIL_BACKEND=smtp
EMAIL_HOST=smtp.gmail.com
EMAIL_PORT=587
EMAIL_USE_TLS=true
EMAIL_HOST_USER=koushikaltacc@gmail.com
EMAIL_HOST_PASSWORD=jyec rlpd llho myhc
DEFAULT_FROM_EMAIL=koushikaltacc@gmail.com

# Static and Media File Settings
STATIC_URL=/static/
MEDIA_URL=/media/
MEDIA_ROOT=./media               # Directory to store media files
STATIC_ROOT=./static             # Directory to store static files

#sqlite
DB_TYPE=sqlite
# Database Configuration (Azure Database)
DB_SERVER=raassr.database.windows.net
DB_PORT=1433
DB_USER=server
DB_PASSWORD=x9ttRfWYFoAK
DB_NAME=RAASDATABASE


# Server Configuration (for Gin)
SERVER_PORT=5000
SERVER_HOST=localhost

# Logging Configuration
LOG_LEVEL=debug

# Rate Limiting Settings (requests per minute)
RATE_LIMIT=100

# Environment (Development or Production)
ENVIRONMENT=development

# JWT Settings (if using JWT for authentication)
JWT_SECRET_KEY=yourJWTSecretKeyHere
JWT_EXPIRATION_TIME=3600  # 1 hour in seconds

# REST Framework-like Settings (for API)
REST_FRAMEWORK_DEFAULT_AUTHENTICATION_CLASSES=Bearer
REST_FRAMEWORK_DEFAULT_PERMISSION_CLASSES=IsAuthenticated
REST_FRAMEWORK_DEFAULT_PAGINATION_CLASS=CustomPagination
REST_FRAMEWORK_DEFAULT_FILTER_BACKENDS=DjangoFilterBackend
REST_FRAMEWORK_DEFAULT_RENDERER_CLASSES=JSONRenderer
REST_FRAMEWORK_DEFAULT_THROTTLE_CLASSES=AnonRateThrottle,UserRateThrottle
REST_FRAMEWORK_DEFAULT_THROTTLE_RATES_ANON=100/min
REST_FRAMEWORK_DEFAULT_THROTTLE_RATES_USER=100/min


EOF

# Create index.html
cat <<EOF > config/config.go
package config

import (
	"log"

	"github.com/spf13/viper"
	"github.com/joho/godotenv"
)

// Config struct holds all configuration fields loaded from .env
type Config struct {
	SecretKey              string
	AuthUserModel          string
	CORSAllowedOrigins     string
	AccessTokenLifetime    int
	RefreshTokenLifetime   int
	RotateRefreshTokens    bool
	BlacklistAfterRotation bool
	AuthHeaderTypes        string
	EmailBackend           string
	EmailHost              string
	EmailPort              int
	EmailUseTLS            bool
	EmailHostUser          string
	EmailHostPassword      string
	DefaultFromEmail       string
	StaticURL              string
	MediaURL               string
	MediaRoot              string
	StaticRoot             string
	DBType				   string
	DBServer               string
	DBPort                 int
	DBUser                 string
	DBPassword             string
	DBName                 string
	ServerPort             int
	ServerHost             string
	LogLevel               string
	RateLimit              int
	Environment            string
	JWTSecretKey           string
	JWTExpirationTime      int
	RestAuthClasses        string
	RestPermissionClasses  string
	RestPaginationClass    string
	RestFilterBackends     string
	RestRendererClasses    string
	RestThrottleClasses    string
	RestThrottleRatesAnon  string
	RestThrottleRatesUser  string
}

// InitConfig initializes viper configuration and loads environment variables
func InitConfig() (*Config, error) {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return nil, err
	}

	// Initialize viper
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	// Create and populate the Config struct
	config := &Config{
		SecretKey:              viper.GetString("SECRET_KEY"),
		AuthUserModel:          viper.GetString("AUTH_USER_MODEL"),
		CORSAllowedOrigins:     viper.GetString("CORS_ALLOWED_ORIGINS"),
		AccessTokenLifetime:    viper.GetInt("ACCESS_TOKEN_LIFETIME"),
		RefreshTokenLifetime:   viper.GetInt("REFRESH_TOKEN_LIFETIME"),
		RotateRefreshTokens:    viper.GetBool("ROTATE_REFRESH_TOKENS"),
		BlacklistAfterRotation: viper.GetBool("BLACKLIST_AFTER_ROTATION"),
		AuthHeaderTypes:        viper.GetString("AUTH_HEADER_TYPES"),
		EmailBackend:           viper.GetString("EMAIL_BACKEND"),
		EmailHost:              viper.GetString("EMAIL_HOST"),
		EmailPort:              viper.GetInt("EMAIL_PORT"),
		EmailUseTLS:            viper.GetBool("EMAIL_USE_TLS"),
		EmailHostUser:          viper.GetString("EMAIL_HOST_USER"),
		EmailHostPassword:      viper.GetString("EMAIL_HOST_PASSWORD"),
		DefaultFromEmail:       viper.GetString("DEFAULT_FROM_EMAIL"),
		StaticURL:              viper.GetString("STATIC_URL"),
		MediaURL:               viper.GetString("MEDIA_URL"),
		MediaRoot:              viper.GetString("MEDIA_ROOT"),
		StaticRoot:             viper.GetString("STATIC_ROOT"),
		DBType:              	viper.GetString("DB_TYPE"),
		DBServer:               viper.GetString("DB_SERVER"),
		DBPort:                 viper.GetInt("DB_PORT"),
		DBUser:                 viper.GetString("DB_USER"),
		DBPassword:             viper.GetString("DB_PASSWORD"),
		DBName:                 viper.GetString("DB_NAME"),
		ServerPort:             viper.GetInt("SERVER_PORT"),
		ServerHost:             viper.GetString("SERVER_HOST"),
		LogLevel:               viper.GetString("LOG_LEVEL"),
		RateLimit:              viper.GetInt("RATE_LIMIT"),
		Environment:            viper.GetString("ENVIRONMENT"),
		JWTSecretKey:           viper.GetString("JWT_SECRET_KEY"),
		JWTExpirationTime:      viper.GetInt("JWT_EXPIRATION_TIME"),
		RestAuthClasses:        viper.GetString("REST_FRAMEWORK_DEFAULT_AUTHENTICATION_CLASSES"),
		RestPermissionClasses:  viper.GetString("REST_FRAMEWORK_DEFAULT_PERMISSION_CLASSES"),
		RestPaginationClass:    viper.GetString("REST_FRAMEWORK_DEFAULT_PAGINATION_CLASS"),
		RestFilterBackends:     viper.GetString("REST_FRAMEWORK_DEFAULT_FILTER_BACKENDS"),
		RestRendererClasses:    viper.GetString("REST_FRAMEWORK_DEFAULT_RENDERER_CLASSES"),
		RestThrottleClasses:    viper.GetString("REST_FRAMEWORK_DEFAULT_THROTTLE_CLASSES"),
		RestThrottleRatesAnon:  viper.GetString("REST_FRAMEWORK_DEFAULT_THROTTLE_RATES_ANON"),
		RestThrottleRatesUser:  viper.GetString("REST_FRAMEWORK_DEFAULT_THROTTLE_RATES_USER"),
	}

	// Optional: Debugging output for configuration (useful during development)
	if config.Environment == "development" {
		log.Printf("Loaded configuration")
	}

	return config, nil
}

EOF

# Create index.html
cat <<EOF > models/db.go
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

	// Determine which database to use
	if cfg.DBType == "sqlite" {
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

	// Run Auto Migration
	log.Println("Starting AutoMigrate...")
	AutoMigrate()
	log.Println("AutoMigrate completed.")

	return DB
}



// AutoMigrate will automatically migrate all models to the database
func AutoMigrate() {
	err := DB.AutoMigrate(
		// Add all your model structs here
		&User{},
		&LinkedinJobLinks{}, // Example model
		// Add more models as needed
	)
	if err != nil {
		log.Fatalf("Error automigrating models: %v", err)
	}
}

// ResetDB drops the entire database and recreates it
func ResetDB() {
	log.Println("Resetting database...")

	// Step 1: Drop all foreign key constraints
	log.Println("Dropping foreign key constraints...")
	var constraints []struct {
		TableName      string
		ConstraintName string
	}

	DB.Raw(\`
		SELECT TABLE_NAME, CONSTRAINT_NAME 
		FROM INFORMATION_SCHEMA.TABLE_CONSTRAINTS 
		WHERE CONSTRAINT_TYPE = 'FOREIGN KEY'
	\`).Scan(&constraints)

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

	log.Println("Database reset completed.")
}

EOF

# Create main.go
cat <<EOF > main.go
package main

import (
	"fmt"
	"log"
	"RAAS/config"
	"RAAS/models"
	"RAAS/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode) 

	// Initialize the configuration
	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatalf("Error initializing config: %v", err)
	}

	// Initialize the database
	db := models.InitDB(cfg)

	// Create a new Gin router
	r := gin.Default()

	// Register all routes
	routes.SetupRoutes(r, db)

	// Start the server
	err = r.Run(fmt.Sprintf("%s:%d", cfg.ServerHost, cfg.ServerPort))
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}


EOF


# Create index.html
cat <<EOF > models/models.go
package models

import (
	"gorm.io/gorm"
)

// DB is the global database variable


// InitDB initializes the database connection and sets up the models

// AutoMigrate will automatically migrate all models to the database


// ResetDB will truncate tables, reset auto increment, and delete data, but keep table structure

// Models 
type User struct {
	ID    uint   \`gorm:"primaryKey"\`
	Name  string \`gorm:"size:100;not null"\`
	Email string \`gorm:"size:100;not null;uniqueIndex"\` // Ensures email stays unique without altering constraints
}

type LinkedinJobLinks struct {
	gorm.Model
	Title       string \`json:"title"\`
	Link        string \`gorm:"type:VARCHAR(1000);uniqueIndex" json:"link"\` // Changed from TEXT to VARCHAR(1000)
	Components  string \`gorm:"type:VARCHAR(2000)" json:"components"\`       // Changed from TEXT to VARCHAR(2000)
	Description string \`gorm:"type:VARCHAR(4000)" json:"description"\`       // Changed from TEXT to VARCHAR(4000)
}

EOF

cat <<EOF > repositories/general_repository.go
package repositories

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"net/url"
)

// GeneralRepository - Generic Repository for CRUD Operations
type GeneralRepository[T any] struct {
	db *gorm.DB
}

// NewGeneralRepository - Returns a new instance of GeneralRepository
func NewGeneralRepository[T any](db *gorm.DB) *GeneralRepository[T] {
	return &GeneralRepository[T]{db: db}
}

// Create - Adds a new record
func (r *GeneralRepository[T]) Create(entity *T) (*T, error) {
	if err := r.db.Create(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

// BulkCreate - Adds multiple records at once
func (r *GeneralRepository[T]) BulkCreate(entities *[]T) error {
	if err := r.db.Create(entities).Error; err != nil {
		return err
	}
	return nil
}

// GetByID - Fetch a record by ID
func (r *GeneralRepository[T]) GetByID(id uint) (*T, error) {
	var entity T
	if err := r.db.First(&entity, id).Error; err != nil {
		return nil, errors.New("record not found")
	}
	return &entity, nil
}

// GetAll - Fetch all records with optional filtering
func (r *GeneralRepository[T]) GetAll(queryParams url.Values) ([]T, error) {
	var entities []T
	query := r.db.Model(&entities)

	// Filtering based on query parameters
	for key, values := range queryParams {
		if len(values) > 0 {
			query = query.Where(fmt.Sprintf("%s = ?", key), values[0])
		}
	}

	// Pagination disabled for now
	// page, _ := strconv.Atoi(queryParams.Get("page"))
	// pageSize, _ := strconv.Atoi(queryParams.Get("page_size"))
	// if page == 0 { page = 1 }
	// if pageSize == 0 { pageSize = 10 }
	// offset := (page - 1) * pageSize
	// query.Offset(offset).Limit(pageSize)

	query.Find(&entities)
	return entities, nil
}

// Update - Updates an existing record
func (r *GeneralRepository[T]) Update(id uint, entity *T) (*T, error) {
	if err := r.db.Model(&entity).Where("id = ?", id).Updates(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

// Delete - Removes a record by ID
func (r *GeneralRepository[T]) Delete(id uint) error {
	var entity T
	if err := r.db.Delete(&entity, id).Error; err != nil {
		return errors.New("failed to delete record")
	}
	return nil
}

EOF

cat <<EOF > controllers/general_controllers.go
package controllers

import (
	"github.com/gin-gonic/gin"
	"RAAS/repositories"
	"net/http"
	"strconv"
)

// GeneralController - Generic Controller for CRUD Operations
type GeneralController[T any] struct {
	repo *repositories.GeneralRepository[T]
}

// NewGeneralController - Returns a new instance of GeneralController
func NewGeneralController[T any](repo *repositories.GeneralRepository[T]) *GeneralController[T] {
	return &GeneralController[T]{repo: repo}
}

// Create - Handles the creation of a new entity
func (gc *GeneralController[T]) Create(c *gin.Context) {
	var entity T
	if err := c.ShouldBindJSON(&entity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	createdEntity, err := gc.repo.Create(&entity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create entity"})
		return
	}
	c.JSON(http.StatusOK, createdEntity)
}

// BulkCreate - Handles the creation of multiple entities
func (gc *GeneralController[T]) BulkCreate(c *gin.Context) {
	var entities []T
	if err := c.ShouldBindJSON(&entities); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err := gc.repo.BulkCreate(&entities)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create entities"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Bulk data inserted successfully"})
}

func (gc *GeneralController[T]) UploadCSV(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file upload"})
		return
	}
	defer file.Close()

	// Parse CSV data
	records, err := parseCSV[T](file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse CSV"})
		return
	}

	// Perform bulk insert
	err = gc.repo.BulkCreate(&records)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert records"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "CSV data uploaded successfully"})
}

// GetByID - Fetch a record by ID
func (gc *GeneralController[T]) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	entity, err := gc.repo.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}

	c.JSON(http.StatusOK, entity)
}

// GetAll - Fetch all records with optional filtering
func (gc *GeneralController[T]) GetAll(c *gin.Context) {
	entities, err := gc.repo.GetAll(c.Request.URL.Query())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch records"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		// "total": total, // Commented out pagination for now
		"data": entities,
	})
}

// Update - Updates a record
func (gc *GeneralController[T]) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var entity T
	if err := c.ShouldBindJSON(&entity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	updatedEntity, err := gc.repo.Update(uint(id), &entity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update record"})
		return
	}

	c.JSON(http.StatusOK, updatedEntity)
}

// Delete - Removes a record by ID
func (gc *GeneralController[T]) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	err = gc.repo.Delete(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Record deleted"})
}

EOF


cat <<EOF > controllers/parsecsv.go
package controllers

import (
	"encoding/csv"
	"errors"
	"io"
	"reflect"
	"strconv"
	"strings"
)

// parseCSV - Parses CSV data into a slice of the given model type
func parseCSV[T any](file io.Reader) ([]T, error) {
	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	// Read all lines
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(lines) < 2 {
		return nil, errors.New("CSV file must have at least one data row")
	}

	// Get headers
	headers := lines[0]
	var records []T

	// Iterate through each row
	for _, row := range lines[1:] {
		var record T
		val := reflect.ValueOf(&record).Elem()

		for i, field := range headers {
			field = strings.ToLower(strings.TrimSpace(field))
			fieldVal := strings.TrimSpace(row[i])

			// Assign values based on field types
			structField := val.FieldByName(strings.Title(field))
			if !structField.IsValid() || !structField.CanSet() {
				continue
			}

			switch structField.Kind() {
			case reflect.String:
				structField.SetString(fieldVal)
			case reflect.Int, reflect.Int64:
				intVal, _ := strconv.Atoi(fieldVal)
				structField.SetInt(int64(intVal))
			default:
				continue
			}
		}

		records = append(records, record)
	}

	return records, nil
}

EOF


cat <<EOF > routes/routes.go
package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"RAAS/controllers"
	"RAAS/repositories"
	"RAAS/models"
)

// SetupRoutes - Registers all routes
func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	// Home Route
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Server is running"})
	})

	// Register CRUD routes for LinkedinJobLinks
	SetupGenericRoutes[models.LinkedinJobLinks](r, db, "/linkedin-job-links")
}


// SetupGenericRoutes - Registers generic CRUD routes for any model
func SetupGenericRoutes[T any](r *gin.Engine, db *gorm.DB, baseRoute string) {
	repo := repositories.NewGeneralRepository[T](db)
	controller := controllers.NewGeneralController[T](repo)

	group := r.Group(baseRoute)
	{
		group.POST("/", controller.Create)
		group.POST("/bulk", controller.BulkCreate)
		group.POST("/upload-csv", controller.UploadCSV) 
		group.GET("/:id", controller.GetByID)
		group.GET("/", controller.GetAll)
		group.PUT("/:id", controller.Update)
		group.DELETE("/:id", controller.Delete)
	}
}

EOF





# .gitignore setup
echo "/bin/" >> .gitignore
echo "/pkg/mod/" >> .gitignore
echo ".env" >> .gitignore

######################################
# 3️⃣ Installing Go packages
######################################

echo -e "\n${GREEN}📦 Installing required Go packages...${RESET}"
go get -u github.com/gin-gonic/gin
go get -u gorm.io/gorm
go get -u github.com/joho/godotenv
go get -u gorm.io/driver/sqlite
go get -u github.com/spf13/viper
go get -u github.com/golang-jwt/jwt/v5
go get -u gorm.io/driver/sqlserver

######################################
# 4️⃣ Setting up Air for live reload
######################################

echo -e "\n${GREEN}💨 Installing Air for auto-reloading...${RESET}\n"
go install github.com/air-verse/air@latest


# Ensure Air is in PATH permanently
AIR_PATH="\$HOME/go/bin"
if [[ ":$PATH:" != *":$AIR_PATH:"* ]]; then
    echo "export PATH=\$PATH:$AIR_PATH" >> ~/.bashrc
    source ~/.bashrc
fi

# Create air.toml
cat <<EOF > air.toml
# Config file for Air
[build]
  cmd = "go build -o tmp\\main.exe ."  # Use Windows path
  bin = "tmp\\main.exe"
  full_bin = "tmp\\main.exe"  # Corrected Windows path
  include_ext = ["go", "tpl", "tmpl", "html"]
  exclude_dir = ["tmp", "vendor", "node_modules"]
  delay = 200


[log]
  time = true
  level = "debug"

[serve]
  watch_dir = ["."]
  restart = true
  kill_delay = "1s"

[screen]
  clear = true

EOF

######################################
# 🎉 Setup Complete
######################################

echo -e "\n${YELLOW}✨ Setup complete! ${RESET}\n"
echo -e "${BLUE}📌 Next Steps:${RESET}"
echo -e "${GREEN}1️⃣  Navigate to the project folder: ${YELLOW}cd $PROJECT_NAME${RESET}"
echo -e "${GREEN}2️⃣  Start the server with auto-reload: ${YELLOW}air${RESET} (or ${YELLOW}~/go/bin/air${RESET} if air isn't recognized)\n\n"

echo -e "\n${GREEN}💨 Installing the build file${RESET}\n"
go build -o tmp/main main.go

