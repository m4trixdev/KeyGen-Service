package repository

import (
	"github.com/google/uuid"
	"github.com/m4trixdev/keygen-service/internal/models"
	"gorm.io/gorm"
)

type KeyRepository struct {
	db *gorm.DB
}

func NewKeyRepository(db *gorm.DB) *KeyRepository {
	return &KeyRepository{db: db}
}

func (r *KeyRepository) Create(key *models.Key) error {
	return r.db.Create(key).Error
}

func (r *KeyRepository) FindByValue(value string) (*models.Key, error) {
	var key models.Key
	err := r.db.Where("value = ?", value).First(&key).Error
	if err != nil {
		return nil, err
	}
	return &key, nil
}

func (r *KeyRepository) FindByID(id uuid.UUID) (*models.Key, error) {
	var key models.Key
	err := r.db.First(&key, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &key, nil
}

func (r *KeyRepository) FindAll(page, size int) ([]models.Key, int64, error) {
	var keys []models.Key
	var total int64

	r.db.Model(&models.Key{}).Count(&total)

	offset := (page - 1) * size
	err := r.db.Offset(offset).Limit(size).Order("created_at desc").Find(&keys).Error
	return keys, total, err
}

func (r *KeyRepository) IncrementUses(id uuid.UUID) error {
	return r.db.Model(&models.Key{}).Where("id = ?", id).UpdateColumn("uses", gorm.Expr("uses + 1")).Error
}

func (r *KeyRepository) Revoke(id uuid.UUID) error {
	return r.db.Model(&models.Key{}).Where("id = ?", id).Update("revoked", true).Error
}

func (r *KeyRepository) LogUsage(log *models.KeyUsageLog) error {
	return r.db.Create(log).Error
}
