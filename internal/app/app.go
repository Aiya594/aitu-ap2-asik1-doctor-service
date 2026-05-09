package app

import (
	"database/sql"
	"log"
	"log/slog"
	"net"
	"os"

	"github.com/Aiya594/doctor-service/internal/cache"
	"github.com/Aiya594/doctor-service/internal/config"
	natspub "github.com/Aiya594/doctor-service/internal/event"
	"github.com/Aiya594/doctor-service/internal/middleware"
	"github.com/Aiya594/doctor-service/internal/repository"
	grpcDoc "github.com/Aiya594/doctor-service/internal/transport/grpc"
	usecase "github.com/Aiya594/doctor-service/internal/use-case"
	"github.com/Aiya594/doctor-service/proto"
	"github.com/golang-migrate/migrate/v4"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type App struct {
	grpcServ    *grpc.Server
	logger      *slog.Logger
	pub         *natspub.Publisher
	db          *sql.DB
	redisClient *redis.Client
}

func NewApp() (*App, error) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	cfg := config.NewConfig()

	runMigrations(cfg.ConnStrDB)

	db, err := cfg.ConnectDB()
	if err != nil {
		return nil, err
	}

	publisher, err := natspub.NewPublisher(cfg.NatsURL)
	if err != nil {
		return nil, err
	}

	// Redis — optional; if unavailable, NoopCache is used
	redisClient := cache.NewRedisClient(logger)
	var cacheRepo cache.CacheRepository
	if redisClient != nil {
		cacheRepo = cache.NewRedisCacheRepository(redisClient, logger)
	} else {
		cacheRepo = cache.NewNoop()
	}

	repo := repository.NewDoctorRepository(db)
	uc := usecase.NewDoctorUseCase(repo, logger, publisher, cacheRepo)
	handler := grpcDoc.NewDoctorServer(uc, logger)

	rateLimiter := middleware.RateLimiterInterceptor(redisClient, logger)
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(rateLimiter),
	)
	proto.RegisterDoctorServiceServer(grpcServer, handler)

	return &App{
		grpcServ:    grpcServer,
		logger:      logger,
		pub:         publisher,
		db:          db,
		redisClient: redisClient,
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
	if a.redisClient != nil {
		a.redisClient.Close()
	}
}

func (a *App) Stop() {
	a.grpcServ.GracefulStop()
}

func runMigrations(dbURL string) {
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}
	m, err := migrate.New("file://migrations", dbURL)
	if err != nil {
		log.Fatal("migration init error:", err)
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatal("migration failed:", err)
	}
	log.Println("migrations applied successfully")
}
