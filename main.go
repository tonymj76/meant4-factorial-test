package main

import (
	"errors"
	"log"
	"net"
	"os"
	"strconv"

	pd "github.com/tonymj76/meant4-factorial/proto"
	"google.golang.org/grpc"
)

type service struct{}

//Calcuate the factorial of n and return the result as a stream
func (s *service) Calculate(req *pd.CalculateRequest, stream pd.Factorial_CalculateServer) error {
	for _, num := range req.Numbers {
		if num < 0 {
			if err := stream.Send(&pd.CalculateResult{
				InputNumber:     num,
				FactorialResult: "value into Numbers shouldn't be negative",
			}); err != nil {
				return err
			}
			return errors.New("value into Numbers shouldn't be negative")
		}
		facNum := factorial(num)
		if err := stream.Send(&pd.CalculateResult{
			InputNumber:     num,
			FactorialResult: strconv.FormatInt(facNum, 10),
		}); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	//geting values from .env
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file ", err)
	// }
	port := os.Getenv("GRPC_PORT")
	// Set-up our gRPC server.
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pd.RegisterFactorialServer(s, &service{})

	log.Println("Running on port ", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// factorial calculate the factor of n
func factorial(n int64) int64 {
	if n < 2 {
		return n
	}
	return n * factorial(n-1)
}

// func factor(n, accu int) int {
// 	if n < 2 {
// 		return accu
// 	}
// 	return factor(n-1, n*accu)
// }
