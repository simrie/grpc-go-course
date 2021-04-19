package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	"grpc-go-course/calculator/calculatorpb"

	"google.golang.org/grpc"
)

// Error if dummy struct does not implement unimplemented[pkg]ServiceServer
type server struct {
	calculatorpb.UnimplementedCalculatorServiceServer
}

func (*server) Sum(ctx context.Context, req *calculatorpb.CalculatorRequest) (*calculatorpb.CalculatorResponse, error) {
	fmt.Printf("Sum function was invoked with %v\n", req)
	num1 := req.GetCalculator().GetNum_1()
	num2 := req.GetCalculator().GetNum_2()
	result := float32(num1 + num2)
	res := &calculatorpb.CalculatorResponse{
		Answer: result,
	}
	return res, nil
}

func (*server) FindPrimes(req *calculatorpb.FindPrimesRequest, stream calculatorpb.CalculatorService_FindPrimesServer) error {
	fmt.Printf("Find Primes server side streaming function was invoked to find prime factors of %v\n", req.GetNum_1())

	// algorithm to find primes
	resetTestFactor := int64(2)
	testFactor := resetTestFactor
	numToFactor := req.GetNum_1()
	for numToFactor >= testFactor {
		if numToFactor%testFactor == 0 {
			stream.Send(&calculatorpb.FindPrimesResponse{
				Prime: testFactor,
			})
			numToFactor = numToFactor / testFactor
			//fmt.Printf("\nnumToFactor is now %d \n", numToFactor)
			testFactor = resetTestFactor
		} else {
			testFactor++
		}
		//fmt.Printf("\ntestFactor is now to %v\n", testFactor)
	}
	return nil
}

func (*server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	fmt.Printf("Client stream function ComputerAverage was invoked\n")

	responses := make([]int64, 0)
	var average float64 = 0
	var numerator int64 = 0
	var denominator int64 = 0

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// we are done
			// computer the average
			average = float64(numerator) / float64(denominator)
			return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
				Average: average,
			})

		}
		if err != nil {
			log.Fatalf("\nerror reading client stream %v\n", err)
		}
		aNum := req.GetNum()
		numerator = numerator + aNum
		denominator++
		responses = append(responses, aNum)
	}
}

func (*server) Div(ctx context.Context, req *calculatorpb.CalculatorRequest) (*calculatorpb.CalculatorResponse, error) {
	fmt.Printf("Div function was invoked with %v\n", req)
	num1 := req.GetCalculator().GetNum_1()
	num2 := req.GetCalculator().GetNum_2()
	result := float32(num1 / num2)
	res := &calculatorpb.CalculatorResponse{
		Answer: result,
	}
	return res, nil
}

func (*server) GetHighestSoFar(stream calculatorpb.CalculatorService_GetHighestSoFarServer) error {
	fmt.Printf("GetHighestSoFar function was invoked as a bi-directional streaming request\n")

	var highestSoFar int32 = 0

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
			return err
		}
		testNum := req.GetNum()
		if testNum > highestSoFar {
			highestSoFar = testNum
			err = stream.Send(&calculatorpb.GetHighestIntResponse{
				Answer: highestSoFar,
			})
			if err != nil {
				//log.Fatalf("Error while sending data to client: %v", err)
				return err
			}
		}

	}
}

func main() {
	fmt.Println("It's time for Numberwang")

	lis, err := net.Listen("tcp", "0.0.0.0:50052")
	if err != nil {
		log.Fatalf("cannot listen to grpc port for tcp: %v", err)
	}

	s := grpc.NewServer()

	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve %v", err)
	}

}
