package usecase

import (
	"github.com/jeagerism/ecommerce-api/domain"
	"github.com/jeagerism/ecommerce-api/entity"
	"github.com/pkg/errors"
)

type productUsecase struct {
	productRepo domain.ProductRepository
	jwtSecret   []byte
}

func NewProductUsecase(r domain.ProductRepository, secret []byte) domain.ProductUsecase {
	return &productUsecase{
		productRepo: r,
		jwtSecret:   secret,
	}
}

// CreateProduct creates a new product
func (u *productUsecase) CreateProduct(shopID uint, input *entity.CreateProduct) error {
	product := &entity.Product{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		ShopID:      shopID, // Associate the product with the shop
	}

	if err := u.productRepo.InsertProduct(product); err != nil {
		return errors.Wrap(err, "[ProductUsecase.CreateProduct] failed to insert product into database")
	}

	return nil
}

func (u *productUsecase) GetAllProductsByShopID(shopID uint) ([]*entity.Product, error) {
	products, err := u.productRepo.FindAllProductsByShopID(shopID)
	if err != nil {
		return nil, errors.Wrap(err, "[ProductUsecase.GetAllProductsByShopID] failed to retrieve products by shop ID")
	}

	return products, nil
}

func (u *productUsecase) UpdateProduct(shopID uint, productID uint, input *entity.UpdateProduct) error {
	product := &entity.Product{
		ID:          productID,
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		ShopID:      shopID, // Ensure the shop ID is set
	}

	if err := u.productRepo.UpdateProduct(product); err != nil {
		return errors.Wrap(err, "[ProductUsecase.UpdateProduct] failed to update product in database")
	}

	return nil
}

func (u *productUsecase) DeleteProduct(shopID uint, productID uint) error {
	if err := u.productRepo.DeleteProduct(shopID, productID); err != nil {
		return errors.Wrap(err, "[ProductUsecase.DeleteProduct] failed to delete product from database")
	}

	return nil
}
