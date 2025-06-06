package domain

type Shop struct {
	ID       uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Name     string `json:"name" gorm:"not null"`
	Email    string `json:"email" gorm:"unique;not null"`
	Password string `json:"-" gorm:"not null"` // ซ่อน password จาก JSON response ยังใช้ใน goได้
	Address  string `json:"address" gorm:"not null"`
}

type CreatShop struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Address  string `json:"address"`
}

type ShopRepository interface {
	Create(CreatShop) error
}

type ShopUsecase interface {
	CreateShop(CreatShop) error
}
