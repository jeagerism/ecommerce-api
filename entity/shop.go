package entity

type Shop struct {
	ID       uint   `gorm:"primaryKey;autoIncrement"`
	Name     string `json:"name" gorm:"type:varchar(50);not null"`
	Email    string `json:"email" gorm:"type:varchar(100);unique;not null"`
	Password string `json:"password" gorm:"type:varchar(100);not null"`
	Address  string `json:"address" gorm:"type:text;not null"`
}

type ShopResponse struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Address string `json:"address"`
}

type CreateShop struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Address  string `json:"address"`
}

type ShopLoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UpdateShop struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Address string `json:"address"`
}
