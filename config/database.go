package config

import (
	"log"

	"github.com/Mitsui515/finsys/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

type DatabaseConfig struct {
	Type     string `json:"type"`
	Path     string `json:"path"`
	LogLevel string `json:"log_level"`
}

func InitDB() *gorm.DB {
	dbConfig := DefaultDBConfig()
	var err error
	logLevel := logger.Info
	switch dbConfig.LogLevel {
	case "silent":
		logLevel = logger.Silent
	case "error":
		logLevel = logger.Error
	case "warn":
		logLevel = logger.Warn
	}
	DB, err = gorm.Open(sqlite.Open(dbConfig.Path), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		log.Fatalf("cannot connect to database: %v", err)
	}
	err = DB.AutoMigrate(&model.Transaction{}, &model.User{}, &model.FraudReport{})
	if err != nil {
		log.Fatalf("fail to automatically migrate: %v", err)
		return nil
	}
	log.Println("successfully init database")
	return DB
}

func DefaultDBConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Type:     "sqlite",
		Path:     "./finsys.db",
		LogLevel: "info",
	}
}
