package main

import (
	"context"
	"fmt"
	"grpc-go-course/greet/greetpb"
	"log"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello Client")

	clientConnectionObject, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Dial error %v", err)
	}

	defer clientConnectionObject.Close()

	// The client generating line below worked to generate a client from the service
	// with greetpb.UnimplementedGreetServiceServer
	// before we added the (*server) Greet function definition to server.go

	// We create the client but we cannot do anything with it yet
	//client := greetpb.GreetServiceClient(clientConnectionObject)

	// Now that the service has a Greet service implemented we do:
	client := greetpb.NewGreetServiceClient(clientConnectionObject)

	fmt.Printf("Created client %f", client)
	doUnary(client)

}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("\n...Starting to do a Unary RPC...")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Boopsie",
			LastName:  "McFeathers",
		},
	}
	// context.Background() initializes a new, non-nil context
	// to be passed between server APIs
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("\nerror while calling Greet RPC: %v", err)
	}
	log.Printf("\nResponse from Greet: %v", res.Result)
}
