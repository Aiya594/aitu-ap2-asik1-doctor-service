package app

import (
	"database/sql"
	"log"
	"log/slog"
	"net"
	"os"

	"github.com/Aiya594/doctor-service/internal/config"
	natspub "github.com/Aiya594/doctor-service/internal/event"
	"github.com/Aiya594/doctor-service/internal/repository"
	grpcDoc "github.com/Aiya594/doctor-service/internal/transport/grpc"
	usecase "github.com/Aiya594/doctor-service/internal/use-case"
	"github.com/Aiya594/doctor-service/proto"
	"github.com/golang-migrate/migrate/v4"
	"google.golang.org/grpc"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type App struct {
	grpcServ *grpc.Server
	logger   *slog.Logger
	pub      *natspub.Publisher
	db       *sql.DB
}

func NewApp() (*App, error) {
	runMigrations()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	cfg := config.NewConfig()

	db, err := cfg.ConnectDB()
	if err != nil {
		return nil, err
	}

	publisher, err := natspub.NewPublisher(cfg.NatsURL)
	if err != nil {
		return nil, err
	}

	repo := repository.NewDoctorRepository(db)
	uc := usecase.NewDoctorUseCase(repo, logger, publisher)
	handler := grpcDoc.NewDoctorServer(uc, logger)

	grpcServer := grpc.NewServer()
	proto.RegisterDoctorServiceServer(grpcServer, handler)

	return &App{
		grpcServ: grpcServer,
		logger:   logger,
		pub:      publisher,
		db:       db,
	}, nil
}

func (a *App) RunServer(port string) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	a.logger.Info("gRPC server starting", "port", port)
	return a.grpcServ.Serve(lis)
}

func (a *App) Close() {
	a.pub.Close()
	a.db.Close()
}

func (a *App) Stop() {
	a.grpcServ.GracefulStop()
}

func runMigrations() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	m, err := migrate.New(
		"file://migrations",
		dbURL,
	)
	if err != nil {
		log.Fatal("migration init error:", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatal("migration failed:", err)
	}

	log.Println("migrations applied successfully")
}
