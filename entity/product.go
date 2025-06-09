package entity

type Product struct {
	ID          uint    `gorm:"primaryKey;autoIncrement"`
	ShopID      uint    `json:"shop_id" gorm:"not null"` // Foreign key to the shop
	Name        string  `json:"name" gorm:"type:varchar(100);not null"`
	Description string  `json:"description" gorm:"type:text;not null"`
	Price       float64 `json:"price" gorm:"type:decimal(10,2);not null"`
}

type CreateProduct struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type UpdateProduct struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}
