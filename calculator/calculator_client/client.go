package main

import (
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/nedemenang/grpc_go/calculator/calculatorpb"
	"google.golang.org/grpc"
	// "github.com/golang/protobuf/protoc-gen-go/grpc"
)

func main() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)

	// doUnary(c)
	// doServerStreaming(c)
	// doClientStreaming(c)
	//doBiDiStreaming(c)
	doErrorUnary(c)
}

func doUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a unary RPC...")
	req := &calculatorpb.SumRequest{
		Digits: &calculatorpb.Digits{
			FirstNum:  100,
			SecondNum: 3,
		},
	}

	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling sum RPC: %v", err)
	}
	log.Printf("Response from sum: %v", res.Result)

}

func doServerStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a server streaming RPC...")
	var n int32 = 1000
	req := &calculatorpb.PrimeNumberRequest{
		PrimeNumber: n,
	}

	resStream, err := c.PrimeNumberDecomposition(context.Background(), req)

	if err != nil {
		log.Fatalf("Error while calling PrimeNumberDecomposition RPC: %v", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while reading stream: %v", err)
		}
		log.Printf("Response from PrimeNumberDecomposition: %v", msg.GetResult())
	}
}

func doClientStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a client streaming RPC....")

	requests := []*calculatorpb.ComputeAverageRequest{
		&calculatorpb.ComputeAverageRequest{
			Number: float64(1),
		},
		&calculatorpb.ComputeAverageRequest{
			Number: float64(2),
		},
		&calculatorpb.ComputeAverageRequest{
			Number: float64(3),
		},
		&calculatorpb.ComputeAverageRequest{
			Number: float64(4),
		},
	}

	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("error while calling compute average: %v", err)
	}

	for _, req := range requests {
		fmt.Printf("Sending req: %v\n", req)
		stream.Send(req)
		time.Sleep(100 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal("Error while recieving from ComputeAverage: %v", err)
	}
	fmt.Printf("Compute Average response %v\n", res)
}

func doBiDiStreaming(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a BiDi streaming RPC....")

	waitc := make(chan struct{})
	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream: %v", err)
		return
	}
	numbers := []int32{1, 5, 3, 6, 2, 20}

	go func() {
		for _, req := range numbers {
			fmt.Printf("Sending message %v\n", req)
			stream.Send(&calculatorpb.FindMaximumRequest{
				Number: req,
			})
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {

			}
			if err != nil {
				log.Fatalf("Error while recieving stream :v%", err)
			}
			fmt.Printf("Recieved %v\n", res.GetNumber())
		}
		close(waitc)
	}()
	<-waitc
}

func doErrorUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a squareroot unary RPC...")
	
	//correct call
	doErrorCall(c, 10)
	

	doErrorCall(c, -2)
}

func doErrorCall(c calculatorpb.CalculatorServiceClient, n int32) {
	res, err := c.SquareRoot(context.Background(), &calculatorpb.SquareRootRequest{Number: n})
	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			fmt.Println(respErr.Message())
			fmt.Println(respErr.Code())
			if respErr.Code() == codes.InvalidArgument{
				fmt.Println("We probably sent a negative number!")
			}
		} else {
			log.Fatalf("Big Error calling SquareRoot: %v", err)
		}
	}
	fmt.Printf("Result of square root of %v: %v\n", n, res.GetNumberRoot())
}
