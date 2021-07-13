package main

import (
	"context"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"testing"
	"time"

	pd "github.com/tonymj76/meant4-factorial/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

var (
	listen     *bufconn.Listener
	bufferSize = 1024 * 1024
)

func bufDialer(context.Context, string) (net.Conn, error) {
	return listen.Dial()
}

// type mockSiteService_Factorial_CalculateServer struct {
// 	grpc.ServerStream
// 	Results []*pd.CalculateRequest
// }

// func (fc *mockSiteService_Factorial_CalculateServer) Send(calReq *pd.CalculateRequest) error {
// 	fc.Results = append(fc.Results, calReq)
// 	return nil
// }

func TestMain(m *testing.M) {
	code := 0
	defer func() {
		os.Exit(code)
	}()
	listen = bufconn.Listen(bufferSize)
	s := grpc.NewServer()
	pd.RegisterFactorialServer(s, &service{})
	go func() {
		if err := s.Serve(listen); err != nil {
			log.Fatalf("server exited with error %v", err)
		}
	}()

	code = m.Run()
}

func TestFactorial(t *testing.T) {
	posibleAns := []int{120, 5040, 0}
	reqData := &pd.CalculateRequest{Numbers: []int64{5, 7, 0}}

	count := 0
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := pd.NewFactorialClient(conn)
	resp, err := client.Calculate(ctx, reqData)
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	for {
		in, err := resp.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("failed to receive a number %v", err)
		}
		num, _ := strconv.Atoi(in.GetFactorialResult())
		if num != posibleAns[count] {
			t.Errorf("factorial of (%d!) was incorrect, got %d but want %d", in.GetInputNumber(), posibleAns[count], num)
		}
		count++
	}
}
