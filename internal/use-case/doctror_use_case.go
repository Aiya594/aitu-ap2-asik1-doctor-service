package usecase

import (
	"log/slog"

	"github.com/Aiya594/doctor-service/internal/repository"
)

type DoctorUsecase struct {
	repo   repository.DoctorRepository
	logger *slog.Logger
}
