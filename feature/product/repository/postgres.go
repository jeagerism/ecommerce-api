package repository

import (
	"github.com/jeagerism/ecommerce-api/entity"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *productRepository {
	return &productRepository{db: db}
}

func (r *productRepository) InsertProduct(data *entity.Product) error {

	if err := r.db.Create(data).Error; err != nil {
		return errors.Wrap(err, "[ProductRepository.InsertProduct]: failed to insert product into database")
	}

	return nil
}

func (r *productRepository) FindAllProductsByShopID(shopID uint) ([]entity.Product, error) {
	var products []entity.Product
	if err := r.db.Where("shop_id = ?", shopID).Find(&products).Error; err != nil {
		return nil, errors.Wrap(err, "[ProductRepository.FindAllProductsByShopID]: failed to retrieve products by shop ID")
	}
	return products, nil
}

func (r *productRepository) UpdateProduct(data *entity.Product) error {
	if err := r.db.Save(data).Error; err != nil {
		return errors.Wrap(err, "[ProductRepository.UpdateProduct]: failed to update product in database")
	}
	return nil
}

func (r *productRepository) DeleteProduct(shopID, productID uint) error {
	if err := r.db.Where("shop_id = ? AND id = ?", shopID, productID).Delete(&entity.Product{}).Error; err != nil {
		return errors.Wrap(err, "[ProductRepository.DeleteProductByShopIDAndProductID]: failed to delete product from database")
	}
	return nil
}
