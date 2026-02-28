package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/m4trixdev/keygen-service/internal/models"
	"github.com/m4trixdev/keygen-service/internal/repository"
)

type CreateKeyInput struct {
	Label     string     `json:"label"`
	MaxUses   int        `json:"max_uses"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	UserID    *uuid.UUID `json:"user_id,omitempty"`
}

type KeyService struct {
	repo *repository.KeyRepository
}

func NewKeyService(repo *repository.KeyRepository) *KeyService {
	return &KeyService{repo: repo}
}

func (s *KeyService) Generate(input CreateKeyInput) (*models.Key, error) {
	value, err := generateKeyValue()
	if err != nil {
		return nil, fmt.Errorf("failed to generate key value: %w", err)
	}

	maxUses := input.MaxUses
	if maxUses == 0 {
		maxUses = -1
	}

	key := &models.Key{
		Value:     value,
		Label:     input.Label,
		MaxUses:   maxUses,
		ExpiresAt: input.ExpiresAt,
		UserID:    input.UserID,
	}

	if err := s.repo.Create(key); err != nil {
		return nil, err
	}

	return key, nil
}

func (s *KeyService) Validate(value, ip string) (*models.Key, error) {
	key, err := s.repo.FindByValue(value)
	if err != nil {
		return nil, errors.New("key not found")
	}

	if !key.IsUsable() {
		if key.Revoked {
			return nil, errors.New("key has been revoked")
		}
		if key.IsExpired() {
			return nil, errors.New("key has expired")
		}
		return nil, errors.New("key usage limit reached")
	}

	if err := s.repo.IncrementUses(key.ID); err != nil {
		return nil, err
	}

	_ = s.repo.LogUsage(&models.KeyUsageLog{
		KeyID: key.ID,
		IP:    ip,
	})

	return key, nil
}

func (s *KeyService) Revoke(id uuid.UUID) error {
	key, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("key not found")
	}
	if key.Revoked {
		return errors.New("key is already revoked")
	}
	return s.repo.Revoke(id)
}

func (s *KeyService) List(page, size int) (map[string]any, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}

	keys, total, err := s.repo.FindAll(page, size)
	if err != nil {
		return nil, err
	}

	pages := int(math.Ceil(float64(total) / float64(size)))
	if pages < 1 {
		pages = 1
	}

	return map[string]any{
		"items": keys,
		"total": total,
		"page":  page,
		"size":  size,
		"pages": pages,
	}, nil
}

func generateKeyValue() (string, error) {
	b := make([]byte, 20)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	raw := hex.EncodeToString(b)
	return fmt.Sprintf("%s-%s-%s-%s", raw[0:8], raw[8:16], raw[16:24], raw[24:32]), nil
}
