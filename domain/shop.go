package domain

import "github.com/jeagerism/ecommerce-api/entity"

type ShopRepository interface {
	InsertShop(*entity.Shop) error
	FindShopByEmail(string) (*entity.Shop, error)
	FindShopByID(id uint) (*entity.Shop, error)
	UpdateShop(*entity.Shop) error
	DeleteShop(id uint) error
}

type ShopUsecase interface {
	RegisterShop(*entity.CreateShop) error
	LoginShop(input *entity.ShopLoginInput) (string, error)
	GetShopDetails(uint) (*entity.ShopResponse, error)
	UpdateShop(uint, *entity.UpdateShop) error
	DeleteShop(uint) error
}
