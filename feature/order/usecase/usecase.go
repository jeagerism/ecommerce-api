package usecase

import (
	"fmt"

	"github.com/jeagerism/ecommerce-api/domain"
	"github.com/jeagerism/ecommerce-api/entity"
	"github.com/pkg/errors"
)

type orderUsecase struct {
	repo domain.OrderRepository
}

func NewOrderUsecase(repo domain.OrderRepository) domain.OrderUsecase {
	return &orderUsecase{
		repo: repo,
	}
}

func (u *orderUsecase) CreateOrder(userID uint, req *entity.CreateOrderRequest) error {
	const pendingStatusID = 1 // สมมติว่า ID=1 คือ Pending
	order := &entity.Order{
		ShopID:    req.ShopID,
		UserID:    userID,
		AddressID: req.AddressID,
		CourierID: req.CourierID,
		StatusID:  pendingStatusID,
	}

	var items []entity.OrderItem
	var totalPrice float64

	for _, i := range req.Items {
		if i.Quantity <= 0 {
			return errors.Errorf("invalid quantity for product ID %d", i.ProductID)
		}

		product, err := u.repo.GetProductByID(i.ProductID)
		if err != nil {
			return errors.Wrapf(err, "product ID %d not found", i.ProductID)
		}

		items = append(items, entity.OrderItem{
			ProductID: i.ProductID,
			Quantity:  i.Quantity,
		})

		totalPrice += product.Price * float64(i.Quantity)
	}

	order.TotalPrice = totalPrice

	// Insert ทั้ง order และ order items ใน transaction
	if err := u.repo.InsertOrderWithItems(order, items); err != nil {
		return errors.Wrap(err, "[OrderUsecase.CreateOrderFromRequest] failed to insert order and items")
	}

	return nil
}

func (u *orderUsecase) FindCouriers() ([]entity.Courier, error) {
	couriers, err := u.repo.FindCouriers()
	if err != nil {
		return nil, errors.Wrap(err, "[OrderUsecase.FindCouriers] failed to find couriers")
	}
	return couriers, nil
}

func (u *orderUsecase) CreateCourier(courier *entity.Courier) error {
	if err := u.repo.InsertCourier(courier); err != nil {
		return errors.Wrap(err, "[OrderUsecase.CreateCourier] failed to create courier")
	}
	return nil
}

func (u *orderUsecase) CreateStatus(status *entity.OrderStatus) error {
	if err := u.repo.InsertStatus(status); err != nil {
		return errors.Wrap(err, "[OrderUsecase.CreateStatus] failed to create order status")
	}
	return nil
}

func (u *orderUsecase) FindOrderByID(orderID, userID uint) (*entity.OrderDetailResponse, error) {
	order, err := u.repo.FindOrderByID(orderID, userID) // ควร preload Items และ Product มาด้วยใน repo
	if err != nil {
		return nil, err
	}

	var itemDetails []entity.OrderItemDetail
	for _, item := range order.Items {
		// ตรวจสอบว่ามีข้อมูล product
		if item.Product == nil {
			return nil, fmt.Errorf("product not found for product id %d", item.ProductID)
		}
		unitPrice := item.Product.Price
		itemDetails = append(itemDetails, entity.OrderItemDetail{
			ProductID:   item.ProductID,
			ProductName: item.Product.Name,
			Quantity:    item.Quantity,
			UnitPrice:   unitPrice,
			TotalPrice:  unitPrice * float64(item.Quantity),
		})
	}

	response := &entity.OrderDetailResponse{
		ID:         order.ID,
		ShopID:     order.ShopID,
		UserID:     order.UserID,
		AddressID:  order.AddressID,
		CourierID:  order.CourierID,
		Status:     order.Status.Name, // preload มาแล้ว
		TotalPrice: order.TotalPrice,
		Items:      itemDetails,
	}

	return response, nil
}

func (u *orderUsecase) UpdateOrderStatus(orderID, statusID uint) error {
	// ตรวจสอบว่า statusID มีอยู่จริง
	_, err := u.repo.FindStatusByID(statusID)
	if err != nil {
		return errors.Wrapf(err, "[OrderUsecase.UpdateOrderStatus] status ID %d not found", statusID)
	}

	// ตรวจสอบว่า orderID มีอยู่จริง
	_, err = u.repo.FindOrderByID(orderID, 0)
	if err != nil {
		return errors.Wrapf(err, "[OrderUsecase.UpdateOrderStatus] order ID %d not found", orderID)
	}

	// อัปเดตสถานะคำสั่งซื้อ
	if err := u.repo.UpdateOrderStatusByID(orderID, statusID); err != nil {
		return errors.Wrapf(err, "[OrderUsecase.UpdateOrderStatus] failed to update order ID %d to status ID %d", orderID, statusID)
	}

	return nil
}

func (u *orderUsecase) FindOrdersByShopID(shopID uint) ([]entity.OrderDetailResponse, error) {
	orders, err := u.repo.FindOrdersByShopID(shopID)
	if err != nil {
		return nil, errors.Wrapf(err, "[OrderUsecase.FindOrdersByShopID] failed to find orders for shop ID %d", shopID)
	}

	if len(orders) == 0 {
		return []entity.OrderDetailResponse{}, nil // คืนค่า empty slice แทน nil
	}

	responses := make([]entity.OrderDetailResponse, 0, len(orders))

	for _, order := range orders {
		var itemDetails []entity.OrderItemDetail
		// loop set items in order
		for _, item := range order.Items {
			if item.Product == nil {
				return nil, fmt.Errorf("product not found for product id %d", item.ProductID)
			}

			unitPrice := item.Product.Price
			itemDetails = append(itemDetails, entity.OrderItemDetail{
				ProductID:   item.ProductID,
				ProductName: item.Product.Name,
				Quantity:    item.Quantity,
				UnitPrice:   unitPrice,
				TotalPrice:  unitPrice * float64(item.Quantity),
			})
		}
		// set order detail response
		responses = append(responses, entity.OrderDetailResponse{
			ID:         order.ID,
			ShopID:     order.ShopID,
			UserID:     order.UserID,
			AddressID:  order.AddressID,
			CourierID:  order.CourierID,
			Status:     order.Status.Name,
			TotalPrice: order.TotalPrice,
			Items:      itemDetails, // set items in order
		})
	}

	return responses, nil
}

func (u *orderUsecase) UpdateOrderStatusByUser(orderReq entity.UpdateOrderStatusByUserRequest, userID uint) error {
	// ตรวจสอบว่า orderID มีอยู่จริง
	order, err := u.repo.FindOrderStatusInfo(orderReq.OrderID, orderReq.ShopID, userID)
	if err != nil {
		return errors.Wrapf(err, "[OrderUsecase.UpdateOrderStatusByUser] order ID %d not found for user ID %d", orderReq.OrderID, userID)
	}
	if order.Status.ID != orderReq.StatusID {
		return errors.Errorf("[OrderUsecase.UpdateOrderStatusByUser] order ID %d does not have status ID %d", orderReq.OrderID, orderReq.StatusID)
	}
	if order.Status.Name != "cancelled" && order.Status.Name != "delivered" {
		if err := u.repo.UpdateOrderStatusByUser(orderReq.OrderID); err != nil {
			return errors.Wrapf(err, "[OrderUsecase.UpdateOrderStatusByUser] failed to update order ID %d to cancelled status", orderReq.OrderID)
		}
		return nil
	}

	return nil
}
