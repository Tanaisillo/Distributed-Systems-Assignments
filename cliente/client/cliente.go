package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	pb "cliente/proto"

	"google.golang.org/grpc"
)

func main() {
	// Read JSON data from a file
	jsonData, err := os.ReadFile("products.json") //json name
	if err != nil {
		log.Fatalf("failed to read JSON file: %v", err)
	}

	// Unmarshal JSON data into the BookstoreRequest message
	var request pb.BookstoreRequest
	if err := json.Unmarshal(jsonData, &request); err != nil {
		log.Fatalf("failed to unmarshal JSON: %v", err)
	}

	// Connect to the gRPC server
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	// Create a gRPC client
	client := pb.NewBookstoreServiceClient(conn)

	// Send the request to the server
	response, err := client.ProcessOrder(context.Background(), &request)
	if err != nil {
		log.Fatalf("failed to call ProcessOrder: %v", err)
	}

	// Print the server's response
	fmt.Printf("El ID de Orden es: %s\n", response.Message)
}
