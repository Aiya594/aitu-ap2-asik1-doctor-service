package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/Aiya594/doctor-service/internal/cache"
	natspub "github.com/Aiya594/doctor-service/internal/event"
	"github.com/Aiya594/doctor-service/internal/model"
	"github.com/Aiya594/doctor-service/internal/repository"
	"github.com/google/uuid"
)

// Doctor Service Rules
//  full_name is required.
//  email is required.
//  email must be unique across all doctors.

type DocUseCase interface {
	CreateDoc(fullName, email, specialization string) (string, error)
	GetDocbyID(id string) (*model.Doctor, error)
	ListDoctors() ([]*model.Doctor, error)
}

type DoctorUsecaseImpl struct {
	repo      repository.DoctorRepository
	publisher natspub.EventPublisher
	logger    *slog.Logger
	cache     cache.CacheRepository
}

func NewDoctorUseCase(repo repository.DoctorRepository, logger *slog.Logger, pub natspub.EventPublisher, cacheRepo cache.CacheRepository) DocUseCase {
	return &DoctorUsecaseImpl{
		repo: repo, logger: logger, publisher: pub, cache: cacheRepo,
	}
}

func (d *DoctorUsecaseImpl) CreateDoc(fullName, email, specialization string) (string, error) {

	fullName = strings.TrimSpace(fullName)
	email = strings.ToLower(strings.TrimSpace(email))
	specialization = strings.ToLower(strings.TrimSpace(specialization))

	if fullName == "" || email == "" || specialization == "" {
		d.logger.Error("failed create a doctor",
			"error", ErrInvalidFields, "full_name", fullName, "email", email, "specialization", specialization)
		return "", fmt.Errorf("full name, email and specialization are required:%w", ErrInvalidFields)
	}

	exists := d.repo.ExistsByEmail(email)

	if exists {
		d.logger.Error("failed create a doctor", "error", ErrAlreadyExists, "email", email)
		return "", ErrAlreadyExists
	}

	id := uuid.New().String()

	doctor := &model.Doctor{
		ID:             id,
		FullName:       fullName,
		Email:          email,
		Specialization: specialization,
	}

	createdAt := time.Now()

	err := d.repo.Create(doctor, createdAt)
	if err != nil {
		d.logger.Error("couldnt create a doctor", "error", err, "ID", doctor.ID)
		return id, err
	}
	d.logger.Info("doctor created successfully")

	// Write-Through: invalidate list cache after successful DB write
	ctx := context.Background()
	if err := d.cache.InvalidateDoctorList(ctx); err != nil {
		d.logger.Error("cache invalidation failed for doctors:list", "error", err)
		// best-effort: do not block
	}

	eventType := model.DoctorCreatedEventName

	event := model.DoctorCreated{
		EventType:      eventType,
		OccurredAt:     createdAt,
		ID:             id,
		Full_name:      fullName,
		Specialization: specialization,
		Email:          email,
	}

	data, err := json.Marshal(event)
	if err != nil {
		d.logger.Error("couldnt marshal event", "error", err)
		return id, err
	}

	err = d.publisher.Publish(eventType, data)
	if err != nil {
		d.logger.Error("couldnt publish event", "error", err)
		return id, err
	}
	d.logger.Info("event published successfully")
	return id, nil
}

func (d *DoctorUsecaseImpl) GetDocbyID(id string) (*model.Doctor, error) {
	id = strings.TrimSpace(id)

	if id == "" {
		d.logger.Error("failed to get doctor", "error", ErrInvalidFields, "id", id)
		return nil, ErrInvalidFields
	}

	ctx := context.Background()

	// Cache-Aside: try cache first
	cached, err := d.cache.GetDoctor(ctx, id)
	if err != nil {
		d.logger.Error("cache read error, falling through to DB", "error", err, "id", id)
	}
	if cached != nil {
		d.logger.Info("cache hit", "key", "doctor:"+id)
		return cached, nil
	}

	// Cache miss — fall through to DB transparently
	doctor, err := d.repo.GetByID(id)
	if err != nil {
		d.logger.Error("failed to get a doctor", "error", err, "id", id)
		return nil, ErrNotFound
	}

	// Populate cache (best-effort)
	if cerr := d.cache.SetDoctor(ctx, doctor); cerr != nil {
		d.logger.Error("cache write failed", "error", cerr, "id", id)
	}
	return doctor, nil

}

func (d *DoctorUsecaseImpl) ListDoctors() ([]*model.Doctor, error) {
	ctx := context.Background()

	// Cache-Aside
	cached, err := d.cache.GetDoctorList(ctx)
	if err != nil {
		d.logger.Error("cache read error for doctors:list, falling through to DB", "error", err)
	}
	if cached != nil {
		d.logger.Info("cache hit", "key", "doctors:list")
		return cached, nil
	}

	doctors, err := d.repo.List()
	if err != nil {
		d.logger.Error("failed to list doctors", "error", err)
		return nil, err
	}

	// Populate cache (best-effort)
	if cerr := d.cache.SetDoctorList(ctx, doctors); cerr != nil {
		d.logger.Error("cache write failed for doctors:list", "error", cerr)
	}
	return doctors, nil
}
