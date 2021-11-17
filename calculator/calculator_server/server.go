package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"grpc/calculator/calculatorpb"
	"io"
	"log"
	"net"
	"time"
)

type server struct{}

func (s server) ComputeAverage(averageServer calculatorpb.CalculatorService_ComputeAverageServer) error {
	fmt.Println("Received ComputeAverage RPC")

	sum := int32(0)
	count := 0
	for {
		req, err := averageServer.Recv()
		if err == io.EOF {
			average := float64(sum) / float64(count)
			res := calculatorpb.ComputeAverageResponse{
				Average: average,
			}
			return averageServer.SendAndClose(&res)
		}
		fmt.Println("received number:", req.GetNumber())
		sum += req.GetNumber()
		count++
	}
}

func (s server) PrimeNumberDecomposition(request *calculatorpb.PrimeNumberDecompositionRequest, decompositionServer calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	fmt.Printf("Received request: %v\n", request)
	number := request.GetNumber()
	divisor := int32(2)

	for number > 1 {
		if number%divisor == 0 {
			decompositionServer.Send(&calculatorpb.PrimeNumberDecompositionResponse{
				PrimeFactor: divisor,
			})
			number = number / divisor
		} else {
			divisor++
			fmt.Printf("Divisor has increated to %d", divisor)
		}
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}

func (s server) Sum(ctx context.Context, request *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	fmt.Printf("Received request: %v\n", request)
	first := request.First
	second := request.Second
	sum := first + second
	res := calculatorpb.SumResponse{
		Result: sum,
	}
	return &res, nil
}

func main() {
	fmt.Println("Welcome to calculator server")

	lis, err := net.Listen("tcp", "0.0.0.0:50051") //default port for grpc
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	// Register reflections
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to server: %v", err)
	}
}
