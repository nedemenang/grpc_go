package main

import (
	"math"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"
	"google.golang.org/grpc/reflection"
	"github.com/nedemenang/grpc_go/calculator/calculatorpb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	fmt.Printf("Sum function was invoked with %v\n", req)
	firstDigit := req.GetDigits().GetFirstNum()
	lastDigit := req.GetDigits().GetSecondNum()
	result := firstDigit + lastDigit

	res := &calculatorpb.SumResponse{
		Result: result,
	}
	return res, nil
}

func (*server) PrimeNumberDecomposition(req *calculatorpb.PrimeNumberRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	fmt.Printf("Prime number decomposition with %v\n", req)
	primeNumber := req.GetPrimeNumber()
	var k int32 = 2
	for primeNumber > 1 {
		if primeNumber%k == 0 {
			res := &calculatorpb.PrimeNumberResponse{
				Result: k,
			}
			primeNumber = primeNumber / k
			stream.Send(res)
			time.Sleep(1000 * time.Millisecond)
		} else {
			k = k + 1
		}
	}
	return nil
}

func (*server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	fmt.Printf("Compute average was invoked with a sreaming request\n")
	numbers := 0.0
	count := 0.0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			average := numbers / count
			return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
				Number: average,
			})
		}
		if err != nil {
			log.Fatalf("Error while reading client stream %v", err)
		}
		numbers += req.GetNumber()
		count++
	}

}

func (*server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {
	fmt.Printf("FindMaximum function was invoked with a streaming request\n")
	comparative := int32(0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading client stream %v", err)
		}

		result := req.GetNumber()
		if comparative > result {
			result = comparative
		}
		sendError := stream.Send(&calculatorpb.FindMaximumResponse{
			Number: result,
		})
		if sendError != nil {
			log.Fatalf("Error while sending data to client: %v", err)
			return err
		}
		comparative = result

	}
}

func(*server) SquareRoot(ctx context.Context, req *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootReponse, error) {
	fmt.Println("Recieve square root RPC")
	number := req.GetNumber()

	if (number < 0) {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Recieved a negative number: %v\n", number),
		)
	}
	return &calculatorpb.SquareRootReponse{
		NumberRoot: math.Sqrt(float64(number)),
	}, nil
}

func main() {
	fmt.Println("Calculator Server")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listed: %v", err)
	}

	s := grpc.NewServer()

	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
