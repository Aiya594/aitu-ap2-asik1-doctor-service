package cache

import (
	"context"

	"github.com/Aiya594/doctor-service/internal/model"
)

// NoopCacheRepository is used when Redis is unavailable. All reads miss and all writes are silent no-ops.
type NoopCacheRepository struct{}

func NewNoop() CacheRepository { return &NoopCacheRepository{} }

func (n *NoopCacheRepository) GetDoctor(_ context.Context, _ string) (*model.Doctor, error) {
	return nil, nil
}
func (n *NoopCacheRepository) SetDoctor(_ context.Context, _ *model.Doctor) error { return nil }
func (n *NoopCacheRepository) GetDoctorList(_ context.Context) ([]*model.Doctor, error) {
	return nil, nil
}
func (n *NoopCacheRepository) SetDoctorList(_ context.Context, _ []*model.Doctor) error { return nil }
func (n *NoopCacheRepository) InvalidateDoctorList(_ context.Context) error             { return nil }
