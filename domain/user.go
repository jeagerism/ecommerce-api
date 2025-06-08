package domain

import "github.com/jeagerism/ecommerce-api/entity"

type UserRepository interface {
	InsertUser(data *entity.User) error
	FindUserByEmail(email string) (*entity.User, error)
	FindUserByID(id uint) (*entity.User, error)
	InsertUserAddress(data *entity.UserAddress) error
	FindUserAddressesByUserID(userID uint) ([]entity.UserAddress, error)
}
type UserUsecase interface {
	CreateUser(input *entity.CreateUser) error
	Login(input *entity.LoginInput) (string, error)
	GetUserByID(id uint) (*entity.User, error)
	CreateUserAddress(userID uint, address *entity.CreateAddress) error
	GetUserAddressesByUserID(id uint) ([]entity.UserAddress, error)
}
