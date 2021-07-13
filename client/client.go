package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	pd "github.com/tonymj76/meant4-factorial/proto"
	"google.golang.org/grpc"
)

func main() {
	// geting values from .env
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file ", err)
	// }
	// // address := os.Getenv("ADDRESS")
	conn, err := grpc.Dial("localhost:5100", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pd.NewFactorialClient(conn)

	//ctx that will timeout in 10 sec
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.Calculate(ctx, &pd.CalculateRequest{Numbers: []int64{5, 7, 65, 0, -3}})
	if err != nil {
		log.Fatalf("did not stream: %v", err)
	}

	for {
		in, err := resp.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("failed to receive a number %v", err)
		}
		fmt.Printf("the factorial of %d is %s \n", in.GetInputNumber(), in.GetFactorialResult())
	}
}
