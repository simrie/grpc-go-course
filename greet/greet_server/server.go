package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"grpc-go-course/greet/greetpb"

	"google.golang.org/grpc"
)

// This code was created following along the Udemy grpc course

// To begin with create a server type to which we will add services
// however this may be replaced later in the course

// Error if dummy struct does not implement unimplementedGreetServiceServer
type server struct {
	greetpb.UnimplementedGreetServiceServer
}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet function was invoked with %v\n", req)
	firstName := req.GetGreeting().GetFirstName()
	result := "How ya doin', " + firstName
	res := &greetpb.GreetResponse{
		Result: result,
	}
	return res, nil
}

func main() {
	fmt.Println("Hello World")

	// Here we test the grpc code generated from greet.proto

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("cannot listen to grpc port for tcp: %v", err)
	}

	s := grpc.NewServer()

	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve %v", err)
	}

}
