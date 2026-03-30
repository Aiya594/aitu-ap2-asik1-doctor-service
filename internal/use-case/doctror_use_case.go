package usecase

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/Aiya594/doctor-service/internal/model"
	"github.com/Aiya594/doctor-service/internal/repository"
	"github.com/google/uuid"
)

// Doctor Service Rules
//  full_name is required.
//  email is required.
//  email must be unique across all doctors.

type DoctorUsecase struct {
	repo   repository.DoctorRepository
	logger *slog.Logger
}

func (d *DoctorUsecase) CreateDoc(fullName, email, specialization string) error {

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

		return fmt.Errorf("full name, email and specialization are required:%w", ErrInvalidFields)
	}

	exists := d.repo.ExistsByEmail(email)

	if exists {
		d.logger.Error("failed create a doctor",
			"error", ErrAlreadyExists,
			"email", email)

		return fmt.Errorf("could not create a doctor:%w", ErrAlreadyExists)
	}

	id := uuid.New().String()

	doctor := &model.Doctor{
		ID:             id,
		FullName:       fullName,
		Email:          email,
		Specialization: specialization,
	}

	err := d.repo.Create(doctor)
	if err != nil {
		d.logger.Error("couldnt create a doctor",
			"error", err,
			"ID", doctor.ID,
			"full_name", doctor.FullName,
			"email", doctor.Email,
			"specialization", doctor.Specialization)
		return err
	}
	d.logger.Info("doctor created succesfully")

	return nil

}

func (d *DoctorUsecase) GetDocbyID(id string) (*model.Doctor, error) {
	id = strings.TrimSpace(id)

	if id == "" {
		d.logger.Error("failed to get doctor",
			"error", ErrInvalidFields,
			"id", id)

		return nil, fmt.Errorf("id is required: %w", ErrInvalidFields)
	}

	doctor, err := d.repo.GetByID(id)
	if err != nil {
		d.logger.Error("failed to get a doctor",
			"error", err,
			"id", id)

		return nil, err
	}

	return doctor, nil

}

func (d *DoctorUsecase) ListDoctors() ([]*model.Doctor, error) {
	doctors, err := d.repo.List()
	if err != nil {
		d.logger.Error("failed to list doctors",
			"error", err)

		return nil, err
	}

	return doctors, nil
}
