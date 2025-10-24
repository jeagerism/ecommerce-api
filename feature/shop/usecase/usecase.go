package usecase

import (
	"os"

	"github.com/jeagerism/ecommerce-api/domain"
	"github.com/jeagerism/ecommerce-api/entity"
	"github.com/jeagerism/ecommerce-api/util"
	"github.com/pkg/errors"
)

type shopUsecase struct {
	repo domain.ShopRepository
}

func NewShopUsecase(r domain.ShopRepository) domain.ShopUsecase {
	return &shopUsecase{
		repo: r,
	}
}

func (u *shopUsecase) RegisterShop(input *entity.CreateShop) error {
	hashedPassword, err := util.HashPassword(input.Password)
	if err != nil {
		return errors.Wrap(err, "[ShopUsecase.CreateShop] failed to hash password")
	}

	secureShop := entity.Shop{
		Name:     input.Name,
		Email:    input.Email,
		Password: hashedPassword,
		Address:  input.Address,
	}

	if err := u.repo.InsertShop(&secureShop); err != nil {
		return errors.Wrap(err, "[ShopUsecase.CreateShop] failed to create shop in repository")
	}

	return nil
}

func (u *shopUsecase) LoginShop(input *entity.ShopLoginInput) (string, error) {
	shop, err := u.repo.FindShopByEmail(input.Email)
	if err != nil {
		return "", errors.Wrap(err, "[ShopUsecase.LoginShop] failed to find shop by email")
	}
	if shop == nil {
		return "", errors.New("[ShopUsecase.LoginShop] invalid email or password")
	}

	if err := util.CheckPasswordHash(shop.Password, input.Password); err != nil {
		return "", errors.New("[ShopUsecase.LoginShop] invalid email or password")
	}

	shopSecret := []byte(os.Getenv("SHOP_SECRET")) // Replace with your actual secret key
	token, err := util.GenerateJWT(shop.ID, "shop", shopSecret)
	if err != nil {
		return "", errors.Wrap(err, "[ShopUsecase.LoginShop] failed to generate JWT")
	}
	return token, nil
}

func (u *shopUsecase) GetShopDetails(shopId uint) (*entity.ShopResponse, error) {
	shop, err := u.repo.FindShopByID(shopId)
	if err != nil {
		return nil, errors.Wrap(err, "[ShopUsecase.GetShopDetails] failed to find shop by ID")
	}
	if shop == nil {
		return nil, nil
	}

	return &entity.ShopResponse{
		ID:      shop.ID,
		Name:    shop.Name,
		Email:   shop.Email,
		Address: shop.Address,
	}, nil
}

func (u *shopUsecase) UpdateShop(shopId uint, input *entity.UpdateShop) error {
	// Find the shop by ID
	shop, err := u.repo.FindShopByID(shopId)
	if err != nil {
		return errors.Wrap(err, "[ShopUsecase.UpdateShop] failed to find shop by ID")
	}
	if shop == nil {
		return errors.Wrap(err, "[ShopUsecase.UpdateShop] shop not found")
	}

	// Create a new entity.Shop object for updating
	updatedShop := &entity.Shop{
		ID:      shop.ID,
		Name:    input.Name,
		Email:   input.Email,
		Address: input.Address,
	}

	// Call the repository to update the shop
	if err := u.repo.UpdateShop(updatedShop); err != nil {
		return errors.Wrap(err, "[ShopUsecase.UpdateShop] failed to update shop")
	}

	return nil
}

func (u *shopUsecase) DeleteShop(shopId uint) error {
	// Call the repository to delete the shop
	if err := u.repo.DeleteShop(shopId); err != nil {
		return errors.Wrap(err, "[ShopUsecase.DeleteShop] failed to delete shop")
	}
	return nil
}
