package domain

import "github.com/jeagerism/ecommerce-api/entity"

type ProductRepository interface {
	InsertProduct(*entity.Product) error
	FindAllProductsByShopID(uint) ([]entity.Product, error)
	UpdateProduct(*entity.Product) error
	DeleteProduct(uint, uint) error
}

type ProductUsecase interface {
	CreateProduct(uint, *entity.CreateProduct) error
	GetAllProductsByShopID(uint) ([]entity.Product, error)
	UpdateProduct(uint, uint, *entity.UpdateProduct) error
	DeleteProduct(uint, uint) error
}
