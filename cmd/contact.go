package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"grpc-contact-manager/services/middlewares"
	"grpc-contact-manager/services/servers"

	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.TraceLevel)
	ctx := context.Background()

	if os.Getenv("ENVIRONMENT") == "" || os.Getenv("ENVIRONMENT") == "development" {
		err := godotenv.Load("./.envs/.env")
		if err != nil {
			panic(err)
		}
	}
	host := os.Getenv("HOST")
	userName := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")
	dbName := os.Getenv("DB_NAME")
	grpcPort := os.Getenv("USER_PORT")
	port := os.Getenv("PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable", host, userName, password, dbName)
	log.Infof("DSN: %s", dsn)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})

	if err != nil {
		panic(err)
	}
	log.Info("DB Connected successfully")

	// Register the prometheus metrics
	middlewares.RegisterPrometheusMetrics()

	server, err := servers.New(db)
	if err != nil {
		panic(err)
	}

	server.Router.Use(middlewares.RecordRequestLatency())
	server.UserRoutes() //setup the user routes
	httpServer, err := server.StartHttp(ctx, port)
	if err != nil {
		panic(err)
	}

	go func() {
		log.Infof("Start HTTP Server on port: %s", port)
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(":3501", nil); err != nil {
			log.Fatal(err)
		}
	}()

	if grpcPort == "" {
		grpcPort = ":5200"
	}
	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		panic(err)
	}
	userGRPC, err := server.StartUserGRPC(ctx)
	if err != nil {
		panic(err)
	}

	go func() {
		log.Infof("Start User GRPC Server on port: %s", grpcPort)
		if err := userGRPC.Serve(lis); err != nil {
			panic(err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server")
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	userGRPC.GracefulStop()
	log.Info("Server stopped successfully")
}
