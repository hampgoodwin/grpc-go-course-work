package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/abelgoodwin1988/grpc-go-course-work/calculator/calculatorpb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) Sum(ctx context.Context, in *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	fmt.Printf("Sum Function Invoked with: %v\n", in)
	ints := in.GetSum()
	var sum int32
	for _, val := range ints {
		sum += val
	}
	res := &calculatorpb.SumResponse{
		Sum: sum,
	}
	return res, nil
}

func (*server) PrimeNumberDecomposition(req *calculatorpb.PrimeNumberDecomponsitionRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	fmt.Printf("PrimeNumberDecomposition function invoked with: %v\n", req)
	number := req.GetNumber()
	k := int32(2)
	for {
		if number < 2 {
			break
		}
		if number%k == 0 {
			stream.Send(&calculatorpb.PrimeNumberDecomponsitionResponse{Decomposition: k})
			number = number / k
		} else {
			k++
		}
	}
	return nil
}

func (*server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	fmt.Println("ComputeAverage was invoked with a stream")
	i := int32(1)
	sum := int32(0)
	average := int32(0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			log.Printf("ComputeAverage reached end of stream\n")
			return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
				Average: average,
			})
		}
		if err != nil {
			log.Fatalf("Failed to receive value from stream: %\n", err)
		}
		sum += req.GetAverageSubject()
		average = sum / i
		log.Printf("New Average: %v\n from adding %v", average, req.GetAverageSubject())
		i++
	}
}

func main() {
	fmt.Printf("CalculatorServiceServer Started\n")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatal("Listener Failed: %v", err)
	}
	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to Server CalculatorService Server: %v", err)
	}
}