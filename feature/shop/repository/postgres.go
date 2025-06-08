package repository

import (
	"github.com/pkg/errors"

	"github.com/jeagerism/ecommerce-api/domain"
	"github.com/jeagerism/ecommerce-api/entity"
	"gorm.io/gorm"
)

type shopRepository struct {
	db *gorm.DB
}

func NewShopRepository(db *gorm.DB) domain.ShopRepository {
	return &shopRepository{db: db}
}

func (r *shopRepository) InsertShop(shop *entity.Shop) error {

	if err := r.db.Create(&shop).Error; err != nil {
		return errors.Wrap(err, "[ShopRepository.CreateShop]: failed to insert shop into database")
	}

	return nil
}

func (r *shopRepository) FindShopByEmail(email string) (*entity.Shop, error) {
	var shop entity.Shop
	if err := r.db.Where("email = ?", email).First(&shop).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "[ShopRepository.FindShopByEmail]: failed to find shop by email")
	}
	return &shop, nil
}

func (r *shopRepository) FindShopByID(id uint) (*entity.Shop, error) {
	var shop entity.Shop
	if err := r.db.First(&shop, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "[ShopRepository.FindByID] failed to find shop")
	}
	return &shop, nil
}

func (r *shopRepository) UpdateShop(shop *entity.Shop) error {
	if err := r.db.Model(&entity.Shop{}).
		Where("id = ?", shop.ID).
		Updates(map[string]interface{}{
			"name":    shop.Name,
			"email":   shop.Email,
			"address": shop.Address,
		}).Error; err != nil {
		return errors.Wrap(err, "[ShopRepository.UpdateShop]: failed to update shop")
	}
	return nil
}

func (r *shopRepository) DeleteShop(id uint) error {
	if err := r.db.Where("id = ?", id).Delete(&entity.Shop{}).Error; err != nil {
		return errors.Wrap(err, "[ShopRepository.DeleteShop]: failed to delete shop")
	}
	return nil
}
