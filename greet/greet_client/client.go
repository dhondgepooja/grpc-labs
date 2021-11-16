package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"grpc/greet/greetpb"
	"io"
	"log"
	"time"
)

func main() {

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Error connecting: %v", err)
	}

	defer conn.Close()

	c := greetpb.NewGreetServiceClient(conn)

	//doUnary(c)
	//doServerStreaming(c)
	//doClientStreaming(c)
	doBiDirectionalStreaming(c)
}

func doUnary(c greetpb.GreetServiceClient) {
	greeting := greetpb.Greeting{FirstName: "Pooja", LastName: "Dhondge"}
	req := greetpb.GreetRequest{Greeting: &greeting}
	res, err := c.Greet(context.Background(), &req)
	if err != nil {
		log.Fatalf("Error calling greet %v", err)
	}
	fmt.Println("Received response: ", res)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	greeting := greetpb.Greeting{FirstName: "Pooja", LastName: "Dhondge"}
	req := greetpb.GreetManyTimesRequest{Greeting: &greeting}
	stream, err := c.GreetManyTimes(context.Background(), &req)
	if err != nil {
		log.Fatalf("Error calling GreetManyTimes %v", err)
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
		fmt.Println("Received response: ", msg.GetResult())
	}
}

func doClientStreaming(c greetpb.GreetServiceClient) {
	greetings := []greetpb.Greeting{
		{FirstName: "Pooja", LastName: "Dhondge"},
		{FirstName: "Pratik", LastName: "Pandit"},
		{FirstName: "Sarvesh", LastName: "Pandit"},
	}

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Error calling LongGreet %v", err)
	}
	for _, greeting := range greetings {
		req := greetpb.LongGreetRequest{Greeting: &greeting}
		fmt.Println("Sending: ", req.Greeting)
		stream.Send(&req)
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Encountered error closing stream %v", err)
	}
	fmt.Println("Response: ", res)
}

func doBiDirectionalStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Client for biDirectionalStreaming")
	greetings := []greetpb.Greeting{
		{FirstName: "Pooja", LastName: "Dhondge"},
		{FirstName: "Pratik", LastName: "Pandit"},
		{FirstName: "Sarvesh", LastName: "Pandit"},
	}

	// we create a stream by involking client
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Encountered error creating stream %v", err)
	}

	waitc := make(chan struct{})
	// we send a bunch of messages to the client ( go routine)
	go func() {
		// function to send a bunch of messages
		for _, greeting := range greetings {
			stream.Send(&greetpb.GreetEveryoneRequest{
				Greeting: &greeting,
			})
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()
	// we receive a bunch of messages from client (go routine)
	go func() {
		// function to receive a bunch of messages
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error reading stream %v", err)
				break
			}
			fmt.Println("Response:", res.GetResult())
		}
		close(waitc)
	}()
	// block until done
	<-waitc
}
