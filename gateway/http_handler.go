package main

import (
	"errors"

	"github.com/gorilla/mux"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"net/http"

	"github.com/amrshaban2005/common"
	pb "github.com/amrshaban2005/common/api"
	"github.com/amrshaban2005/oms-gateway/gateway"
)

type handler struct {
	gateway gateway.OrdersGateway
}

func NewHandler(gateway gateway.OrdersGateway) *handler {
	return &handler{gateway}
}

func (h *handler) reigsterRoutes(r *mux.Router) {

	r.HandleFunc("/api/customers/{customerID}/orders", h.HandleCreateOrder).Methods("POST")

}

func (h *handler) HandleCreateOrder(w http.ResponseWriter, r *http.Request) {

	customerID := r.PathValue("customerID")

	var items []*pb.ItemsWithQuantity

	if err := common.ReadJSON(r, &items); err != nil {
		common.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validateItems(items); err != nil {
		common.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	o, err := h.gateway.CreateOrder(r.Context(), &pb.CreateOrderRequest{
		CustomerID: customerID,
		Items:      items,
	})

	rStatus := status.Convert(err)

	if rStatus != nil {
		if rStatus.Code() != codes.InvalidArgument {
			common.WriteError(w, http.StatusBadRequest, rStatus.Message())
			return
		}

		common.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	common.WriteJSON(w, http.StatusOK, o)
}

func validateItems(items []*pb.ItemsWithQuantity) error {
	if len(items) == 0 {
		return common.ErrNoItem
	}

	for _, i := range items {
		if i.ID == "" {
			return errors.New("you have to proivde valid item ID")
		}

		if i.Quantity <= 0 {
			return errors.New("you have to proivde valid quantity")
		}
	}
	return nil
}
