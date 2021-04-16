package main

import (
	"context"
	"fmt"
	"grpc-go-course/calculator/calculatorpb"
	"io"
	"log"

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

	doFindPrimesClientStreaming(client)
}

func doFindPrimesClientStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("...Calling Find Primes client streaming RPC...")

	testNumbers := []int64{10, 11, 21, 39, 60, 188, 231, 348}
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
