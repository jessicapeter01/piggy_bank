package database

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/jessicapeter01/piggy_bank/app/models"
	"github.com/jessicapeter01/piggy_bank/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

var (
	// DBConn for gorm
	DBConn *gorm.DB
	// SessionStore for session storage
	SessionStore *session.Store
	cfg          *config.DatabaseConfig
)

type Database struct {
	*gorm.DB
}

// InitGormDb with gorm models
func Setup() {
	if DBConn != nil {
		return
	}

	cfg = config.GetInstance().GetDatabaseConfig()

	if err := connect(); err != nil {
		log.Panicf("error could not connect database (%s)", err.Error())
	}

	log.Println("Connected to DB")

	if err := DBConn.AutoMigrate(&models.User{},
		&models.Goal{},
		&models.Transaction{}); err != nil {
		log.Panicf("Failed to automigrate %s", err.Error())
	}

	log.Println("Ran Auto Migrate")
}

func connect() error {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,         // Disable color
		},
	)
	var err error
	DBConn, err = gorm.Open(sqlite.Open(cfg.Default.DBName), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   newLogger,
	})
	if err = DBConn.Use(
		dbresolver.Register(dbresolver.Config{}).
			SetConnMaxIdleTime(time.Hour).
			SetConnMaxLifetime(24 * time.Hour).
			SetMaxIdleConns(100).
			SetMaxOpenConns(200),
	); err != nil {
		log.Fatalf("error: registering resolver failed '%s'", err.Error())
	}
	return nil
}
