package main

import (
	"context"
	"fmt"
	"grpc-go-course/greet/greetpb"
	"log"
	"time"

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

	doClientStreaming(client)

}

func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("\n...Starting to do a Client Streaming RPC...")

	// We don't need to pass a request, just a context
	// the streaming will be part of the context
	// It will return a client stream
	stream, err := c.LongRequestGreet(context.Background())
	if err != nil {
		log.Fatalf("Client streaming LongGreet error %v\n", err)
	}
	// We can do stream.Send() on the client stream as much as we want
	// and when done, we do stream.CloseAndRecv() which is
	// probably what the server is waiting for when it checks io.EOF

	nameList := []string{"Suzie", "ButtonEyes", "Wanda", "Ignatio", "Bugsy"}

	for _, s := range nameList {
		req := &greetpb.GreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: fmt.Sprintf("%s", s),
			},
		}
		fmt.Printf("\nSending request to LongGreet %s\n", s)
		stream.Send(req)
		time.Sleep(100 * time.Millisecond)
	}

	response, err := stream.CloseAndRecv()
	if err != nil {
		fmt.Printf("\nError receiving client stream LongGreet response. %v\n", err)
	}
	// In client streaming the client gets a single response
	// which was in this case defined as a string
	fmt.Printf("\n%s\n", response)
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
