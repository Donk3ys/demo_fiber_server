package main

import (
	"context"
	"fiber_demo/greetpb"
	handler "fiber_demo/handlers"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

// for gRPC
type svr struct{
	greetpb.UnimplementedGreetServiceServer
}

func main() {
	// Create TCP Connection for gRPC
	listen, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create gRPC server
	s := grpc.NewServer()
	// Register service to server
	greetpb.RegisterGreetServiceServer(s, &svr{})


	// Serve on connection
	err = s.Serve(listen)
	if err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}



	// Set environmental variables
	host, port := getEnvs()

	// TODO Add db

	// Create App
	app := fiber.New()

	// Route Handler
	gh := handler.GeneralHandler {}
  app.Get("/", gh.SayHi)
	app.Get("/:id", gh.GetPersonMatchingId)
	app.Post("/", gh.PersonCreds)

	// Run Server
	log.Fatal(app.Listen(fmt.Sprintf("%s:%s", host, port)))
}


// Checks if envs are set in os or .env file
func getEnvs() (string, string) {
	host := os.Getenv("SERVER_HOST")
	port := os.Getenv("SERVER_PORT")

	if host == "" && port == "" {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal("Environment variables not set")
		}

		host = os.Getenv("SERVER_HOST")
		port = os.Getenv("SERVER_PORT")
	}

	return host, port
}



// Recieves a greeting<Greeting> -> Responds with result<string>	
// Func Greet implements greetpb.Greet : Found in _grpc.pb.go
func (*svr) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet func called with %v", req)

	firstName := req.GetGreeting().GetFirstName()
	lastName := req.GetGreeting().GetLastName()

	result := "Hello " + firstName + " " + lastName; 
	res := greetpb.GreetResponse {
		Result: result,
	}

	return &res, nil
}

// Recieves a greeting<Greeting> -> Responds with result<string> or nil/err to close connection	
func (*svr) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	fmt.Printf("GreetEveryone streram func called")

	for {
		// Stream recieves a request
		req, err := stream.Recv()
		if err == io.EOF {	// if connection closed
			return nil
		}
		if err != nil {
			log.Fatalf("Error while recieving client stream: %v", err)
			return err
		}

		// Get data from request
		firstName := req.GetGreeting().GetFirstName()
		lastName := req.GetGreeting().GetLastName()

		// create result & response
		result := "Hello " + firstName + " " + lastName
		res := greetpb.GreetEveryoneResponse{ Result: result }

		// send response
		err = stream.Send(&res)
		if err != nil {
			log.Fatalf("Error while sending data to client: %v", err)
			return err
		}
	} 

}
