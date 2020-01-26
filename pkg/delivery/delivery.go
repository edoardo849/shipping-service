package delivery

import (
	"bytes"
	"encoding/json"
	"fmt"
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

//Deliver delivers the order
// TODO: accept http interface to that it becomes testable
// TODO: accept a global variable for the shipping service address
func Deliver(o order.OrderCreateReq) (int64, error) {

	priceFloat, err := strconv.ParseFloat(o.ShippingLines[0].Price, 64)
	if err != nil {
		return 0, err
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
		return 0, err
	}

	decoder := json.NewDecoder(r.Body)

	var response response
	err = decoder.Decode(&response)
	if err != nil {
		fmt.Println("Could not parse the response", err)

		return 0, err
	}

	return response.ID, nil
}
