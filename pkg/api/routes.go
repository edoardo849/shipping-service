package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/edoardo849/bezos/pkg/order"
)

func handleOrdersCreate(os order.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		decoder := json.NewDecoder(r.Body)

		var req order.OrdersCreateReq
		err := decoder.Decode(&req)
		if err != nil {
			log.Println("Error while parsing the request", err.Error())
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		err = os.AddOrders(req)
		if err != nil {
			log.Println("Error while saving the orders", err.Error())
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		respondWithJSON(w, http.StatusOK, Status{true})
		return
	}
}
