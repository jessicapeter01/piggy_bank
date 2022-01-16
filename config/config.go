package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/spf13/viper"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/sqlite3"
	"github.com/gofiber/template/html"
)

type Config struct {
	*viper.Viper
	O            interface{}
	errorHandler fiber.ErrorHandler
	fiber        *fiber.Config
	database     *DatabaseConfig
}

var instantiated *Config
var once sync.Once

var defaultErrorHandler = func(c *fiber.Ctx, err error) error {
	// Status code defaults to 500
	code := fiber.StatusInternalServerError

	// Set error message
	message := err.Error()

	// Check if it's a fiber.Error type
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	// Return HTTP response
	c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
	c.Status(code)

	// Render default error view
	err = c.Render("errors/"+strconv.Itoa(code), fiber.Map{"message": message})
	if err != nil {
		return c.SendString(message)
	}
	return err
}

func GetInstance() *Config {
	once.Do(func() {
		instantiated = &Config{
			Viper: viper.New(),
		}

		// Set default configurations
		instantiated.setDefaults()

		// Select the .env file
		instantiated.SetConfigName(".env")
		instantiated.SetConfigType("dotenv")
		instantiated.AddConfigPath(".")

		// Automatically refresh environment variables
		instantiated.AutomaticEnv()

		// Read configuration
		if err := instantiated.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				fmt.Println("failed to read configuration:", err.Error())
				os.Exit(1)
			}
		}

		instantiated.SetErrorHandler(defaultErrorHandler)

		// Set Fiber configurations
		instantiated.setFiberConfig()

		instantiated.setDatabaseConfig()
	})
	return instantiated
}

func (config *Config) SetErrorHandler(errorHandler fiber.ErrorHandler) {
	config.errorHandler = errorHandler
}

func (config *Config) setDefaults() {
	// Set default App configuration
	config.SetDefault("APP_ADDR", ":3000")
	config.SetDefault("APP_ENV", "local")

	// Set default database configuration
	config.SetDefault("DB_DRIVER", "sqlite3")
	config.SetDefault("DB_HOST", "localhost")
	config.SetDefault("DB_USERNAME", "")
	config.SetDefault("DB_PASSWORD", "")
	config.SetDefault("DB_PORT", 5432)
	config.SetDefault("DB_NAME", "piggy_bank")

	// Set default session configuration
	config.SetDefault("MW_FIBER_SESSION_STORAGE_PROVIDER", "sqlite3")
	config.SetDefault("MW_FIBER_SESSION_STORAGE_HOST", "localhost")
	config.SetDefault("MW_FIBER_SESSION_STORAGE_PORT", 6379)
	config.SetDefault("MW_FIBER_SESSION_STORAGE_USERNAME", "")
	config.SetDefault("MW_FIBER_SESSION_STORAGE_PASSWORD", "")
	config.SetDefault("MW_FIBER_SESSION_STORAGE_DATABASE", "piggy_bank")
	config.SetDefault("MW_FIBER_SESSION_STORAGE_TABLE", "sessions")
	config.SetDefault("MW_FIBER_SESSION_COOKIENAME", "session_id")
	config.SetDefault("MW_FIBER_SESSION_COOKIEDOMAIN", "")
	config.SetDefault("MW_FIBER_SESSION_COOKIEPATH", "")
	config.SetDefault("MW_FIBER_SESSION_COOKIEHTTPONLY", false)
	config.SetDefault("MW_FIBER_SESSION_COOKIESAMESITE", "Lax")
	config.SetDefault("MW_FIBER_SESSION_STORAGE_RESET", false)
	config.SetDefault("MW_FIBER_SESSION_COOKIESECURE", true)
	config.SetDefault("MW_FIBER_SESSION_EXPIRATION", "24h")
	config.SetDefault("MW_FIBER_SESSION_STORAGE_GCINTERVAL", "10s")

	// Set default Fiber configuration
	config.SetDefault("FIBER_PREFORK", false)
	config.SetDefault("FIBER_SERVERHEADER", "")
	config.SetDefault("FIBER_STRICTROUTING", false)
	config.SetDefault("FIBER_CASESENSITIVE", false)
	config.SetDefault("FIBER_IMMUTABLE", false)
	config.SetDefault("FIBER_UNESCAPEPATH", false)
	config.SetDefault("FIBER_ETAG", false)
	config.SetDefault("FIBER_BODYLIMIT", 4194304)
	config.SetDefault("FIBER_CONCURRENCY", 262144)
	config.SetDefault("FIBER_VIEWS", "html")
	config.SetDefault("FIBER_VIEWS_DIRECTORY", "resources/views")
	config.SetDefault("FIBER_VIEWS_RELOAD", false)
	config.SetDefault("FIBER_VIEWS_DEBUG", false)
	config.SetDefault("FIBER_VIEWS_LAYOUT", "embed")
	config.SetDefault("FIBER_VIEWS_DELIMS_L", "{{")
	config.SetDefault("FIBER_VIEWS_DELIMS_R", "}}")
	config.SetDefault("FIBER_READTIMEOUT", 0)
	config.SetDefault("FIBER_WRITETIMEOUT", 0)
	config.SetDefault("FIBER_IDLETIMEOUT", 0)
	config.SetDefault("FIBER_READBUFFERSIZE", 4096)
	config.SetDefault("FIBER_WRITEBUFFERSIZE", 4096)
	config.SetDefault("FIBER_COMPRESSEDFILESUFFIX", ".fiber.gz")
	config.SetDefault("FIBER_PROXYHEADER", "")
	config.SetDefault("FIBER_GETONLY", false)
	config.SetDefault("FIBER_DISABLEKEEPALIVE", false)
	config.SetDefault("FIBER_DISABLEDEFAULTDATE", false)
	config.SetDefault("FIBER_DISABLEDEFAULTCONTENTTYPE", false)
	config.SetDefault("FIBER_DISABLEHEADERNORMALIZING", false)
	config.SetDefault("FIBER_DISABLESTARTUPMESSAGE", false)
	config.SetDefault("FIBER_REDUCEMEMORYUSAGE", false)
}

