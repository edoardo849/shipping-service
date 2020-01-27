package delivery

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/edoardo849/bezos/pkg/order"
)

// Delivery holds the delivery information
type Delivery struct {
	Client   string  `json:"client"`
	Price    float64 `json:"price"`
	Customer string  `json:"customer"`
	Parcel   Parcel  `json:"parcel"`
	Address  Address `json:"address"`
}

// Parcel is
type Parcel struct {
	Weight  float64 `json:"client"`
	Courier string  `json:"courier"`
}

//Address is
type Address struct {
	Postcode string `json:"postcode"`
}

type response struct {
	Success bool  `json:"success"`
	ID      int64 `json:"id"`
}

// Service is the delivery service
type Service struct {
	httpClient   *http.Client
	addr         string
	orderService order.Service
}

//New is the factory that initialises the Delivery service
func New(httpClient *http.Client, addr string, orderService order.Service) Service {
	return Service{
		httpClient:   httpClient,
		addr:         addr,
		orderService: orderService,
	}
}

//Start starts the delivery service
func (s *Service) Start(orderChan chan order.CreateReq) func() {

	stop := make(chan struct{}, 1)
	deliveryChan := make(chan order.Dispatch, 1)

	go func() {
		log.Println("Starting the delivery service 🚚")
		for {
			select {
			case o := <-orderChan:
				log.Println("Received order:", o.ID)
				res, err := s.deliver(o)
				if err != nil {
					log.Println("Error while calling the delivery service", err.Error())
					continue
				}
				status := "success"
				if !res.Success {
					status = "error"
				}
				deliveryChan <- order.Dispatch{
					OrderID:      o.ID,
					DispatcherID: res.ID,
					Status:       status,
				}
			case d := <-deliveryChan:
				log.Println("Recevied delivery", d)
				err := s.orderService.Despatch(d)
				if err != nil {
					log.Println("Could not save the despatch information", err.Error())
					// delete the delivery here
				}
				log.Println("Delivery information saved, package on the way 📦")

			case <-stop:
				log.Println("Received stop signal, stopping the delivery service")
				close(orderChan)
				close(deliveryChan)
				return
			}
		}
	}()

	return func() {
		stop <- struct{}{}
	}
}

//Deliver delivers the order
// TODO: accept http interface to that it becomes testable
// TODO: accept a global variable for the shipping service address
func (s *Service) deliver(o order.CreateReq) (response, error) {

	priceFloat, err := strconv.ParseFloat(o.ShippingLines[0].Price, 64)
	if err != nil {
		return response{}, err
	}

	wFloat := float64(o.TotalWeightGrams)

	jsonData := Delivery{
		Client:   "Luke Skywalker LTD",
		Price:    priceFloat,
		Customer: o.Email,
		Parcel: Parcel{
			Weight:  wFloat * 0.00220462262185,
			Courier: o.ShippingLines[0].Title,
		},
		Address: Address{
			Postcode: o.ShippingAddress.PostCode,
		},
	}

	jsonValue, _ := json.Marshal(jsonData)
	r, err := http.Post("http://delivery:8081/api/orders/create", "application/json", bytes.NewBuffer(jsonValue))

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		return response{}, err
	}

	decoder := json.NewDecoder(r.Body)

	var res response
	err = decoder.Decode(&res)
	if err != nil {
		fmt.Println("Could not parse the response", err)
		return response{}, err
	}

	return res, nil
}
