package api

import (
	"encoding/json"
	"fmt"
	"l0/models"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
)

type orders struct {
	orderModel models.OrderModel
}

func (o *orders) orders(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ordersUids := o.orderModel.ListOfUids()
	ans, err := json.Marshal(ordersUids)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: http api orders %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(ans)
}

func (o *orders) ordersUid(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	uid := ps.ByName("uid")
	order, err := o.orderModel.GetByUid(uid)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: http api orders uids %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ans, err := json.Marshal(order)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: http api orders uid %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(ans)
}

func MakeHandler(orderModel models.OrderModel) http.Handler {
	handlers := orders{
		orderModel: orderModel,
	}

	// Использовал первый попавшейся роутер с большим количеством звезд
	router := httprouter.New()
	router.ServeFiles("/static/*filepath", http.Dir("static"))
	router.GET("/api/v1/orders", handlers.orders)
	router.GET("/api/v1/orders/:uid", handlers.ordersUid)

	return router
}
