package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	_ "github.com/go-sql-driver/mysql"
	infraDatabase "github.com/lucas/clean-orders-challenge/internal/infra/database"
	infraGrpc "github.com/lucas/clean-orders-challenge/internal/infra/grpc"
	infraGraphql "github.com/lucas/clean-orders-challenge/internal/infra/graphql"
	infraWeb "github.com/lucas/clean-orders-challenge/internal/infra/web"
	"github.com/lucas/clean-orders-challenge/internal/usecase"
	pb "github.com/lucas/clean-orders-challenge/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	dsn := getEnv("DB_DSN", "root:root@tcp(mysql:3306)/orders?parseTime=true")
	db, err := waitForDatabase(dsn, 60*time.Second)
	if err != nil {
		log.Fatalf("database unavailable: %v", err)
	}
	defer db.Close()

	if err := runMigrations(db, "migrations"); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	orderRepository := infraDatabase.NewOrderRepository(db)
	createOrderUseCase := usecase.NewCreateOrderUseCase(orderRepository)
	listOrdersUseCase := usecase.NewListOrdersUseCase(orderRepository)

	go startRESTServer(createOrderUseCase, listOrdersUseCase)
	go startGraphQLServer(createOrderUseCase, listOrdersUseCase)
	startGRPCServer(listOrdersUseCase)
}

func waitForDatabase(dsn string, timeout time.Duration) (*sql.DB, error) {
	deadline := time.Now().Add(timeout)
	var lastErr error

	for time.Now().Before(deadline) {
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			lastErr = err
			time.Sleep(time.Second)
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		err = db.PingContext(ctx)
		cancel()
		if err == nil {
			log.Println("database connection established")
			return db, nil
		}

		_ = db.Close()
		lastErr = err
		log.Printf("waiting for database: %v", err)
		time.Sleep(2 * time.Second)
	}

	if lastErr == nil {
		lastErr = errors.New("timeout waiting for database")
	}
	return nil, lastErr
}

func runMigrations(db *sql.DB, directory string) error {
	files, err := filepath.Glob(filepath.Join(directory, "*.sql"))
	if err != nil {
		return err
	}

	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		if _, err := db.Exec(string(content)); err != nil {
			return fmt.Errorf("%s: %w", file, err)
		}
		log.Printf("migration applied: %s", file)
	}

	return nil
}

func startRESTServer(createUC *usecase.CreateOrderUseCase, listUC *usecase.ListOrdersUseCase) {
	port := getEnv("WEB_PORT", "8000")
	mux := http.NewServeMux()
	handler := infraWeb.NewOrderHandler(createUC, listUC)

	mux.HandleFunc("/order", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handler.Create(w, r)
		case http.MethodGet:
			handler.List(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Printf("REST server running on :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("REST server failed: %v", err)
	}
}

func startGraphQLServer(createUC *usecase.CreateOrderUseCase, listUC *usecase.ListOrdersUseCase) {
	port := getEnv("GRAPHQL_PORT", "8080")
	mux := http.NewServeMux()
	mux.Handle("/query", infraGraphql.NewHandler(createUC, listUC))

	log.Printf("GraphQL server running on :%s/query", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("GraphQL server failed: %v", err)
	}
}

func startGRPCServer(listUC *usecase.ListOrdersUseCase) {
	port := getEnv("GRPC_PORT", "50051")
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen on gRPC port %s: %v", port, err)
	}

	server := grpc.NewServer()
	pb.RegisterOrderServiceServer(server, infraGrpc.NewOrderService(listUC))
	reflection.Register(server)

	log.Printf("gRPC server running on :%s", port)
	if err := server.Serve(listener); err != nil {
		log.Fatalf("gRPC server failed: %v", err)
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