func (config *Config) setDatabaseConfig() {
	config.database = &DatabaseConfig{
		Default: DatabaseDriver{
			Driver:   config.GetString("DB_DRIVER"),
			Host:     config.GetString("DB_HOST"),
			Username: config.GetString("DB_USERNAME"),
			Password: config.GetString("DB_PASSWORD"),
			DBName:   config.GetString("DB_NAME"),
			Port:     config.GetInt("DB_PORT"),
		},
	}
}

func (config *Config) GetDatabaseConfig() *DatabaseConfig {
	return config.database
}

func (config *Config) setFiberConfig() {
	config.fiber = &fiber.Config{
		Prefork:                   config.GetBool("FIBER_PREFORK"),
		ServerHeader:              config.GetString("FIBER_SERVERHEADER"),
		StrictRouting:             config.GetBool("FIBER_STRICTROUTING"),
		CaseSensitive:             config.GetBool("FIBER_CASESENSITIVE"),
		Immutable:                 config.GetBool("FIBER_IMMUTABLE"),
		UnescapePath:              config.GetBool("FIBER_UNESCAPEPATH"),
		ETag:                      config.GetBool("FIBER_ETAG"),
		BodyLimit:                 config.GetInt("FIBER_BODYLIMIT"),
		Concurrency:               config.GetInt("FIBER_CONCURRENCY"),
		Views:                     config.getFiberViewsEngine(),
		ReadTimeout:               config.GetDuration("FIBER_READTIMEOUT"),
		WriteTimeout:              config.GetDuration("FIBER_WRITETIMEOUT"),
		IdleTimeout:               config.GetDuration("FIBER_IDLETIMEOUT"),
		ReadBufferSize:            config.GetInt("FIBER_READBUFFERSIZE"),
		WriteBufferSize:           config.GetInt("FIBER_WRITEBUFFERSIZE"),
		CompressedFileSuffix:      config.GetString("FIBER_COMPRESSEDFILESUFFIX"),
		ProxyHeader:               config.GetString("FIBER_PROXYHEADER"),
		GETOnly:                   config.GetBool("FIBER_GETONLY"),
		ErrorHandler:              config.errorHandler,
		DisableKeepalive:          config.GetBool("FIBER_DISABLEKEEPALIVE"),
		DisableDefaultDate:        config.GetBool("FIBER_DISABLEDEFAULTDATE"),
		DisableDefaultContentType: config.GetBool("FIBER_DISABLEDEFAULTCONTENTTYPE"),
		DisableHeaderNormalizing:  config.GetBool("FIBER_DISABLEHEADERNORMALIZING"),
		DisableStartupMessage:     config.GetBool("FIBER_DISABLESTARTUPMESSAGE"),
		ReduceMemoryUsage:         config.GetBool("FIBER_REDUCEMEMORYUSAGE"),
	}
}

