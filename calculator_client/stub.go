package main

import (
	"calc/calculatorpb"
	"context"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func main() {

	calculatorServiceConn, err := grpc.Dial("localhost:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("Dial error:", err.Error())
		return
	}

	stub := calculatorpb.NewCalculatorServiceClient(calculatorServiceConn)

	// DoUnary(stub)
	// DoServerStream(stub)
	// DoClientStream(stub)
	// DoBiDirectionalStream(stub)
	// DoPerfectNumber(stub)
	DoTotalNumber(stub)
	// DoFindMinimum(stub)
	// DoAdd(stub)
	// DoSquareRoot(stub)
}

func DoUnary(stub calculatorpb.CalculatorServiceClient) {

	log.Println("Request Sum ------------->")

	resp, err := stub.Sum(context.Background(), &calculatorpb.SumRequest{
		FirstNumber:  10,
		SecondNumber: 20,
	})
	if err != nil {
		log.Println("Error while sum function:", err.Error())
		return
	}

	log.Printf("Response Sum <------------- %+v\n", resp)
}

func DoServerStream(stub calculatorpb.CalculatorServiceClient) {

	log.Println("Request PrimeNumberDecompostion ------------->")

	stream, err := stub.PrimeNumberDecompostion(context.Background(), &calculatorpb.PrimeNumberDecompostionRequest{
		PrimeFactor: 40,
	})
	if err != nil {
		log.Println("Error while PrimeNumberDecompostion stream function:", err.Error())
		return
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Println("Error while PrimeNumberDecompostion function:", err.Error())
			return
		}

		log.Printf("Response PrimeNumberDecompostion <------------- %+v\n", resp)
	}
}

func DoClientStream(stub calculatorpb.CalculatorServiceClient) {

	log.Println("Request ComputeAverage ------------->")

	stream, err := stub.ComputeAverage(context.Background())
	if err != nil {
		log.Println("Error while ComputeAverage stream function:", err.Error())
		return
	}

	numbers := []int32{2, 6, 1, 9, 2, 8, 5}

	for _, number := range numbers {
		stream.Send(&calculatorpb.ComputeAverageRequest{
			Number: number,
		})
		time.Sleep(time.Second * 1)
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Println("Error while ComputeAverage stream close and recv function:", err.Error())
		return
	}

	log.Printf("Response ComputeAverage <------------- %+v\n", resp)
}

func DoBiDirectionalStream(stub calculatorpb.CalculatorServiceClient) {

	log.Println("Request FindMaximum ------------->")

	stream, err := stub.FindMaximum(context.Background())
	if err != nil {
		log.Println("Error while FindMaximum stream function:", err.Error())
		return
	}

	waitch := make(chan struct{})

	go func() {
		numbers := []int32{2, 6, 1, 9, 2, 8, 5}

		for _, number := range numbers {
			stream.Send(&calculatorpb.FindMaximumRequest{
				Number: number,
			})
		}

		stream.CloseSend()
	}()

	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				break
			}

			if err != nil {
				log.Println("Error while FindMaximum stream  recv function:", err.Error())
				return
			}

			log.Printf("Response FindMaximum <------------- %+v\n", resp)
		}
		waitch <- struct{}{}
	}()

	<-waitch
}

func DoPerfectNumber(stub calculatorpb.CalculatorServiceClient) {

	log.Println("Request PerfectNumber ------------->")

	stream, err := stub.PerfectNumber(context.Background(), &calculatorpb.PerfectNumberRequest{
		Min: 1,
		Max: 1000,
	})
	if err != nil {
		log.Println("Error while PerfectNumber stream function:", err.Error())
		return
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Println("Error while PerfectNumber function:", err.Error())
			return
		}

		log.Printf("Response PerfectNumber <------------- %+v\n", resp)
	}
}

func DoTotalNumber(stub calculatorpb.CalculatorServiceClient) {

	log.Println("Request TotalNumber ------------->")

	stream, err := stub.TotalNumber(context.Background())
	if err != nil {
		log.Println("Error while TotalNumber stream function:", err.Error())
		return
	}

	numbers := []int32{2, 6, 1, 9, 2, 8, 5}

	for _, number := range numbers {
		stream.Send(&calculatorpb.TotalNumberRequest{
			Number: number,
		})
		time.Sleep(time.Second * 1)
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Println("Error while TotalNumber stream close and recv function:", err.Error())
		return
	}

	log.Printf("Response TotalNumber <------------- %+v\n", resp)
}

func DoFindMinimum(stub calculatorpb.CalculatorServiceClient) {

	log.Println("Request FindMinimum ------------->")

	stream, err := stub.FindMinimum(context.Background())
	if err != nil {
		log.Println("Error while FindMinimum stream function:", err.Error())
		return
	}

	numbers := []int32{2, 6, 1, 9, 2, 8, 5}

	for _, number := range numbers {
		stream.Send(&calculatorpb.FindMinimumRequest{
			Number: number,
		})
		time.Sleep(time.Second * 1)
	}

	stream.CloseSend()

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Println("Error while FindMinimum stream recv function:", err.Error())
			return
		}

		log.Printf("Response FindMinimum <------------- %+v\n", resp)
	}
}

func DoSquareRoot(stub calculatorpb.CalculatorServiceClient) {

	log.Println("Request SquareRoot ------------->")

	resp, err := stub.SquareRoot(context.Background(), &calculatorpb.SquareRootRequest{
		Number: 25,
	})
	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok && statusErr.Code() == codes.InvalidArgument {
			log.Printf("Invalid argument error: %v", statusErr.Message())
		} else {
			log.Println("Error while SquareRoot function:", err.Error())
		}
		return
	}

	log.Printf("Response SquareRoot <------------- %+v\n", resp)
}

func DoAdd(stub calculatorpb.CalculatorServiceClient) {

	log.Println("Request Add ------------->")

	resp, err := stub.Add(context.Background(), &calculatorpb.AddRequest{
		FirstNumber:  30,
		SecondNumber: 40,
	})
	if err != nil {
		log.Println("Error while Add function:", err.Error())
		return
	}

	log.Printf("Response Add <------------- %+v\n", resp)
}
