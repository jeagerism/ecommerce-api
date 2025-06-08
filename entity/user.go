package entity

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Username string `json:"username" gorm:"unique;not null"`
	Email    string `json:"email" gorm:"unique;not null"`
	Password string `json:"-" gorm:"not null"` // Password should not be exposed in JSON responses
}
type CreateUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type UserAddress struct {
	ID            uint   `json:"id" gorm:"primaryKey"`
	UserID        uint   `json:"user_id" gorm:"not null;index"`
	RecipientName string `json:"recipient_name" gorm:"not null"`
	City          string `json:"city" gorm:"not null"`
	Zipcode       string `json:"zipcode" gorm:"not null"`
	Phone         string `json:"phone" gorm:"not null"`
}
type CreateAddress struct {
	RecipientName string `json:"recipient_name"`
	City          string `json:"city"`
	Zipcode       string `json:"zipcode"`
	Phone         string `json:"phone"`
}
type RegisterUserWithAddressInput struct {
	Username string        `json:"username"`
	Email    string        `json:"email"`
	Password string        `json:"password"`
	Address  CreateAddress `json:"address"`
}
