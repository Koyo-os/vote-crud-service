package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Koyo-os/vote-crud-service/internal/config"
	"github.com/Koyo-os/vote-crud-service/internal/repository"
	"github.com/Koyo-os/vote-crud-service/internal/server"
	"github.com/Koyo-os/vote-crud-service/internal/service"
	"github.com/Koyo-os/vote-crud-service/pkg/api/protobuf"
	"github.com/Koyo-os/vote-crud-service/pkg/logger"
	"github.com/Koyo-os/vote-crud-service/pkg/retrier"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func getService(logger *logger.Logger) (*service.Service, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	logger.Info("connecting to mariadb...", zap.String("dsn", dsn))

	db, err := retrier.Connect(10, 10, func() (*gorm.DB, error) {
		return gorm.Open(mysql.Open(dsn))
	})
	if err != nil{
		return nil, err
	}

	repository := repository.NewRepository(db, logger)

	service := service.NewService(repository)

	return service, nil
}

func main() {
	logcfg := logger.Config{
		AppName: "vote-crud-service",
		LogLevel: "debug",
		LogFile: "app.log",
		AddCaller: false,
	}

	if err := logger.Init(logcfg);err != nil{
		log.Fatal(err)
		
		return
	}

	logger := logger.Get()

	cfg := config.NewConfig()

	logger.Info("vote-crud-service starting...")

	grpcServer := grpc.NewServer()

	service, err := getService(logger)
	if err != nil{
		logger.Error("failed get service", zap.Error(err))
		
		return
	}

	server := server.NewServer(service, logger)

	protobuf.RegisterVoteServiceServer(grpcServer, server)
}