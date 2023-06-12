package main

import (
	"log"
	"os"

	"github.com/fekuna/go-rest-clean-architecture/config"
	"github.com/fekuna/go-rest-clean-architecture/internal/server"
	"github.com/fekuna/go-rest-clean-architecture/pkg/db/minioS3"
	"github.com/fekuna/go-rest-clean-architecture/pkg/db/postgres"
	"github.com/fekuna/go-rest-clean-architecture/pkg/logger"
	"github.com/fekuna/go-rest-clean-architecture/pkg/utils"
)

// @title Go Example REST API
// @version 1.0
// @description Example Golang REST API
// @contact_name Alfan Almunawar
// @contact_url https://github.com/fekuna
// @contact_email almunawar.alfan@gmail.com
// @BasePath /api/v1
func main() {
	log.Println("Starting api server")

	configPath := utils.GetConfigPath(os.Getenv("config"))

	cfgFile, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("LoadConfig: %v", err)
	}

	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatalf("ParseConfig: %v", err)
	}

	appLogger := logger.NewApiLogger(cfg)

	appLogger.InitLogger()
	appLogger.Infof("AppVersion: %s, LogLeve: %s, Mode: %s, SSL: %v", cfg.Server.AppVersion, cfg.Logger.Level, cfg.Server.Mode, cfg.Server.SSL)

	// psqlDB, err
	psqlDB, err := postgres.NewPsqlDB(cfg)
	if err != nil {
		appLogger.Fatalf("Postgresql init: %s", err)
	} else {
		appLogger.Infof("Postgresql connected, Status: %#v", psqlDB.Stats())
	}
	defer psqlDB.Close()

	MinioConfig, err := minioS3.NewMinioS3Client(cfg.Minio.Endpoint, cfg.Minio.MinioAccessKey, cfg.Minio.MinioSecretKey, cfg.Minio.UseSSL)

	s := server.NewServer(cfg, appLogger, psqlDB, MinioConfig)
	if err = s.Run(); err != nil {
		log.Fatal(err)
	}
}
