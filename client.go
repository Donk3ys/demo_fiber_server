package main

import (
	"context"
	"fmt"
	"fiber_demo/greetpb"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
)

func main() {

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	// Close connection on client exit
	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)
	//fmt.Printf("Created client: %f", c)

	//unaryCall(c)

	biDiCall(c)

}



func unaryCall(c greetpb.GreetServiceClient) {
	// Call Greet Service
	req := greetpb.GreetRequest {
		Greeting: &greetpb.Greeting{
			FirstName: "David",
			LastName: "Gericke",
		},
	}
	res, err :=	c.Greet(context.Background(), &req)
	if err != nil {
		log.Fatalf("Error while calling Greet RPC: %v", err)
	}

	log.Printf("Response from Greet: %v", res.Result)
}


func biDiCall(c greetpb.GreetServiceClient) {
	fmt.Println("Starting BiDi Streaming RPC...")

	// Invoke stream connection -> client to server
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream: %v", err)
		return
	}

	// Create blocking func
	waitc := make(chan struct{}) // channel that doesnt take any data

	// Send messages to client in a separate thread (go routine)
	go func()  {
		requests := mockGreetRequests()
		for _, req := range requests {
			fmt.Printf("Sending message %v\n", req)
			stream.Send(req)
			time.Sleep(2 * time.Second)
		}

		// Close the stream once all requests sent
		stream.CloseSend()
	}() // Calls annon func

	// Recieve messages from server in a separate thread (go routine)
	go func()  {
		// Listen to incomming stream
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				// close wait channel once stream send/rec finished
				break;
			}
			if err != nil {
				log.Fatalf("Error while recieving %v", err)
				break;
			}

			// Print response result
			fmt.Printf("Recieved: %v\n\n", res.GetResult())
		}

		close(waitc)
	}() // Calls annon func

	// block until wait channel is closed
	<-waitc

	log.Println("Stream closed...")
}

func mockGreetRequests() []*greetpb.GreetEveryoneRequest {
	return []*greetpb.GreetEveryoneRequest{
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "David",
				LastName: "Gericke",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Tom",
				LastName: "Tomson",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Frank",
				LastName: "Frankson",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Bob",
				LastName: "Bobson",
			},
		},
	}

}
