package repository

import (
	"ordersystem/model"
	"time"
)

type DatabaseHandler struct {
	// drinks represent all available drinks
	drinks []model.Drink
	// orders serves as order history
	orders []model.Order
}

func NewDatabaseHandler() *DatabaseHandler {
	// Init the drinks slice with some test data
	drinks := []model.Drink{
		{ID: 1, Name: "Beer", Price: 2.00, Description: "Hagenberger Gold"},
		{ID: 2, Name: "Spritzer", Price: 1.40, Description: "Wine with soda"},
		{ID: 3, Name: "Coffee", Price: 1.50, Description: "Freshly brewed"},
	}

	// Init orders slice with some test data
	orders := []model.Order{
		{DrinkID: 1, CreatedAt: time.Now().Add(-30 * time.Minute), Amount: 2},
		{DrinkID: 2, CreatedAt: time.Now().Add(-20 * time.Minute), Amount: 1},
		{DrinkID: 1, CreatedAt: time.Now().Add(-10 * time.Minute), Amount: 3},
		{DrinkID: 3, CreatedAt: time.Now(), Amount: 1},
	}

	return &DatabaseHandler{
		drinks: drinks,
		orders: orders,
	}
}

func (db *DatabaseHandler) GetDrinks() []model.Drink {
	return db.drinks
}

func (db *DatabaseHandler) GetOrders() []model.Order {
	return db.orders
}

func (db *DatabaseHandler) GetTotalledOrders() map[uint64]uint64 {
	// calculate total orders
	// key = DrinkID, value = Amount of orders
	totalledOrders := make(map[uint64]uint64)
	for _, order := range db.orders {
		totalledOrders[order.DrinkID] += order.Amount
	}
	return totalledOrders
}

func (db *DatabaseHandler) AddOrder(order *model.Order) {
	// add order to db.orders slice
	db.orders = append(db.orders, *order)
}
