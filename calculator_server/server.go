package main

import (
	"calc/calculatorpb"
	"context"
	"io"
	"log"
	"math"
	"net"
	"time"

	"errors"

	"google.golang.org/grpc"
)

func main() {

	lis, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Println("Listen error:", err.Error())
		return
	}

	s := grpc.NewServer()

	calculatorpb.RegisterCalculatorServiceServer(s, &Service{})

	log.Println("Listenining RPC Server 9090....")
	err = s.Serve(lis)
	if err != nil {
		log.Println("Serve error:", err.Error())
		return
	}
}

type Service struct {
	calculatorpb.UnimplementedCalculatorServiceServer
}

func (s *Service) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {

	log.Printf("Sum function request ~~~~~~~~~~~~>>> %+v\n", req)

	var result = req.GetFirstNumber() + req.GetSecondNumber()

	resp := calculatorpb.SumResponse{SumResult: result}

	return &resp, nil
}

func (s *Service) PrimeNumberDecompostion(req *calculatorpb.PrimeNumberDecompostionRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompostionServer) error {

	log.Printf("PrimeNumberDecompostion function request ~~~~~~~~~~~~>>> %+v\n", req)

	number := req.PrimeFactor
	divisor := int64(2)

	for number > 1 {
		if number%divisor == 0 {
			stream.Send(&calculatorpb.PrimeNumberDecompostionResponse{
				Number: divisor,
			})
			time.Sleep(time.Second * 1)
			number /= divisor
		} else {
			divisor++
		}
	}

	return nil
}

func (s *Service) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	sum := int32(0)
	count := 0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			average := float64(sum) / float64(count)
			err := stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
				Average: average,
			})
			if err != nil {
				log.Println("ComputeAverage stream send and close error:", err.Error())
				return err
			}
			return nil
		}

		if err != nil {
			log.Println("Compute Average stream recieve error:", err.Error())
			return err
		}

		sum += req.Number
		count++

		log.Printf("ComputeAverage function request ~~~~~~~~~~~~>>> %+v\n", req)
	}
}

func (s *Service) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {
	var max int32
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Println("FindMaximum stream recieve error:", err.Error())
			return err
		}

		if req.Number > max {
			max = req.Number
			err := stream.Send(&calculatorpb.FindMaximumResponse{
				Maximum: max,
			})
			if err != nil {
				log.Println("FindMaximum stream send error:", err.Error())
				return err
			}
		}

		log.Printf("FindMaximum function request ~~~~~~~~~~~~>>> %+v\n", req)
		time.Sleep(time.Second * 1)
	}
}

func (s *Service) PerfectNumber(req *calculatorpb.PerfectNumberRequest, stream calculatorpb.CalculatorService_PerfectNumberServer) error {
	log.Printf("PerfectNumber function request ~~~~~~~~~~~~>>> %+v\n", req)

	perfMin := req.GetMin()
	perfMax := req.GetMax()

	for perfNum := perfMin; int32(perfNum) <= int32(perfMax); perfNum++ {
		perfsum := int32(0)
		for i := int32(1); i < perfNum; i++ {
			if perfNum%i == 0 {
				perfsum += i
			}
		}
		if perfNum == perfsum {
			if err := stream.Send(&calculatorpb.PerfectNumberResponse{PerfectNumber: perfNum}); err != nil {
				log.Printf("Error sending response: %v", err)
				return err
			}
		}
	}

	return nil
}

func (s *Service) TotalNumber(stream calculatorpb.CalculatorService_TotalNumberServer) error {
	log.Println("TotalNumber function started...")

	var total int32
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			response := &calculatorpb.TotalNumberResponse{Total: total}
			if err := stream.SendAndClose(response); err != nil {
				log.Println("TotalNumber stream send and close error:", err.Error())
				return err
			}
			log.Println("TotalNumber function completed.")
			return nil
		}

		if err != nil {
			log.Println("TotalNumber stream receive error:", err.Error())
			return err
		}

		total += req.Number
		log.Printf("TotalNumber function received ~~~~~~~~~~~~>>> %+v\n", req)
	}
}

func (s *Service) FindMinimum(stream calculatorpb.CalculatorService_FindMinimumServer) error {
	log.Println("FindMinimum function started...")

	var min int32
	firstItem := true

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			log.Println("FindMinimum function completed.")
			return nil
		}

		if err != nil {
			log.Println("FindMinimum stream receive error:", err.Error())
			return err
		}

		if firstItem || req.Number < min {
			min = req.Number
			firstItem = false
		}

		response := &calculatorpb.FindMinimumResponse{Minimum: min}
		if err := stream.Send(response); err != nil {
			log.Println("FindMinimum stream send error:", err.Error())
			return err
		}

		log.Printf("FindMinimum function received ~~~~~~~~~~~~>>> %+v\n", req)
	}
}

func (s *Service) SquareRoot(ctx context.Context, req *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootResponse, error) {
	log.Printf("SquareRoot function request ~~~~~~~~~~~~>>> %+v\n", req)

	number := req.GetNumber()

	if number < 0 {
		return nil, errors.New(
			"cannot calculate square root of a negative number: %v",
		)
	}

	result := float64(number)
	result = math.Sqrt(result)

	resp := &calculatorpb.SquareRootResponse{SquareRoot: result}
	return resp, nil
}

func (s *Service) Add(ctx context.Context, req *calculatorpb.AddRequest) (*calculatorpb.AddResponse, error) {
	log.Printf("Add function request ~~~~~~~~~~~~>>> %+v\n", req)

	result := req.GetFirstNumber() + req.GetSecondNumber()

	resp := &calculatorpb.AddResponse{SumResult: result}
	return resp, nil
}
