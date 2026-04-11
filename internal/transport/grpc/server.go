package grpcDoc

import (
	"context"
	"log/slog"

	usecase "github.com/Aiya594/doctor-service/internal/use-case"
	"github.com/Aiya594/doctor-service/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DoctorGRPCServer struct {
	proto.UnimplementedDoctorServiceServer
	logger *slog.Logger
	svc    usecase.DocUseCase
}

func NewDoctorServer(svc usecase.DocUseCase, logger *slog.Logger) *DoctorGRPCServer {
	return &DoctorGRPCServer{logger: logger, svc: svc}
}

func (s *DoctorGRPCServer) CreateDoctor(ctx context.Context, req *proto.CreateDoctorRequest) (*proto.DoctorResponse, error) {
	id, err := s.svc.CreateDoc(
		req.GetFullName(),
		req.GetEmail(),
		req.GetSpecialization(),
	)
	if err != nil {
		s.logger.Error("CreateDoctor failed", "error", err)
		return nil, status.Error(codes.Internal, "failed to create doctor")
	}

	return &proto.DoctorResponse{
		Id:             id,
		FullName:       req.GetFullName(),
		Specialization: req.GetSpecialization(),
		Email:          req.GetEmail(),
	}, nil
}

func (s *DoctorGRPCServer) GetDoctor(ctx context.Context, req *proto.GetDoctorRequest) (*proto.DoctorResponse, error) {
	doc, err := s.svc.GetDocbyID(req.GetId())
	if err != nil {
		s.logger.Error("GetDoctor failed", "error", err)
		return nil, status.Error(codes.NotFound, "doctor not found")
	}

	return &proto.DoctorResponse{
		Id:             doc.ID,
		FullName:       doc.FullName,
		Specialization: doc.Specialization,
		Email:          doc.Email,
	}, nil
}

func (s *DoctorGRPCServer) ListDoctors(ctx context.Context, req *proto.ListDoctorsRequest) (*proto.ListDoctorsResponse, error) {
	docs, err := s.svc.ListDoctors()
	if err != nil {
		s.logger.Error("ListDoctors failed", "error", err)
		return nil, status.Error(codes.Internal, "failed to list doctors")
	}

	result := make([]*proto.DoctorResponse, 0, len(docs))

	for _, d := range docs {
		result = append(result, &proto.DoctorResponse{
			Id:             d.ID,
			FullName:       d.FullName,
			Specialization: d.Specialization,
			Email:          d.Email,
		})
	}

	return &proto.ListDoctorsResponse{
		Doctors: result,
	}, nil
}
