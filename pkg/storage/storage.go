package storage

import (
	"database/sql"
	"log"

	"github.com/edoardo849/bezos/pkg/delivery"
	"github.com/edoardo849/bezos/pkg/order"
)

// New is the factory that creates a storage
func New(db *sql.DB) *Storage {
	return &Storage{
		db: db,
	}
}

// Storage keeps the logic to perform CRUD operations
// against the Database
// TODO: accept http interface to pass to the Shipping Service
type Storage struct {
	db *sql.DB
}

//AddOrders saves a user to the storage
//TODO: pass a channel of order.OrdersCreateReq so that we can create deliveries
func (s *Storage) AddOrders(o order.OrdersCreateReq) error {
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

	deliveryTx, err := tx.Prepare("INSERT INTO api_db.order_delivery(order_id, delivery_id) VALUES(?,?)")
	if err != nil {
		log.Println("Failed to prepare the delivery transaction", err.Error())
		return err
	}

	oLen := len(o)
	for i := 0; i < oLen; i++ {

		log.Printf("Processing order %d \n", i)
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

		deliveryID, err := delivery.Deliver(currOrder)
		if err != nil {
			log.Println("Failed to deliver the order", err.Error())
			tx.Rollback()
			return err
		}

		// Do not rollback if a delivery wasn't saved
		//TODO remove this out of here and transform this into concurrent requests with goroutines through a channel of Orders
		_, err = deliveryTx.Exec(currOrder.ID, deliveryID)
		if err != nil {
			log.Println("Failed to save a delivery", err.Error())
			return err
		}

	}

	return tx.Commit()
}
