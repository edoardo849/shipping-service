package order

import (
	"database/sql"
	"log"
)

//ToShippingLine is the shipping line in the Warehouse where the order's items
//are stored
type ToShippingLine struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
	Price string `json:"price"`
}

//ToShippingLineDB is
type ToShippingLineDB struct {
	ShippingLineID int64  `json:"shipping_line_id"`
	OrderID        int64  `json:"order_id"`
	Title          string `json:"title"`
	Price          string `json:"price"`
}

//ShippingAddress is the address of the shipment
type ShippingAddress struct {
	FirstName string `json:"first_name"`
	Address1  string `json:"address1"`
	City      string `json:"city"`
	PostCode  string `json:"postcode"`
}

//ShippingAddressDB is the address of the shipment
type ShippingAddressDB struct {
	OrderID int64 `json:"order_id"`
	ShippingAddress
}

//Order is an order
type Order struct {
	ID               int64  `json:"id"`
	Email            string `json:"email"`
	TotalPrice       string `json:"total_price"`
	TotalWeightGrams int    `json:"total_weight_grams"`
	OrderNumber      int    `json:"order_number"`
}

//Dispatch holds the information of the despatcher
type Dispatch struct {
	OrderID      int64  `json:"order_id"`
	DispatcherID int64  `json:"dispatcher_id"`
	Status       string `json:"status"`
}

//CreateReq is the request to create an order
type CreateReq struct {
	Order
	ShippingAddress ShippingAddress  `json:"shipping_address"`
	ShippingLines   []ToShippingLine `json:"shipping_lines"`
}

//OrdersCreateReq is the REST request to create a list of orders
type OrdersCreateReq []CreateReq

// New is the factory that creates an order service
func New(db *sql.DB, deliverOrderChan chan CreateReq) Service {
	return Service{
		db:               db,
		deliverOrderChan: deliverOrderChan,
	}
}

// Service keeps the logic to perform CRUD operations
// against the Database
// TODO: accept http interface to pass to the Shipping Service
type Service struct {
	db               *sql.DB
	deliverOrderChan chan CreateReq
}

//AddOrders saves a user to the storage
//TODO: pass a channel of order.OrdersCreateReq so that we can create deliveries
func (s Service) AddOrders(o OrdersCreateReq) error {
	tx, err := s.db.Begin()
	if err != nil {
		log.Println("Failed to start a transaction", err.Error())
		return err
	}

	orderTx, err := tx.Prepare("INSERT INTO api_db.order(id, email, total_price, total_weight_grams, order_number) VALUES(?,?,?,?,?)")
	if err != nil {
		log.Println("Failed to prepare the order transaction", err.Error())
		return err
	}

	addrTx, err := tx.Prepare("INSERT INTO api_db.order_shipping_address(order_id, first_name, address1, postcode) VALUES(?, ?, ?, ?)")
	if err != nil {
		log.Println("Failed to prepare the Address transaction", err.Error())
		return err
	}

	shippingLineTx, err := tx.Prepare("INSERT INTO api_db.order_to_shipping_line(order_id, shipping_line_id, title, price) VALUES(?, ?, ?, ?)")
	if err != nil {
		log.Println("Failed to prepare the Shipping Line transaction", err.Error())

		return err
	}

	oLen := len(o)
	for i := 0; i < oLen; i++ {

		log.Printf("Processing order %d 📤\n", i)
		currOrder := o[i]
		_, err = orderTx.Exec(
			currOrder.ID,
			currOrder.Email,
			currOrder.TotalPrice,
			currOrder.TotalWeightGrams,
			currOrder.OrderNumber,
		)
		if err != nil {
			tx.Rollback()
			log.Println("Failed to save an order", err.Error())

			return err
		}

		shippingAddr := currOrder.ShippingAddress

		// insert the shipping adddress
		_, err := addrTx.Exec(currOrder.ID, shippingAddr.FirstName, shippingAddr.Address1, shippingAddr.PostCode)
		if err != nil {
			log.Println("Failed to save the shipping address", err.Error())
			return err
		}

		slLen := len(currOrder.ShippingLines)
		for j := 0; j < slLen; j++ {

			sl := currOrder.ShippingLines[j]
			_, err := shippingLineTx.Exec(currOrder.ID, sl.ID, sl.Title, sl.Price)
			if err != nil {
				log.Println("Failed to save a shipping line", err.Error())
				tx.Rollback()
				return err
			}
		}

		s.deliverOrderChan <- currOrder

	}

	return tx.Commit()
}

// Despatch saves the despatcher's information
func (s Service) Despatch(d Dispatch) error {
	tx, err := s.db.Begin()
	if err != nil {
		log.Println("Failed to start a transaction", err.Error())
		return err
	}

	deliveryTx, err := tx.Prepare("INSERT INTO api_db.order_delivery(order_id, delivery_id) VALUES(?,?)")
	if err != nil {
		log.Println("Failed to prepare the delivery transaction", err.Error())
		return err
	}

	// Do not rollback if a delivery wasn't saved
	//TODO remove this out of here and transform this into concurrent requests with goroutines through a channel of Orders
	_, err = deliveryTx.Exec(d.OrderID, d.DispatcherID)
	if err != nil {
		log.Println("Failed to save a delivery", err.Error())
		return err
	}

	return nil
}
