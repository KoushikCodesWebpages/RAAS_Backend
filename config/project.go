package config

import "github.com/spf13/viper"

type ProjectConfig struct {


	AuthUserModel               	string
    FrontendBaseUrl             	string

	CORSAllowedOrigins        		string
	AuthHeaderTypes             	string


	
    JWTSecretKey                    string
    JWTExpirationTime               int
    AccessTokenLifetime         	int
    RefreshTokenLifetime        	int
    RotateRefreshTokens         	bool
    BlacklistAfterRotation      	bool


	SecretKey                 		string
	StaticURL						string
	MediaURL                    	string
    MediaRoot                   	string
    StaticRoot                  	string


	RestAuthClasses                 string
    RestPermissionClasses           string
    RestPaginationClass             string
    RestFilterBackends              string
    RestRendererClasses             string
    RestThrottleClasses             string
    RestThrottleRatesAnon           string
    RestThrottleRatesUser           string


	


}

func LoadProjectConfig() (*ProjectConfig, error) {
    ProjectConfig := &ProjectConfig{

		StaticURL:              viper.GetString("STATIC_URL"),
        MediaURL:               viper.GetString("MEDIA_URL"),
        MediaRoot:              viper.GetString("MEDIA_ROOT"),
        StaticRoot:             viper.GetString("STATIC_ROOT"),
		RestAuthClasses:        viper.GetString("REST_FRAMEWORK_DEFAULT_AUTHENTICATION_CLASSES"),
        RestPermissionClasses:  viper.GetString("REST_FRAMEWORK_DEFAULT_PERMISSION_CLASSES"),
        RestPaginationClass:    viper.GetString("REST_FRAMEWORK_DEFAULT_PAGINATION_CLASS"),
        RestFilterBackends:     viper.GetString("REST_FRAMEWORK_DEFAULT_FILTER_BACKENDS"),
        RestRendererClasses:    viper.GetString("REST_FRAMEWORK_DEFAULT_RENDERER_CLASSES"),
        RestThrottleClasses:    viper.GetString("REST_FRAMEWORK_DEFAULT_THROTTLE_CLASSES"),
        RestThrottleRatesAnon:  viper.GetString("REST_FRAMEWORK_DEFAULT_THROTTLE_RATES_ANON"),
        RestThrottleRatesUser:  viper.GetString("REST_FRAMEWORK_DEFAULT_THROTTLE_RATES_USER"),
		


    }

    return ProjectConfig, nil
}
