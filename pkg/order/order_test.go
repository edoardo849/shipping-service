package order_test

import (
	"database/sql"
	"log"
	"testing"

	"github.com/edoardo849/bezos/pkg/order"
	_ "github.com/go-sql-driver/mysql"
)

// a successful case
func TestShouldCreateOrders(t *testing.T) {
	db, err := sql.Open("mysql", "docker:docker@tcp(localhost:3306)/api_db")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	doc := make(chan order.CreateReq, 1)

	orderService := order.New(db, doc)

	shippingAddr := order.ShippingAddress{
		FirstName: "edo",
		Address1:  "bla",
		PostCode:  "bla",
	}

	currOrder := order.Order{
		ID:    42353463,
		Email: "bla",
	}

	ocreq := order.OrdersCreateReq{
		order.CreateReq{
			currOrder,
			shippingAddr,
			[]order.ToShippingLine{},
		},
	}

	if err := orderService.AddOrders(ocreq); err != nil {
		t.Errorf("The service should not return an error: %s ", err)
	}

	query := db.QueryRow("SELECT id, email FROM api_db.order WHERE id=?", currOrder.ID)

	var orderID int64
	var email string
	query.Scan(&orderID, &email)

	if orderID != currOrder.ID {
		t.Errorf("Expected OrderID: %d got %d ", currOrder.ID, orderID)
	}

	if email != currOrder.Email {
		t.Errorf("Expected email: %s got %s", currOrder.Email, email)
	}

}
