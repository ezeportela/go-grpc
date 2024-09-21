package main

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/ezeportela/go-grpc/testpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cc, err := grpc.NewClient("localhost:5061", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	defer cc.Close()

	c := testpb.NewTestServiceClient(cc)

	// t := doUnary(c)
	// doClientStreaming(c, t.Id)
	// doServerStreaming(c, t.Id)
	doBidirectionalStreaming(c)
}

func doUnary(c testpb.TestServiceClient) *testpb.Test {
	req := &testpb.GetTestRequest{
		Id: "t1",
	}

	res, err := c.GetTest(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling GetTest RPC: %v", err)
	}

	log.Printf("Response from GetTest: %v", res)

	return res
}

func doClientStreaming(c testpb.TestServiceClient, testId string) {
	questions := []*testpb.Question{
		{
			Id:       "t1q1",
			Question: "What is Golang?",
			Answer:   "A programming language",
			TestId:   testId,
		},
		{
			Id:       "t1q2",
			Question: "What is a struct?",
			Answer:   "A data structure",
			TestId:   testId,
		},
		{
			Id:       "t1q3",
			Question: "What is a pointer?",
			Answer:   "A reference to a memory address",
			TestId:   testId,
		},
	}

	stream, err := c.SetQuestion(context.Background())
	if err != nil {
		log.Fatalf("error while calling SetQuestion RPC: %v", err)
	}

	for _, question := range questions {
		log.Println("Sending question:", question)
		if err := stream.Send(question); err != nil {
			log.Fatalf("error while sending question: %v", err)
		}
		time.Sleep(3 * time.Second)
	}

	msg, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while closing stream: %v", err)
	}

	log.Printf("Response from SetQuestion: %v", msg)
}

func doServerStreaming(c testpb.TestServiceClient, testId string) {
	req := &testpb.GetStudentsPerTestRequest{
		TestId: testId,
	}

	stream, err := c.GetStudentsPerTest(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling GetQuestions RPC: %v", err)
	}

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}

		log.Printf("Response from GetStudentsPerRequest: %v", msg)
	}
}

func doBidirectionalStreaming(c testpb.TestServiceClient) {
	answer := testpb.TakeTestRequest{
		Answer: "A programming language",
	}

	numberOfQuestions := 4

	waitChannel := make(chan struct{})

	stream, err := c.TakeTest(context.Background())
	if err != nil {
		log.Fatalf("error while calling TakeTest RPC: %v", err)
	}

	go func() {
		for i := 0; i < numberOfQuestions; i++ {
			log.Println("Sending answer:", answer)
			if err := stream.Send(&answer); err != nil {
				log.Fatalf("error while sending answer: %v", err)
			}
			time.Sleep(3 * time.Second)
		}
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("error while reading stream: %v", err)
				break
			}

			log.Printf("Response from TakeTest: %v", res)
		}
		close(waitChannel)
	}()

	<-waitChannel
}
