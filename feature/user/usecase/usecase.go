package usecase

import (
	"github.com/jeagerism/ecommerce-api/domain"
	"github.com/jeagerism/ecommerce-api/entity"
	"github.com/jeagerism/ecommerce-api/util"
	"github.com/pkg/errors"
)

type userUsecase struct {
	repo      domain.UserRepository
	jwtSecret []byte
}

func NewUserUsecase(repo domain.UserRepository, jwtSecret []byte) domain.UserUsecase {
	return &userUsecase{
		repo:      repo,
		jwtSecret: jwtSecret,
	}
}

func (u *userUsecase) CreateUser(input *entity.CreateUser) error {
	// Hash the password before storing it
	hashedPassword, err := util.HashPassword(input.Password)
	if err != nil {
		return errors.Wrap(err, "[UserUsecase.CreateUser] failed to hash password")
	}

	user := &entity.User{
		Username: input.Username,
		Email:    input.Email,
		Password: hashedPassword,
	}

	// Insert the user into the repository
	if err := u.repo.InsertUser(user); err != nil {
		return errors.Wrap(err, "[UserUsecase.CreateUser] failed to insert user into database")
	}

	return nil
}

func (u *userUsecase) Login(input *entity.LoginInput) (string, error) {
	// Find user by email
	user, err := u.repo.FindUserByEmail(input.Email)
	if err != nil {
		return "", errors.Wrap(err, "[UserUsecase.Login] failed to find user by email")
	}

	// Check if user exists
	if user == nil {
		return "", errors.Wrap(errors.New("user not found"), "[UserUsecase.Login] invalid email or password")
	}

	// Check password
	if err := util.CheckPasswordHash(user.Password, input.Password); err != nil {
		return "", errors.Wrap(err, "[UserUsecase.Login] invalid email or password")
	}

	// Generate JWT token
	token, err := util.GenerateJWT(user.ID, user.Email, "user", u.jwtSecret)
	if err != nil {
		return "", errors.Wrap(err, "[UserUsecase.Login] failed to generate JWT")
	}

	return token, nil
}

func (u *userUsecase) GetUserByID(id uint) (*entity.User, error) {
	user, err := u.repo.FindUserByID(id)
	if err != nil {
		return nil, errors.Wrap(err, "[UserUsecase.GetUserByID] failed to find user by ID")
	}

	if user == nil {
		return nil, errors.Wrap(errors.New("user not found"), "[UserUsecase.GetUserByID] user not found")
	}

	return user, nil
}

func (u *userUsecase) CreateUserAddress(userID uint, address *entity.CreateAddress) error {
	// Create UserAddress entity
	userAddress := &entity.UserAddress{
		UserID:        userID,
		RecipientName: address.RecipientName,
		City:          address.City,
		Zipcode:       address.Zipcode,
		Phone:         address.Phone,
	}

	// Insert the user address into the repository
	if err := u.repo.InsertUserAddress(userAddress); err != nil {
		return errors.Wrap(err, "[UserUsecase.CreateUserAddress] failed to insert user address into database")
	}

	return nil
}

func (u *userUsecase) GetUserAddressesByUserID(id uint) ([]entity.UserAddress, error) {
	address, err := u.repo.FindUserAddressesByUserID(id)
	if err != nil {
		return nil, errors.Wrap(err, "[UserUsecase.GetUserAddressByID] failed to find user address by ID")
	}

	if address == nil {
		return nil, errors.Wrap(errors.New("address not found"), "[UserUsecase.GetUserAddressByID] address not found")
	}

	return address, nil
}
