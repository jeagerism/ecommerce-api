package usecase

import (
	"github.com/jeagerism/ecommerce-api/domain"
	"golang.org/x/crypto/bcrypt"
)

type shopUsecase struct {
	repo domain.ShopRepository
}

func NewShopUsecase(repo domain.ShopRepository) domain.ShopUsecase {
	return &shopUsecase{repo: repo}
}

// ฟังก์ชัน hash password
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (u *shopUsecase) CreateShop(input domain.CreatShop) error {
	hashedPassword, err := hashPassword(input.Password)
	if err != nil {
		return err
	}

	// สร้าง struct ใหม่ พร้อมใส่ password ที่ hash แล้ว
	secureShop := domain.CreatShop{
		Name:     input.Name,
		Email:    input.Email,
		Password: hashedPassword,
		Address:  input.Address,
	}

	// ส่งให้ repository บันทึกลง DB
	return u.repo.Create(secureShop)
}
