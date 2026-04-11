package app

import (
	"log/slog"
	"net"
	"os"

	"github.com/Aiya594/doctor-service/internal/repository"
	grpcDoc "github.com/Aiya594/doctor-service/internal/transport/grpc"
	usecase "github.com/Aiya594/doctor-service/internal/use-case"
	"github.com/Aiya594/doctor-service/proto"
	"google.golang.org/grpc"
)

type App struct {
	grpcSrev *grpc.Server
	logger   *slog.Logger
}

func NewApp() *App {
	repo := repository.NewDocRepo()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	usecase := usecase.NewDoctorUseCase(repo, logger)

	handler := grpcDoc.NewDoctorServer(usecase, logger)

	grpcServer := grpc.NewServer()

	proto.RegisterDoctorServiceServer(grpcServer, handler)

	return &App{grpcSrev: grpcServer, logger: logger}
}

func (a *App) RunServer(port string) error {

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	err = a.grpcSrev.Serve(lis)
	if err != nil {
		return err
	}

	a.logger.Info("gRPC server started", "port", port)

	return nil

}
