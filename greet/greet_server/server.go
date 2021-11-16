package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"grpc/greet/greetpb"
	"io"
	"log"
	"net"
	"strconv"
	"time"
)

type server struct{}

func (s server) GreetEveryone(everyoneServer greetpb.GreetService_GreetEveryoneServer) error {
	// Server decides when to end so it doesn't have to wait for all messages from client before closing
	// we will keep it simple here though

	fmt.Println("GreetEveryone function was invoked")
	for {
		req, err := everyoneServer.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error reading stream %v", err)
			return err
		}
		firstName := req.GetGreeting().GetFirstName()
		result := "Hello " + firstName
		err = everyoneServer.Send(&greetpb.GreetEveryoneResponse{
			Result: result,
		})
		if err != nil {
			log.Fatalf("Error sending response %v", err)
			return err
		}
	}
}

func (s server) LongGreet(greetServer greetpb.GreetService_LongGreetServer) error {
	fmt.Println("LongGreet function was invoked")
	result := ""
	for {
		msg, err := greetServer.Recv()
		if err == io.EOF {
			return greetServer.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
		}
		if err != nil {
			log.Fatalf("Encountered error in LongGreet %v", err)
			return err
		}
		firstName := msg.GetGreeting().GetFirstName()
		lastName := msg.GetGreeting().GetLastName()
		result += "Hello " + firstName + " " + lastName + "! "
	}
}

func (s server) GreetManyTimes(request *greetpb.GreetManyTimesRequest, timesServer greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Printf("GreetManyTimes function was invoked with %v\n", request)

	firstName := request.GetGreeting().GetFirstName()
	lastName := request.GetGreeting().GetLastName()
	for i := 0; i < 10; i++ {
		result := "Hello " + firstName + " " + lastName + " number " + strconv.Itoa(i)
		res := greetpb.GreetManyTimesResponse{
			Result: result,
		}
		timesServer.Send(&res)
		time.Sleep(1000 * time.Millisecond)
	}

	return nil
}

func (s server) Greet(ctx context.Context, request *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet function was invoked with %v\n", request)
	firstName := request.GetGreeting().GetFirstName()
	lastName := request.GetGreeting().GetLastName()

	result := "Hello " + firstName + " " + lastName
	res := greetpb.GreetResponse{
		Result: result,
	}
	return &res, nil
}

func main() {
	fmt.Println("Hello")

	lis, err := net.Listen("tcp", "0.0.0.0:50051") //default port for grpc
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to server: %v", err)
	}
}
