package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/Aiya594/doctor-service/internal/model"
	"github.com/lib/pq"
)

type DoctorRepository interface {
	Create(doctor *model.Doctor, created_at time.Time) error
	GetByID(id string) (*model.Doctor, error)
	List() ([]*model.Doctor, error)
	ExistsByEmail(email string) bool
}

type PostgresDoctorRepository struct {
	db *sql.DB
}

func NewDoctorRepository(db *sql.DB) DoctorRepository {
	return &PostgresDoctorRepository{db: db}
}

func (r *PostgresDoctorRepository) Create(doctor *model.Doctor, created_at time.Time) error {
	query := `
		INSERT INTO doctors (id, full_name, specialization, email, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(
		query,
		doctor.ID,
		doctor.FullName,
		doctor.Specialization,
		doctor.Email,
		created_at,
	)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" {
				return ErrAlreadyExists
			}
		}
		return err
	}

	return nil
}
func (r *PostgresDoctorRepository) GetByID(id string) (*model.Doctor, error) {
	query := `
		SELECT id, full_name, specialization, email
		FROM doctors
		WHERE id = $1
	`

	var d model.Doctor

	err := r.db.QueryRow(query, id).Scan(
		&d.ID,
		&d.FullName,
		&d.Specialization,
		&d.Email,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &d, nil
}

func (r *PostgresDoctorRepository) List() ([]*model.Doctor, error) {
	query := `
		SELECT id, full_name, specialization, email
		FROM doctors
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var doctors []*model.Doctor

	for rows.Next() {
		var d model.Doctor

		err := rows.Scan(
			&d.ID,
			&d.FullName,
			&d.Specialization,
			&d.Email,
		)
		if err != nil {
			return nil, err
		}

		doctors = append(doctors, &d)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return doctors, nil
}

func (r *PostgresDoctorRepository) ExistsByEmail(email string) bool {
	query := `
		SELECT 1 FROM doctors WHERE email = $1 LIMIT 1
	`

	var exists int
	err := r.db.QueryRow(query, email).Scan(&exists)

	return err == nil
}
