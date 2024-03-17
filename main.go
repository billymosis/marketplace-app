package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/billymosis/marketplace-app/db"
	"github.com/billymosis/marketplace-app/handler/api"
	as "github.com/billymosis/marketplace-app/store/account"
	ps "github.com/billymosis/marketplace-app/store/product"
	us "github.com/billymosis/marketplace-app/store/user"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {

	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	host := os.Getenv("DB_HOST")
	database := os.Getenv("DB_NAME")
	port, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	user := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	maxOpenConnections := 100

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion("ap-southeast-1"),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				os.Getenv("S3_ID"), os.Getenv("S3_SECRET_KEY"), "",
			)))

	s3Client := s3.NewFromConfig(cfg)

	db, err := db.Connection("postgres", host, database, user, password, port, maxOpenConnections)
	if err != nil {
		log.Fatal(err)
	}
	validate := validator.New()

	userStore := us.NewUserStore(db, validate)
	productStore := ps.NewProductStore(db, validate)
	accountStore := as.NewAccountStore(db, validate)

	r := api.New(userStore, productStore, accountStore, s3Client)
	h := r.Handler()

	logrus.Info("application starting")

	log.Println("application starting")

	go func() {
		s := http.Server{
			Addr:           ":8000",
			Handler:        h,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20, //1mb
		}

		err := s.ListenAndServe()
		if err != nil {
			log.Println("application failed to start")
			panic(err)
		}
	}()
	log.Println("application started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Info("application shutting down")

	log.Println("database closing")
	db.Close()
	log.Println("database closed")
}
