package main

import (
	"context"
	"fmt"
	"grpc-go-course/calculator/calculatorpb"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello Numberwanger")

	clientConnectionObject, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Dial error %v", err)
	}

	defer clientConnectionObject.Close()

	client := calculatorpb.NewCalculatorServiceClient(clientConnectionObject)

	fmt.Println("Created client")
	//doSumUnary(client)
	//doDivUnary(client)

	doFindPrimesServerStreaming(client)
	doCalculateAverageClientStreaming(client)
	doGetHighestSoFarBiDiStreaming(client)
}

func doGetHighestSoFarBiDiStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("...Calling GetHighestSoFar bidirectional streaming RPC...")

	// we create a stream by invoking the client
	stream, err := c.GetHighestSoFar(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream: %v", err)
		return
	}

	testNumbers := []int32{1, 5, 3, 6, 2, 20}

	// Create a channel
	waitc := make(chan struct{})

	// we send a bunch of messages to the client (go routine)
	go func() {
		// function to send a bunch of messages
		for _, testNum := range testNumbers {
			req := &calculatorpb.GetHighestIntRequest{
				Num: testNum,
			}
			fmt.Printf("Sending message: %v\n", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	// we receive a bunch of messages from the client (go routine)
	go func() {
		// function to receive a bunch of messages
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving: %v", err)
				break
			}
			fmt.Printf("Received: %v\n", res.GetAnswer())
		}
		close(waitc)
	}()

	// block until everything is done
	<-waitc
}

func doCalculateAverageClientStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("...Calling CalculateAverage client streaming RPC...")

	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("Client streaming ComputeAverage error %v\n", err)
	}

	numsToAverage := []int64{3, 10, 42}

	for _, num := range numsToAverage {
		req := &calculatorpb.ComputeAverageRequest{
			Num: num,
		}
		fmt.Println("Sending request to CalculateAverage")
		stream.Send(req)
		time.Sleep(100 * time.Millisecond)
	}
	response, err := stream.CloseAndRecv()
	if err != nil {
		fmt.Printf("\nError receiving client stream CalculateAverage response. %v\n", err)
	}
	// In client streaming the client gets a single response
	// which was in this case defined as a float64
	fmt.Printf("\nThe average of %v is %f\n", numsToAverage, response.GetAverage())
}

func doFindPrimesServerStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("...Calling Find Primes server streaming RPC...")

	testNumbers := []int64{10, 11, 21, 39, 60, 188, 231, 348, 56789, 109, 521, 321654987, 419, 85297}
	for _, testNumber := range testNumbers {

		req := &calculatorpb.FindPrimesRequest{
			Num_1: testNumber,
		}
		primeFactors := make([]int64, 0)
		stream, err := c.FindPrimes(context.Background(), req)
		if err != nil {
			log.Fatalf("error while calling FindPrimes RPC: %v", err)
		}
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Something happened: %v", err)
			}
			primeFactors = append(primeFactors, res.GetPrime())
		}
		if len(primeFactors) == 1 {
			fmt.Printf("%d is prime \n", req.GetNum_1())
		} else {
			fmt.Printf("Prime factors of %d are:  %v\n", req.GetNum_1(), primeFactors)
		}

	}
}

func doSumUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("\n...Calling Sum Unary RPC...")
	req := &calculatorpb.CalculatorRequest{
		Calculator: &calculatorpb.Calculator{
			Num_1: 24,
			Num_2: 48,
		},
	}

	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("\nerror while calling Sum RPC: %v", err)
	}
	log.Printf("\nThe sum of %d and %d is: %v", req.Calculator.Num_1, req.Calculator.Num_2, res.Answer)
}

func doDivUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("\n...Calling Div Unary RPC...")
	req := &calculatorpb.CalculatorRequest{
		Calculator: &calculatorpb.Calculator{
			Num_1: 48,
			Num_2: 24,
		},
	}

	res, err := c.Div(context.Background(), req)
	if err != nil {
		log.Fatalf("\nerror while calling Div RPC: %v", err)
	}
	log.Printf("\n%d divided by %d is: %v", req.Calculator.Num_1, req.Calculator.Num_2, res.Answer)
}
