package usecase

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

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
}

func NewDoctorUseCase(repo repository.DoctorRepository, logger *slog.Logger, pub natspub.EventPublisher) DocUseCase {
	return &DoctorUsecaseImpl{
		repo: repo, logger: logger, publisher: pub,
	}
}

func (d *DoctorUsecaseImpl) CreateDoc(fullName, email, specialization string) (string, error) {
	//validation
	fullName = strings.TrimSpace(fullName)
	email = strings.ToLower(strings.TrimSpace(email))
	specialization = strings.ToLower(strings.TrimSpace(specialization))

	if fullName == "" || email == "" || specialization == "" {
		d.logger.Error("failed create a doctor",
			"error", ErrInvalidFields,
			"full_name", fullName,
			"email", email,
			"specialization", specialization)

		return "", fmt.Errorf("full name, email and specialization are required:%w", ErrInvalidFields)
	}

	exists := d.repo.ExistsByEmail(email)

	if exists {
		d.logger.Error("failed create a doctor",
			"error", ErrAlreadyExists,
			"email", email)

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
		d.logger.Error("couldnt create a doctor",
			"error", err,
			"ID", doctor.ID,
			"full_name", doctor.FullName,
			"email", doctor.Email,
			"specialization", doctor.Specialization)
		return id, err
	}
	d.logger.Info("doctor created succesfully")

	event := map[string]interface{}{
		"event_type":     model.DoctorCreated,
		"occurred_at":    createdAt.UTC().Format(time.RFC3339),
		"id":             doctor.ID,
		"full_name":      doctor.FullName,
		"specialization": doctor.Specialization,
		"email":          doctor.Email,
	}

	//event publishing
	data, err := json.Marshal(event)
	if err != nil {
		d.logger.Error("couldnt marshal event",
			"error", err,
			"event_type", model.DoctorCreated,
			"ID", doctor.ID,
			"full_name", doctor.FullName,
			"email", doctor.Email,
			"specialization", doctor.Specialization)
		return id, err
	}
	err = d.publisher.Publish(model.DoctorCreated, data)
	if err != nil {
		d.logger.Error("couldnt publish event",
			"error", err,
			"event_type", model.DoctorCreated,
			"ID", doctor.ID,
			"full_name", doctor.FullName,
			"email", doctor.Email,
			"specialization", doctor.Specialization)
		return id, err
	}

	d.logger.Info("event published succesfully")

	return id, nil

}

func (d *DoctorUsecaseImpl) GetDocbyID(id string) (*model.Doctor, error) {
	id = strings.TrimSpace(id)

	if id == "" {
		d.logger.Error("failed to get doctor",
			"error", ErrInvalidFields,
			"id", id)

		return nil, ErrInvalidFields
	}

	doctor, err := d.repo.GetByID(id)
	if err != nil {
		d.logger.Error("failed to get a doctor",
			"error", err,
			"id", id)

		return nil, ErrNotFound
	}

	return doctor, nil

}

func (d *DoctorUsecaseImpl) ListDoctors() ([]*model.Doctor, error) {
	doctors, err := d.repo.List()
	if err != nil {
		d.logger.Error("failed to list doctors",
			"error", err)

		return nil, err
	}

	return doctors, nil
}
