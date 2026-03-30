package repository

import (
	"sync"

	"github.com/Aiya594/doctor-service/internal/model"
)

type DoctorRepository interface {
	Create(doctor *model.Doctor) error
	GetByID(id string) (*model.Doctor, error)
	List() ([]*model.Doctor, error)
	ExistsByEmail(email string) bool
}

type InMemoryDoctorRepository struct {
	mu      sync.RWMutex
	doctors map[string]*model.Doctor
}

func NewDocRepo() DoctorRepository {
	return &InMemoryDoctorRepository{
		doctors: make(map[string]*model.Doctor),
	}
}

func (i *InMemoryDoctorRepository) Create(doctor *model.Doctor) error {

	i.mu.Lock()
	defer i.mu.Unlock()

	i.doctors[doctor.ID] = doctor

	return nil
}

func (i *InMemoryDoctorRepository) GetByID(id string) (*model.Doctor, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	doctor, exists := i.doctors[id]
	if !exists {
		return nil, ErrNotFound
	}
	return doctor, nil
}

func (i *InMemoryDoctorRepository) List() ([]*model.Doctor, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()
	docs := make([]*model.Doctor, 0, len(i.doctors))
	for _, d := range i.doctors {
		docs = append(docs, d)
	}
	return docs, nil
}

func (i *InMemoryDoctorRepository) ExistsByEmail(email string) bool {
	i.mu.RLock()
	defer i.mu.RUnlock()

	for _, d := range i.doctors {
		if d.Email == email {
			return true
		}
	}

	return false
}