func (config *Config) GetFiberConfig() *fiber.Config {
	return config.fiber
}

func (config *Config) GetSessionConfig() session.Config {
	store := sqlite3.New(sqlite3.Config{
		Database:   config.GetString("MW_FIBER_SESSION_STORAGE_DATABASE"),
		Table:      config.GetString("MW_FIBER_SESSION_STORAGE_TABLE"),
		Reset:      config.GetBool("MW_FIBER_SESSION_STORAGE_RESET"),
		GCInterval: config.GetDuration("MW_FIBER_SESSION_STORAGE_GCINTERVAL"),
	})
	return session.Config{
		Expiration:     config.GetDuration("MW_FIBER_SESSION_EXPIRATION"),
		Storage:        store,
		KeyLookup:      fmt.Sprintf("cookie:%s", config.GetString("MW_FIBER_SESSION_COOKIENAME")),
		CookieDomain:   config.GetString("MW_FIBER_SESSION_COOKIEDOMAIN"),
		CookiePath:     config.GetString("MW_FIBER_SESSION_COOKIEPATH"),
		CookieSecure:   config.GetBool("MW_FIBER_SESSION_COOKIESECURE"),
		CookieHTTPOnly: config.GetBool("MW_FIBER_SESSION_COOKIEHTTPONLY"),
		CookieSameSite: config.GetString("MW_FIBER_SESSION_COOKIESAMESITE"),
	}
}

func (config *Config) GetCorsConfig() cors.Config {
	return cors.Config{
		AllowOrigins:     config.GetString("MW_FIBER_CORS_ALLOWORIGINS"),
		AllowMethods:     config.GetString("MW_FIBER_CORS_ALLOWMETHODS"),
		AllowHeaders:     config.GetString("MW_FIBER_CORS_ALLOWHEADERS"),
		AllowCredentials: config.GetBool("MW_FIBER_CORS_ALLOWCREDENTIALS"),
		ExposeHeaders:    config.GetString("MW_FIBER_CORS_EXPOSEHEADERS"),
		MaxAge:           config.GetInt("MW_FIBER_CORS_MAXAGE"),
	}
}

func (config *Config) GetCsrfConfig() csrf.Config {
	store := sqlite3.New(sqlite3.Config{
		Database:   config.GetString("MW_FIBER_CSRF_STORAGE_DATABASE"),
		Table:      config.GetString("MW_FIBER_CSRF_STORAGE_TABLE"),
		Reset:      config.GetBool("MW_FIBER_CSRF_STORAGE_RESET"),
		GCInterval: config.GetDuration("MW_FIBER_CSRF_STORAGE_GCINTERVAL"),
	})
	return csrf.Config{
		KeyLookup:      config.GetString("MW_FIBER_CSRF_KEYLOOKUP"),
		Expiration:     config.GetDuration("MW_FIBER_CSRF_EXPIRATION"),
		Storage:        store,
		CookieName:     config.GetString("MW_FIBER_CSRF_COOKIENAME"),
		CookieDomain:   config.GetString("MW_FIBER_CSRF_COOKIEDOMAIN"),
		CookiePath:     config.GetString("MW_FIBER_CSRF_COOKIEPATH"),
		CookieSecure:   config.GetBool("MW_FIBER_CSRF_COOKIESECURE"),
		CookieHTTPOnly: config.GetBool("MW_FIBER_CSRF_COOKIEHTTPONLY"),
		CookieSameSite: config.GetString("MW_FIBER_CSRF_COOKIESAMESITE"),
	}
}

func (config *Config) getFiberViewsEngine() fiber.Views {
	// Set file extension dynamically to FIBER_VIEWS
	if config.GetString("FIBER_VIEWS_EXTENSION") == "" {
		config.Set("FIBER_VIEWS_EXTENSION", ".html")
	}
	viewsEngine := html.New(config.GetString("FIBER_VIEWS_DIRECTORY"), config.GetString("FIBER_VIEWS_EXTENSION"))
	viewsEngine.AddFunc(
		"ToLower", func(s string) string {
			return strings.ToLower(s)
		},
	)
	viewsEngine.Reload(config.GetBool("FIBER_VIEWS_RELOAD")).
		Debug(config.GetBool("FIBER_VIEWS_DEBUG")).
		Layout(config.GetString("FIBER_VIEWS_LAYOUT")).
		Delims(config.GetString("FIBER_VIEWS_DELIMS_L"), config.GetString("FIBER_VIEWS_DELIMS_R"))
	return viewsEngine
}
