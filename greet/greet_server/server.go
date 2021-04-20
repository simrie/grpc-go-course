package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"grpc-go-course/greet/greetpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// This code was created following along the Udemy grpc course

// To begin with create a server type to which we will add services
// however this may be replaced later in the course

// Error if dummy struct does not implement unimplementedGreetServiceServer
type server struct {
	greetpb.UnimplementedGreetServiceServer
}

func (*server) LongRequestGreet(stream greetpb.GreetService_LongRequestGreetServer) error {
	// signature here was copied from under greet_grpc.pb.go's
	// type GreetServiceServer interface
	fmt.Printf("LongRequestGreet function was invoked with a stream\n")
	// client stream can theoretically return whenever it wants
	// but we are going to try to wait for the end of the requests
	// by doing stream.Recv() a bunch of times then stream.SendAndClose()
	responses := make([]string, 0)
	for {
		// the req comes from stream.recv
		req, err := stream.Recv()
		if err == io.EOF {
			// Client should send stream.CloseAndRecv()
			// so we know we have finished reading the client stream
			// SendAndClose returns an error so if this doesn't work
			// the error will be returns on the steam
			return stream.SendAndClose(&greetpb.GreetResponse{
				Result: fmt.Sprintf("%v", responses),
			})
		}
		if err != nil {
			log.Fatalf("\nerror reading client stream %v\n", err)
		}
		aGreeting := fmt.Sprintf("Hiya %s;", req.GetGreeting().GetFirstName())
		responses = append(responses, aGreeting)
	}
}

func (*server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	fmt.Printf("GreetEveryone function was invoked as a bi-directional streaming request\n")

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
			return err
		}
		firstName := req.GetGreeting().GetFirstName()
		result := "Hello " + firstName + "! "

		err = stream.Send(&greetpb.GreetEveryoneResponse{
			Result: result,
		})
		if err != nil {
			//log.Fatalf("Error while sending data to client: %v", err)
			return err
		}
	}

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

func (*server) GreetWithDeadline(ctx context.Context, req *greetpb.GreetWithDeadlineRequest) (*greetpb.GreetWithDeadlineResponse, error) {
	fmt.Printf("GreetWithDeadline function was invoked with %v\n", req)

	// Check to see if the timeout occurred
	for i := 0; i < 3; i++ {
		if ctx.Err() == context.Canceled {
			msg := "Client cancelled the request (timeout)"
			fmt.Println(msg)
			return nil, status.Error(codes.DeadlineExceeded, msg)
		}
		time.Sleep(1 * time.Second)
	}
	firstName := req.GetGreeting().GetFirstName()
	result := "Glad you made it, " + firstName
	res := &greetpb.GreetWithDeadlineResponse{
		Result: result,
	}
	return res, nil

}

func main() {
	fmt.Println("Ohayoo-san! Robo-greeta de gozaimasu.")

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
