package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	SecretKey              string
	AuthUserModel          string
	CORSAllowedOrigins     string
	ClientId               string
	ClientSecret           string
	FrontendBaseUrl        string
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
	DBType                 string
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

func InitConfig() (*Config, error) {
	// Load .env file only if not running on Railway (or similar env)
	if _, exists := os.LookupEnv("RAILWAY_ENVIRONMENT"); !exists {
		err := godotenv.Load()
		if err != nil {
			log.Println("No .env file found, relying on system environment variables")
		}
	}

	// Initialize viper
	viper.AutomaticEnv() // Automatically bind to environment variables

	// Create and populate the Config struct
	config := &Config{
		SecretKey:              viper.GetString("SECRET_KEY"),
		AuthUserModel:          viper.GetString("AUTH_USER_MODEL"),
		CORSAllowedOrigins:     viper.GetString("CORS_ALLOWED_ORIGINS"),
		ClientId:               viper.GetString("CLIENT_ID"),
		ClientSecret:           viper.GetString("CLIENT_SECRET"),
		FrontendBaseUrl:        viper.GetString("FRONTEND_BASE_URL"),
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
		DBType:                 viper.GetString("DB_TYPE"),
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

	// if config.Environment == "development" {
	// 	log.Printf("Loaded configuration: %+v", config)
	// }

	return config, nil
}

//for development
// package config

// import (
// 	"log"
// 	"os"

// 	"github.com/joho/godotenv"
// 	"github.com/spf13/viper"
// )

// type Config struct {
// 	SecretKey              string
// 	AuthUserModel          string
// 	CORSAllowedOrigins     string
// 	ClientId               string
// 	ClientSecret           string
// 	FrontendBaseUrl        string
// 	AccessTokenLifetime    int
// 	RefreshTokenLifetime   int
// 	RotateRefreshTokens    bool
// 	BlacklistAfterRotation bool
// 	AuthHeaderTypes        string
// 	EmailBackend           string
// 	EmailHost              string
// 	EmailPort              int
// 	EmailUseTLS            bool
// 	EmailHostUser          string
// 	EmailHostPassword      string
// 	DefaultFromEmail       string
// 	StaticURL              string
// 	MediaURL               string
// 	MediaRoot              string
// 	StaticRoot             string
// 	DBType                 string
// 	DBServer               string
// 	DBPort                 int
// 	DBUser                 string
// 	DBPassword             string
// 	DBName                 string
// 	ServerPort             int
// 	ServerHost             string
// 	LogLevel               string
// 	RateLimit              int
// 	Environment            string
// 	JWTSecretKey           string
// 	JWTExpirationTime      int
// 	RestAuthClasses        string
// 	RestPermissionClasses  string
// 	RestPaginationClass    string
// 	RestFilterBackends     string
// 	RestRendererClasses    string
// 	RestThrottleClasses    string
// 	RestThrottleRatesAnon  string
// 	RestThrottleRatesUser  string
// }

// func InitConfig() (*Config, error) {
// 	// Load .env file only if not running on Railway (or similar env)
// 	if _, exists := os.LookupEnv("RAILWAY_ENVIRONMENT"); !exists {
// 		err := godotenv.Load()
// 		if err != nil {
// 			log.Println("No .env file found, relying on system environment variables")
// 		}
// 	}

// 	// Initialize viper
// 	viper.AutomaticEnv() // Automatically bind to environment variables

// 	// Create and populate the Config struct
// 	config := &Config{
// 		SecretKey:              viper.GetString("SECRET_KEY"),
// 		AuthUserModel:          viper.GetString("AUTH_USER_MODEL"),
// 		CORSAllowedOrigins:     viper.GetString("CORS_ALLOWED_ORIGINS"),
// 		ClientId:               viper.GetString("CLIENT_ID"),
// 		ClientSecret:           viper.GetString("CLIENT_SECRET"),
// 		FrontendBaseUrl:        viper.GetString("FRONTEND_BASE_URL"),
// 		AccessTokenLifetime:    viper.GetInt("ACCESS_TOKEN_LIFETIME"),
// 		RefreshTokenLifetime:   viper.GetInt("REFRESH_TOKEN_LIFETIME"),
// 		RotateRefreshTokens:    viper.GetBool("ROTATE_REFRESH_TOKENS"),
// 		BlacklistAfterRotation: viper.GetBool("BLACKLIST_AFTER_ROTATION"),
// 		AuthHeaderTypes:        viper.GetString("AUTH_HEADER_TYPES"),
// 		EmailBackend:           viper.GetString("EMAIL_BACKEND"),
// 		EmailHost:              viper.GetString("EMAIL_HOST"),
// 		EmailPort:              viper.GetInt("EMAIL_PORT"),
// 		EmailUseTLS:            viper.GetBool("EMAIL_USE_TLS"),
// 		EmailHostUser:          viper.GetString("EMAIL_HOST_USER"),
// 		EmailHostPassword:      viper.GetString("EMAIL_HOST_PASSWORD"),
// 		DefaultFromEmail:       viper.GetString("DEFAULT_FROM_EMAIL"),
// 		StaticURL:              viper.GetString("STATIC_URL"),
// 		MediaURL:               viper.GetString("MEDIA_URL"),
// 		MediaRoot:              viper.GetString("MEDIA_ROOT"),
// 		StaticRoot:             viper.GetString("STATIC_ROOT"),
// 		DBType:                 viper.GetString("DB_TYPE"),
// 		DBServer:               viper.GetString("DB_SERVER"),
// 		DBPort:                 viper.GetInt("DB_PORT"),
// 		DBUser:                 viper.GetString("DB_USER"),
// 		DBPassword:             viper.GetString("DB_PASSWORD"),
// 		DBName:                 viper.GetString("DB_NAME"),
// 		ServerPort:             viper.GetInt("SERVER_PORT"),
// 		ServerHost:             viper.GetString("SERVER_HOST"),
// 		LogLevel:               viper.GetString("LOG_LEVEL"),
// 		RateLimit:              viper.GetInt("RATE_LIMIT"),
// 		Environment:            viper.GetString("ENVIRONMENT"),
// 		JWTSecretKey:           viper.GetString("JWT_SECRET_KEY"),
// 		JWTExpirationTime:      viper.GetInt("JWT_EXPIRATION_TIME"),
// 		RestAuthClasses:        viper.GetString("REST_FRAMEWORK_DEFAULT_AUTHENTICATION_CLASSES"),
// 		RestPermissionClasses:  viper.GetString("REST_FRAMEWORK_DEFAULT_PERMISSION_CLASSES"),
// 		RestPaginationClass:    viper.GetString("REST_FRAMEWORK_DEFAULT_PAGINATION_CLASS"),
// 		RestFilterBackends:     viper.GetString("REST_FRAMEWORK_DEFAULT_FILTER_BACKENDS"),
// 		RestRendererClasses:    viper.GetString("REST_FRAMEWORK_DEFAULT_RENDERER_CLASSES"),
// 		RestThrottleClasses:    viper.GetString("REST_FRAMEWORK_DEFAULT_THROTTLE_CLASSES"),
// 		RestThrottleRatesAnon:  viper.GetString("REST_FRAMEWORK_DEFAULT_THROTTLE_RATES_ANON"),
// 		RestThrottleRatesUser:  viper.GetString("REST_FRAMEWORK_DEFAULT_THROTTLE_RATES_USER"),
// 	}

// 	if config.Environment == "development" {
// 		log.Printf("Loaded configuration: %+v", config)
// 	}

// 	return config, nil
// }