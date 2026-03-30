package app

import (
	"log/slog"
	"os"

	"github.com/Aiya594/doctor-service/internal/repository"
	httpdoc "github.com/Aiya594/doctor-service/internal/transport/http"
	usecase "github.com/Aiya594/doctor-service/internal/use-case"
	"github.com/gin-gonic/gin"
)

type App struct {
	router *gin.Engine
}

func NewApp() *App {
	repo := repository.NewDocRepo()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	usecase := usecase.NewDoctorUseCase(repo, logger)

	handler := httpdoc.NewDocHandler(usecase)

	r := gin.Default()

	httpdoc.RegisterRoutes(r, handler)

	return &App{router: r}
}

func (a *App) RunServer(port string) {

	a.router.Run(":" + port)
}
