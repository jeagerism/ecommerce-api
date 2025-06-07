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

func (r *shopRepository) InsertShop(data *entity.CreateShop) error {
	shop := entity.Shop{
		Name:     data.Name,
		Email:    data.Email,
		Password: data.Password,
		Address:  data.Address,
	}

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

func (r *shopRepository) FindShopByID(id uint) (*entity.GetShop, error) {
	var shop entity.GetShop
	if err := r.db.Where("id = ?", id).First(&shop).Error; err != nil { // Explicitly specify the table name
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		// if ids instead of id make this error bro!
		return nil, errors.Wrap(err, "[ShopRepository.FindByID]: failed to find shop by ID")
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
