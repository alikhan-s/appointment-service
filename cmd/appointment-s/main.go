package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alikhan-s/appointment-s/internal/client"
	"github.com/alikhan-s/appointment-s/internal/repository"
	transport "github.com/alikhan-s/appointment-s/internal/transport/http"
	"github.com/alikhan-s/appointment-s/internal/usecase"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoClient.Disconnect(context.Background())

	db := mongoClient.Database("appointment_db")

	doctorServiceURL := "http://localhost:8081"

	docClient := client.NewDoctorHTTPClient(doctorServiceURL)
	repo := repository.NewAppointmentMongoRepo(db)
	usecaseLayer := usecase.NewAppointmentUseCase(repo, docClient)

	router := gin.Default()
	transport.NewAppointmentHandler(router, usecaseLayer)

	srv := &http.Server{
		Addr:    ":8082",
		Handler: router,
	}

	go func() {
		log.Println("Appointment Service is running on port 8082")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down Appointment server...")

	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	if err := srv.Shutdown(ctxShutdown); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Appointment Server exiting")
}
