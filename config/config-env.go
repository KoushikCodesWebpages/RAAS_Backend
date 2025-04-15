package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
    "strings"
    "github.com/spf13/viper"
)

var Cfg *Config

type Config struct {
    SecretKey              string
    AuthUserModel          string
    CORSAllowedOrigins     string
    ClientId               string
    ClientSecret           string
    FrontendBaseUrl        string


    AuthHeaderTypes        string
    AccessTokenLifetime    int
    RefreshTokenLifetime   int
    RotateRefreshTokens    bool
    BlacklistAfterRotation bool



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

    AzureStorageAccount    string
    AzureStorageKey        string
    AzureBlobContainer     string

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

    // New fields for Hugging Face models and API key
    HFModelForMS1               string
    HFModelForMS2               string
    HFModelForMS3               string
    HFModelForMS4               string
    HFModelForMS5               string
    HFModelForMS6               string
    HFModelForMS7               string
    HFModelForMS8               string
    HFModelForMS9               string
    HFModelForMS10              string

    // CV MODELS
    HFModelForCL1               string
    HFModelForCL2               string
    HFModelForCL3               string
    HFModelForCL4               string
    HFModelForCL5               string
    HFModelForCL6               string
    HFModelForCL7               string
    HFModelForCL8               string
    HFModelForCL9               string
    HFModelForCL10              string

    HFBaseAPIUrl                string

    HFAPIKey                    string
    
    //COVERLETTER GENERATION
    CL_Url                      string
    CV_Url                      string

    GEN_API_KEY                 string
}   
func RemoveSystemEnv() {
    for _, pair := range os.Environ() {
        kv := strings.SplitN(pair, "=", 2)
        if len(kv) != 2 {
            continue
        }
        os.Unsetenv(kv[0])
    }
}

func InitConfig() error {

    //RemoveSystemEnv()
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
    Cfg = &Config{
        SecretKey:              viper.GetString("SECRET_KEY"),
        AuthUserModel:          viper.GetString("AUTH_USER_MODEL"),
        CORSAllowedOrigins:     viper.GetString("CORS_ALLOWED_ORIGINS"),
        ClientId:               viper.GetString("CLIENT_ID"),
        ClientSecret:           viper.GetString("CLIENT_SECRET"),
        FrontendBaseUrl:        viper.GetString("FRONTEND_BASE_URL"),

        JWTSecretKey:           viper.GetString("JWT_SECRET_KEY"),
        JWTExpirationTime:      viper.GetInt("JWT_EXPIRATION_TIME"),
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

        AzureStorageAccount: viper.GetString("AZURE_STORAGE_ACCOUNT"),
        AzureStorageKey:    viper.GetString("AZURE_STORAGE_KEY"),
        AzureBlobContainer: viper.GetString("AZURE_BLOB_CONTAINER"),

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

        RestAuthClasses:        viper.GetString("REST_FRAMEWORK_DEFAULT_AUTHENTICATION_CLASSES"),
        RestPermissionClasses:  viper.GetString("REST_FRAMEWORK_DEFAULT_PERMISSION_CLASSES"),
        RestPaginationClass:    viper.GetString("REST_FRAMEWORK_DEFAULT_PAGINATION_CLASS"),
        RestFilterBackends:     viper.GetString("REST_FRAMEWORK_DEFAULT_FILTER_BACKENDS"),
        RestRendererClasses:    viper.GetString("REST_FRAMEWORK_DEFAULT_RENDERER_CLASSES"),
        RestThrottleClasses:    viper.GetString("REST_FRAMEWORK_DEFAULT_THROTTLE_CLASSES"),
        RestThrottleRatesAnon:  viper.GetString("REST_FRAMEWORK_DEFAULT_THROTTLE_RATES_ANON"),
        RestThrottleRatesUser:  viper.GetString("REST_FRAMEWORK_DEFAULT_THROTTLE_RATES_USER"),

        // Load Hugging Face model and API key values from environment variables
        HFModelForMS1:               viper.GetString("HF_MODEL_FOR_MS_1"),
        HFModelForMS2:               viper.GetString("HF_MODEL_FOR_MS_2"),
        HFModelForMS3:               viper.GetString("HF_MODEL_FOR_MS_3"),
        HFModelForMS4:               viper.GetString("HF_MODEL_FOR_MS_4"),
        HFModelForMS5:               viper.GetString("HF_MODEL_FOR_MS_5"),
        HFModelForMS6:               viper.GetString("HF_MODEL_FOR_MS_6"),
        HFModelForMS7:               viper.GetString("HF_MODEL_FOR_MS_7"),
        HFModelForMS8:               viper.GetString("HF_MODEL_FOR_MS_8"),
        HFModelForMS9:               viper.GetString("HF_MODEL_FOR_MS_9"),
        HFModelForMS10:              viper.GetString("HF_MODEL_FOR_MS_10"),
       
        //CV GEN

        HFModelForCL1: viper.GetString("HF_MODEL_FOR_CL_1"),
        HFModelForCL2: viper.GetString("HF_MODEL_FOR_CL_2"),
        HFModelForCL3: viper.GetString("HF_MODEL_FOR_CL_3"),
        HFModelForCL4: viper.GetString("HF_MODEL_FOR_CL_4"),
        HFModelForCL5: viper.GetString("HF_MODEL_FOR_CL_5"),
        HFModelForCL6: viper.GetString("HF_MODEL_FOR_CL_6"),
        HFModelForCL7: viper.GetString("HF_MODEL_FOR_CL_7"),
        HFModelForCL8: viper.GetString("HF_MODEL_FOR_CL_8"),
        HFModelForCL9: viper.GetString("HF_MODEL_FOR_CL_9"),
        HFModelForCL10: viper.GetString("HF_MODEL_FOR_CL_10"),


        HFAPIKey:               viper.GetString("HF_API_KEY"),

        HFBaseAPIUrl:     viper.GetString("HF_BASE_API_URL"),

        //GENERATION

        CL_Url: viper.GetString("COVER_LETTER_API_URL"),
        CV_Url: viper.GetString("CV_RESUME_API_URL"), 
        
        GEN_API_KEY: viper.GetString("COVER_CV_API_KEY"),
    }

    return nil
}



