package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/Aiya594/doctor-service/internal/model"
	"github.com/redis/go-redis/v9"
)

type CacheRepository interface {
	GetDoctor(ctx context.Context, id string) (*model.Doctor, error)
	SetDoctor(ctx context.Context, doctor *model.Doctor) error
	GetDoctorList(ctx context.Context) ([]*model.Doctor, error)
	SetDoctorList(ctx context.Context, doctors []*model.Doctor) error
	InvalidateDoctorList(ctx context.Context) error
}

type RedisCacheRepository struct {
	client *redis.Client
	ttl    time.Duration
	logger *slog.Logger
}

func NewRedisCacheRepository(client *redis.Client, logger *slog.Logger) CacheRepository {
	ttlSec := 60
	if v := os.Getenv("CACHE_TTL_SECONDS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			ttlSec = n
		}
	}
	return &RedisCacheRepository{
		client: client,
		ttl:    time.Duration(ttlSec) * time.Second,
		logger: logger,
	}
}

func doctorKey(id string) string {
	return fmt.Sprintf("doctor:%s", id)
}

const doctorListKey = "doctors:list"

func (r *RedisCacheRepository) GetDoctor(ctx context.Context, id string) (*model.Doctor, error) {
	val, err := r.client.Get(ctx, doctorKey(id)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil // cache miss
		}
		return nil, err
	}
	var doc model.Doctor
	if err := json.Unmarshal([]byte(val), &doc); err != nil {
		return nil, err
	}
	return &doc, nil
}

func (r *RedisCacheRepository) SetDoctor(ctx context.Context, doctor *model.Doctor) error {
	data, err := json.Marshal(doctor)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, doctorKey(doctor.ID), data, r.ttl).Err()
}

func (r *RedisCacheRepository) GetDoctorList(ctx context.Context) ([]*model.Doctor, error) {
	val, err := r.client.Get(ctx, doctorListKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, err
	}
	var docs []*model.Doctor
	if err := json.Unmarshal([]byte(val), &docs); err != nil {
		return nil, err
	}
	return docs, nil
}

func (r *RedisCacheRepository) SetDoctorList(ctx context.Context, doctors []*model.Doctor) error {
	data, err := json.Marshal(doctors)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, doctorListKey, data, r.ttl).Err()
}

func (r *RedisCacheRepository) InvalidateDoctorList(ctx context.Context) error {
	return r.client.Del(ctx, doctorListKey).Err()
}
