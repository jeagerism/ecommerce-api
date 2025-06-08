package repository

import (
	"github.com/jeagerism/ecommerce-api/domain"
	"github.com/jeagerism/ecommerce-api/entity"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) InsertUser(data *entity.User) error {
	if err := r.db.Create(data).Error; err != nil {
		return errors.Wrap(err, "[UserRepository.InsertUser]: failed to insert user into database")
	}
	return nil
}

func (r *userRepository) FindUserByEmail(email string) (*entity.User, error) {
	var user entity.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // User not found
		}
		return nil, errors.Wrap(err, "[UserRepository.FindUserByEmail]: failed to find user by email")
	}
	return &user, nil
}

func (r *userRepository) FindUserByID(id uint) (*entity.User, error) {
	var user entity.User
	if err := r.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil //user not found , no error
		}
		return nil, errors.Wrap(err, "[UserRepository.FindUserByID]: failed to find user by ID") // error in DB operation
	}
	return &user, nil
}

func (r *userRepository) FindUserAddressesByUserID(userID uint) ([]entity.UserAddress, error) {
	var addresses []entity.UserAddress
	if err := r.db.Where("user_id = ?", userID).Find(&addresses).Error; err != nil {
		return nil, errors.Wrap(err, "[UserRepository.FindUserAddressesByUserID]: failed to find user addresses by user ID")
	}
	return addresses, nil
}

func (r *userRepository) InsertUserAddress(address *entity.UserAddress) error {
	if err := r.db.Create(address).Error; err != nil {
		return errors.Wrap(err, "[UserRepository.InsertUserAddress]: failed to insert user address into database")
	}
	return nil
}
