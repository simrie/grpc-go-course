package main

import (
	"context"
	"fmt"
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
