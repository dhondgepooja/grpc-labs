package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"grpc/calculator/calculatorpb"
	"io"
	"log"
)

func main() {

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Error connecting: %v", err)
	}

	defer conn.Close()

	c := calculatorpb.NewCalculatorServiceClient(conn)

	//doUnary(c)
	//doServerStreaming(c)
	doClientStreaming(c)
}

func doUnary(c calculatorpb.CalculatorServiceClient) {

	req := calculatorpb.SumRequest{First: 10, Second: 15}
	res, err := c.Sum(context.Background(), &req)
	if err != nil {
		log.Fatalf("Error calling greet %v", err)
	}
	fmt.Println("Received response: ", res)
}

func doServerStreaming(c calculatorpb.CalculatorServiceClient) {
	req := calculatorpb.PrimeNumberDecompositionRequest{Number: 12}
	stream, err := c.PrimeNumberDecomposition(context.Background(), &req)
	if err != nil {
		log.Fatalf("Error calling greet %v", err)
	}
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			// we have reached the end of stream
			break
		}
		if err != nil {
			log.Fatalf("Error reading stream: %v", err)
		}
		fmt.Println("Received response: ", msg.GetPrimeFactor())
	}
}

func doClientStreaming(c calculatorpb.CalculatorServiceClient) {
	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("Error opening stream %v", err)
	}

	numbers := []int32{3, 5, 9, 54, 23}

	for _, number := range numbers {
		stream.Send(&calculatorpb.ComputeAverageRequest{
			Number: number,
		})
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error closing stream %v", err)
	}
	fmt.Println("Response:", res.Average)
}
