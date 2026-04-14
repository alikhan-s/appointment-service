package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alikhan-s/appointment-s/internal/client"
	"github.com/alikhan-s/appointment-s/internal/repository"
	transport "github.com/alikhan-s/appointment-s/internal/transport/grpc"
	"github.com/alikhan-s/appointment-s/internal/usecase"
	pb "github.com/alikhan-s/appointment-s/proto"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
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

	doctorServiceURL := "localhost:8081"
	conn, err := grpc.NewClient(doctorServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Did not connect to Doctor Service: %v", err)
	}
	defer conn.Close()

	docClient := client.NewDoctorGRPCClient(conn)
	repo := repository.NewAppointmentMongoRepo(db)
	usecaseLayer := usecase.NewAppointmentUseCase(repo, docClient)

	grpcServer := grpc.NewServer()
	handler := transport.NewAppointmentHandler(usecaseLayer)
	pb.RegisterAppointmentServiceServer(grpcServer, handler)

	listener, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// For production we need to delete gRPC Reflection for security (and urls migrate to .env file)
	reflection.Register(grpcServer)

	go func() {
		log.Println("Appointment gRPC Service is running on port 8082")
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("Listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down Appointment gRPC server...")

	grpcServer.GracefulStop()
	log.Println("Appointment Server exiting")
}
