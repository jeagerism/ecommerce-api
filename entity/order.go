package entity

type Order struct {
	ID         uint        `gorm:"primaryKey"`
	ShopID     uint        `gorm:"not null"`
	Shop       Shop        `gorm:"foreignKey:ShopID"`
	UserID     uint        `gorm:"not null"`
	User       User        `gorm:"foreignKey:UserID"`
	AddressID  uint        `gorm:"not null"`
	Address    UserAddress `gorm:"foreignKey:AddressID"`
	CourierID  uint        `gorm:"not null"`
	Courier    Courier     `gorm:"foreignKey:CourierID"`
	StatusID   uint        `gorm:"not null"`
	Status     OrderStatus `gorm:"foreignKey:StatusID"`
	TotalPrice float64     `gorm:"not null"`

	Items []OrderItem `gorm:"foreignKey:OrderID"` // <-- เพิ่มบรรทัดนี้
}

type OrderItem struct {
	ID        uint     `gorm:"primaryKey"`
	OrderID   uint     `gorm:"not null"`
	ProductID uint     `gorm:"not null"`
	Product   *Product `gorm:"foreignKey:ProductID"`
	Quantity  int      `gorm:"not null"`
}

type OrderStatus struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"not null"`
}

type Courier struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"not null"`
}

type CreateOrderRequest struct {
	ShopID    uint `json:"shop_id"`
	AddressID uint `json:"address_id"`
	CourierID uint `json:"courier_id"`
	Items     []struct {
		ProductID uint `json:"product_id"`
		Quantity  int  `json:"quantity"`
	} `json:"items"`
}

type OrderSummary struct {
	ID         uint    `json:"id"`
	ShopID     uint    `json:"shop_id"`
	UserID     uint    `json:"user_id"`
	Status     string  `json:"status"`
	TotalPrice float64 `json:"total_price"`
	CreatedAt  string  `json:"created_at"`
}

type Pagination struct {
	Page        int `json:"page"`
	PerPage     int `json:"per_page"`
	TotalPages  int `json:"total_pages"`
	TotalOrders int `json:"total_orders"`
}

type OrderListResponse struct {
	Orders     []OrderSummary `json:"orders"`
	Pagination Pagination     `json:"pagination"`
}

type OrderItemDetail struct {
	ProductID   uint    `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	TotalPrice  float64 `json:"total_price"`
}

type OrderDetailResponse struct {
	ID         uint              `json:"id"`
	ShopID     uint              `json:"shop_id"`
	UserID     uint              `json:"user_id"`
	AddressID  uint              `json:"address_id"`
	CourierID  uint              `json:"courier_id"`
	Status     string            `json:"status"`
	TotalPrice float64           `json:"total_price"`
	CreatedAt  string            `json:"created_at"`
	Items      []OrderItemDetail `json:"items"`
}

type UpdateOrderStatusRequest struct {
	OrderID  uint `json:"order_id"`
	StatusID uint `json:"status_id"`
}

type GetOrderStatus struct {
	ID       uint `gorm:"primaryKey"`
	StatusID uint
	Status   OrderStatus `gorm:"foreignKey:StatusID"`
}

type GetOrderStatusName struct {
	ID   uint
	Name string
}

type UpdateOrderStatusByUserRequest struct {
	OrderID  uint `json:"order_id"`
	ShopID   uint `json:"shop_id"`
	StatusID uint `json:"status_id"`
}
