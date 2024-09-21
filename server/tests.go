package server

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/ezeportela/go-grpc/models"
	"github.com/ezeportela/go-grpc/repositories"
	"github.com/ezeportela/go-grpc/studentpb"
	"github.com/ezeportela/go-grpc/testpb"
)

type TestServer struct {
	repo repositories.Repository
	testpb.UnimplementedTestServiceServer
}

func NewTestServer(repo repositories.Repository) *TestServer {
	return &TestServer{
		repo: repo,
	}
}

func (s *TestServer) GetTest(ctx context.Context, req *testpb.GetTestRequest) (*testpb.Test, error) {
	test, err := s.repo.GetTest(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &testpb.Test{
		Id:   test.Id,
		Name: test.Name,
	}, nil
}

func (s *TestServer) SetTest(ctx context.Context, req *testpb.Test) (*testpb.SetTestResponse, error) {
	test := &models.Test{
		Id:   req.Id,
		Name: req.Name,
	}

	if err := s.repo.SetTest(ctx, test); err != nil {
		return nil, err
	}

	return &testpb.SetTestResponse{
		Id:   test.Id,
		Name: test.Name,
	}, nil
}

func (s *TestServer) SetQuestion(stream testpb.TestService_SetQuestionServer) error {
	questions := make([]*models.Question, 0)
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error reading stream: %v", err)
			return err
		}
		questions = append(questions, &models.Question{
			Id:       msg.Id,
			TestId:   msg.TestId,
			Question: msg.Question,
			Answer:   msg.Answer,
		})
	}

	for _, question := range questions {
		if err := s.repo.SetQuestion(stream.Context(), question); err != nil {
			return stream.SendAndClose(&testpb.SetQuestionResponse{
				Ok: false,
			})
		}
	}

	return stream.SendAndClose(&testpb.SetQuestionResponse{
		Ok: true,
	})
}

func (s *TestServer) EnrollStudents(stream testpb.TestService_EnrollStudentsServer) error {
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&testpb.SetQuestionResponse{
				Ok: true,
			})
		}
		if err != nil {
			log.Fatalf("Error reading stream: %v", err)
			return stream.SendAndClose(&testpb.SetQuestionResponse{
				Ok: false,
			})
		}

		err = s.repo.SetEnrollment(stream.Context(), &models.Enrollment{
			StudentId: msg.StudentId,
			TestId:    msg.TestId,
		})
		if err != nil {
			return stream.SendAndClose(&testpb.SetQuestionResponse{
				Ok: false,
			})
		}
	}
}

func (s *TestServer) GetStudentsPerTest(req *testpb.GetStudentsPerTestRequest, stream testpb.TestService_GetStudentsPerTestServer) error {
	students, err := s.repo.GetStudentsPerTest(stream.Context(), req.TestId)
	if err != nil {
		return err
	}

	for _, student := range students {
		if err := stream.Send(&studentpb.Student{
			Id:   student.Id,
			Name: student.Name,
			Age:  student.Age,
		}); err != nil {
			return err
		}
		time.Sleep(3 * time.Second)
	}

	return nil
}
