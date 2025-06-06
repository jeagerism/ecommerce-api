package repository

import (
	"github.com/jeagerism/ecommerce-api/domain"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type shopRepository struct {
	db *gorm.DB
}

func NewShopRepository(db *gorm.DB) domain.ShopRepository {
	return &shopRepository{db: db}
}

func (r *shopRepository) Create(data domain.CreatShop) error {
	shop := domain.Shop{
		Name:     data.Name,
		Email:    data.Email,
		Password: data.Password, // ต้อง hash มาแล้ว
		Address:  data.Address,
	}

	if err := r.db.Create(&shop).Error; err != nil {
		logrus.Errorln("[Repository Shop : failed to create new shop]")
		return err
	}

	return nil
}
