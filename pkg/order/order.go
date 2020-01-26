package order

//OrderToShippingLine is the shipping line in the Warehouse where the order's items
//are stored
type OrderToShippingLine struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
	Price string `json:"price"`
}

//OrderToShippingLineDB is
type OrderToShippingLineDB struct {
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

//OrderCreateReq
type OrderCreateReq struct {
	Order
	ShippingAddress ShippingAddress       `json:"shipping_address"`
	ShippingLines   []OrderToShippingLine `json:"shipping_lines"`
}

//OrdersCreateReq is the REST request to create a list of orders
type OrdersCreateReq []OrderCreateReq
