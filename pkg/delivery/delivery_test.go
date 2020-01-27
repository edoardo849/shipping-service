package delivery_test

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/DATA-DOG/go-sqlmock"

	"github.com/edoardo849/bezos/pkg/delivery"
	"github.com/edoardo849/bezos/pkg/order"
)

func testingHTTPClient(handler http.Handler) (*http.Client, func()) {
	s := httptest.NewServer(handler)

	cli := &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, network, _ string) (net.Conn, error) {
				return net.Dial(network, s.Listener.Addr().String())
			},
		},
	}

	return cli, s.Close
}

const (
	okResponse = `{
		"success": true,
		 "id": 999
	}`
)

type OrderMock struct{
	db *sql.DB
}



func (om OrderMock) AddOrders(o order.OrdersCreateReq) error {
	return nil
}
func (om OrderMock) Despatch(d order.Dispatch) error {
	return nil
}

func TestService_success(t *testing.T) {

	mockOrder := order.Order{
		Email: "test@test.com",
	}

	mockOrderReq := order.CreateReq{
		mockOrder,
		order.ShippingAddress{},
		[]order.ToShippingLine{},
	}

	exit := make(chan struct{})

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var req order.CreateReq
		decoder.Decode(&req)

		if req.Email != mockOrder.Email {
			t.Errorf("Email was incorrect, got: %s, want: %s.", req.Email, mockOrder.Email)
		}

		w.Write([]byte(okResponse))

		exit <- struct{}{}
	})

	client, stopClient := testingHTTPClient(h)
	defer stopClient()

	delService := delivery.New(client, "", OrderMock{

	})
	orderChan := make(chan order.CreateReq, 1)
	stopDeliveryService := delService.Start(orderChan)
	defer stopDeliveryService()

	orderChan <- mockOrderReq

	<-exit

}
