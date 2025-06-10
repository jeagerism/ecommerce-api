package repository

import (
	"github.com/jeagerism/ecommerce-api/domain"
	"github.com/jeagerism/ecommerce-api/entity"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) domain.OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) InsertCourier(courier *entity.Courier) error {
	if err := r.db.Create(courier).Error; err != nil {
		return errors.Wrap(err, "[OrderRepository.InsertCourier]: failed to insert courier into database")
	}
	return nil
}

func (r *orderRepository) FindCouriers() ([]entity.Courier, error) {
	var couriers []entity.Courier
	if err := r.db.Find(&couriers).Error; err != nil {
		return nil, errors.Wrap(err, "[OrderRepository.FindCouriers]: failed to find couriers")
	}
	return couriers, nil
}

func (r *orderRepository) GetProductByID(productID uint) (*entity.Product, error) {
	var product entity.Product
	if err := r.db.First(&product, productID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.Wrapf(err, "product id %d not found", productID)
		}
		return nil, errors.Wrap(err, "[OrderRepository.GetProductByID]: failed to get product by ID")
	}
	return &product, nil
}
func (r *orderRepository) InsertOrderWithItems(order *entity.Order, items []entity.OrderItem) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Insert Order
		if err := tx.Create(order).Error; err != nil {
			return errors.Wrap(err, "[OrderRepository.InsertOrderWithItems]: failed to insert order")
		}

		// Insert OrderItems
		for i := range items {
			items[i].OrderID = order.ID
		}

		if err := tx.Create(&items).Error; err != nil {
			return errors.Wrap(err, "[OrderRepository.InsertOrderWithItems]: failed to insert order items")
		}

		return nil // success → commit
	})
}

func (r *orderRepository) InsertStatus(status *entity.OrderStatus) error {
	if err := r.db.Create(status).Error; err != nil {
		return errors.Wrap(err, "[OrderRepository.InsertStatus]: failed to insert order status into database")
	}
	return nil
}

func (r *orderRepository) FindOrderByID(orderID, userID uint) (*entity.Order, error) {
	var order entity.Order
	if err := r.db.
		Preload("Shop").
		Preload("User").
		Preload("Address").
		Preload("Courier").
		Preload("Status").
		Preload("Items.Product").
		Where("id = ? AND user_id = ?", orderID, userID). // ✅ เช็กว่า user นี้เป็นเจ้าของ order
		First(&order).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.Wrapf(err, "order id %d not found for user id %d", orderID, userID)
		}
		return nil, errors.Wrap(err, "[OrderRepository.FindOrderByID]")
	}

	return &order, nil
}

func (r *orderRepository) UpdateOrderStatusByID(orderID, statusID uint) error {
	if err := r.db.Model(&entity.Order{}).Where("id = ?", orderID).Update("status_id", statusID).Error; err != nil {
		return errors.Wrapf(err, "[OrderRepository.UpdateOrderStatusByID]: failed to update status for order id %d", orderID)
	}
	return nil
}

func (r *orderRepository) FindStatusByID(statusID uint) (*entity.OrderStatus, error) {
	var status entity.OrderStatus
	if err := r.db.First(&status, statusID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.Wrapf(err, "[OrderRepository.FindStatusByID]: order status id %d not found", statusID)
		}
		return nil, errors.Wrap(err, "[OrderRepository.FindStatusByID]: failed to find order status by ID")
	}
	return &status, nil
}

func (r *orderRepository) FindOrdersByShopID(shopID uint) ([]entity.Order, error) {
	var orders []entity.Order

	if err := r.db.
		Preload("Shop").
		Preload("User").
		Preload("Address").
		Preload("Courier").
		Preload("Status").
		Preload("Items.Product").
		Where("shop_id = ?", shopID).
		Find(&orders).Error; err != nil {

		return nil, errors.Wrap(err, "[OrderRepository.FindOrdersByShopID]")
	}

	return orders, nil
}

func (r *orderRepository) UpdateOrderStatusByUser(orderID uint) error {
	if err := r.db.Model(&entity.Order{}).
		Where("id = ?", orderID). // ✅ เช็กว่า user นี้เป็นเจ้าของ order
		Update("status_id", 5).Error; err != nil {
		return errors.Wrapf(err, "[OrderRepository.UpdateOrderStatusByUser]: failed to update status for order id %d", orderID)
	}
	return nil
}

func (r *orderRepository) FindOrderStatusInfo(orderID, shopID, userID uint) (entity.GetOrderStatus, error) {
	var order entity.GetOrderStatus

	if err := r.db.
		Model(&entity.Order{}).
		Select("orders.id, orders.status_id, order_statuses.name").
		Joins("JOIN order_statuses ON orders.status_id = order_statuses.id").
		Where("orders.id = ? AND orders.shop_id = ? AND orders.user_id = ?", orderID, shopID, userID).
		First(&order).Error; err != nil {
		return order, err
	}

	return order, nil
}
