package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"bronirovanie"
	"bronirovanie/configs"
	"bronirovanie/pkg/handler"
	"bronirovanie/pkg/repository"
	"bronirovanie/pkg/service"

	_ "bronirovanie/docs"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// @title Reservation App
// @description API server for Reservation application
// @host localhost:8443
// @BasePath /
func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := configs.InitConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(viper.GetString("environment.DATABASE_URL"))
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	err = repository.CreateTable(db)
	if err != nil {
		logrus.Fatal(err)
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(bronirovanie.Server)
	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			logrus.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()
	logrus.Print(viper.GetString("port"))

	logrus.Print("Bronirovanie app started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("Bronirovanie app shutting down")

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}

}
