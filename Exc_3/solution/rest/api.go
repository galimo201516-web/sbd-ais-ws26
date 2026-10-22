package rest

import (
	"encoding/json"
	"net/http"
	"ordersystem/model"
	"ordersystem/repository"
	"time"

	"github.com/go-chi/render"
)

// GetMenu 			godoc
// @tags 			Menu
// @Description 	Returns the menu of all drinks
// @Produce  		json
// @Success 		200 {array} model.Drink
// @Router 			/api/menu [get]
func GetMenu(db *repository.DatabaseHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get slice from db
		drinks := db.GetDrinks()
		render.Status(r, http.StatusOK)
		render.JSON(w, r, drinks)
	}
}

// GetOrders 		godoc
// @tags 			Order
// @Description 	Returns all orders
// @Produce  		json
// @Success 		200 {array} model.Order
// @Router 			/api/order/all [get]
func GetOrders(db *repository.DatabaseHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get slice from db
		orders := db.GetOrders()
		render.Status(r, http.StatusOK)
		render.JSON(w, r, orders)
	}
}

// GetOrdersTotal 	godoc
// @tags 			Order
// @Description 	Returns how many of each drink have been ordered (key = drink ID)
// @Produce  		json
// @Success 		200 {object} map[string]uint64
// @Router 			/api/order/totalled [get]
func GetOrdersTotal(db *repository.DatabaseHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get map from db
		totalledOrders := db.GetTotalledOrders()
		render.Status(r, http.StatusOK)
		render.JSON(w, r, totalledOrders)
	}
}

// PostOrder 		godoc
// @tags 			Order
// @Description 	Adds an order to the db
// @Accept 			json
// @Param 			b body model.Order true "Order"
// @Produce  		json
// @Success 		200
// @Failure     	400
// @Router 			/api/order [post]
func PostOrder(db *repository.DatabaseHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// declare empty order struct
		var order model.Order
		err := json.NewDecoder(r.Body).Decode(&order)
		// handle error and render Status 400
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, "invalid order")
			return
		}
		// no time sent: use the current time
		if order.CreatedAt.IsZero() {
			order.CreatedAt = time.Now()
		}
		// add to db
		db.AddOrder(&order)
		render.Status(r, http.StatusOK)
		render.JSON(w, r, "ok")
	}
}
