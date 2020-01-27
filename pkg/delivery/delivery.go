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
	orderService order.ServiceInt
}

//New is the factory that initialises the Delivery service
func New(httpClient *http.Client, addr string, orderService order.ServiceInt) Service {
	return Service{
		httpClient:   httpClient,
		addr:         addr,
		orderService: orderService,
	}
}

//Start starts the delivery service
func (s *Service) Start(orderChan chan order.CreateReq) func() {

	stop := make(chan struct{}, 1)
	deliveryChan := make(chan order.Dispatch, 10)

	go func() {
		log.Println("Starting the delivery service 🚚")
		for {
			select {
			case o := <-orderChan:
				log.Println("Received order:", o.ID)
				res, err := s.deliver(o)
				if err != nil {
					log.Println("Error while calling the delivery service", err.Error())
					break
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
					//TODO delete the delivery here
					break
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

	jsonValue, err := json.Marshal(jsonData)
	if err != nil {
		log.Printf("Could not parse the delivery data: %s \n", err.Error())
		return response{}, err
	}
	log.Printf("Sending the delivery for Order #%d to %s \n", o.ID, s.addr)
	r, err := s.httpClient.Post(fmt.Sprintf("%s/api/orders/create", s.addr), "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Printf("The HTTP request failed with error %s\n", err)
		return response{}, err
	}

	log.Printf("Status Code %d\n", r.StatusCode)
	if r.StatusCode != http.StatusOK {
		return response{}, fmt.Errorf("Received Status Code %d", r.StatusCode)
	}

	decoder := json.NewDecoder(r.Body)

	var res response
	err = decoder.Decode(&res)
	if err != nil {
		log.Println("Could not parse the response", err)
		return response{}, err
	}
	log.Println("Order delivered", o.ID)

	return res, nil
}
