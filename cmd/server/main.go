package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/lazhari/blog-grpc/blogpb"
	"github.com/lazhari/blog-grpc/internal"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

func main() {
	// if we crash the go code, we get the file name and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fmt.Println("Blog Service gRPC")

	fmt.Println("Connecting to MongoDB")
	// Connect to MongoDB
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	collection := client.Database("blogdb").Collection("blog")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	tls := false
	opts := []grpc.ServerOption{}

	if tls {
		certFile := "ssl/server.crt"
		keyFile := "ssl/server.pem"

		creds, sslErr := credentials.NewServerTLSFromFile(certFile, keyFile)

		if sslErr != nil {
			log.Fatalf("Failed loanding certificates: %v", sslErr)
			return
		}

		opts = append(opts, grpc.Creds(creds))
	}

	s := grpc.NewServer(opts...)
	blogpb.RegisterBlogServiceServer(s, &internal.Server{
		Collection: collection,
	})
	// Register reflection
	reflection.Register(s)

	go func() {
		fmt.Println("Starting Server....")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Wait for control C to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	// Block until a signal is received
	<-ch
	fmt.Println("Stopping the server")
	s.Stop()
	fmt.Println("Closing the listener")
	lis.Close()
	fmt.Println("Closing Mongodb connection")
	client.Disconnect(context.TODO())
	fmt.Println("End of program")
}
