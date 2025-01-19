package main

import (
	"log"
	"net"
	"product-microservice/config"
	"product-microservice/db"
	"product-microservice/internal/domain"
	"product-microservice/internal/repository"
	"product-microservice/internal/service"
	grpcTransport "product-microservice/internal/transport/grpc"
	pb "product-microservice/proto/product"
	// sp "product-microservice/proto/subscription"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to the database
	database, err := db.ConnectDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Migrate models
	if err := migrateModels(database); err != nil {
		log.Fatalf("Database migration failed: %v", err)
	}

	// Initialize repositories and services
	productRepo := repository.ProductRepositoryImpl{DB: database}  // Ensure the repo is properly initialized
	// subscriptionRepo := repository.NewSubscriptionRepository(database) // Assuming this repo is initialized correctly

	productService := service.NewProductService(&productRepo)  // Initialize the service
	// subscriptionService := service.NewSubscriptionService(subscriptionRepo)
	//subscriptionService := service.NewSubscriptionService(subscriptionRepo)

	// Start gRPC server
	listener, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", cfg.GRPCPort, err)
	}

	server := grpc.NewServer()
	pb.RegisterProductServiceServer(server, grpcTransport.NewProductHandler(productService))
	// sp.RegisterSubscriptionServiceServer(server, grpcTransport.NewSubscriptionHandler(subscriptionService))

	log.Printf("gRPC server running on port %s", cfg.GRPCPort)
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}

func migrateModels(db *gorm.DB) error {
	log.Println("Starting database migration...")
	err := db.AutoMigrate(
		&domain.Product{},      
		&domain.SubscriptionPlan{}, 
	)
	if err == nil {
		log.Println("Database migrated successfully")
	}
	return err
}
