package http

import (
	"context"
	"encoding/json"
	"kafkaExperiments/internal/domain/order/entities"
	"net/http"
)

type Server struct {
	mux *http.ServeMux
	ctl ordersController
}

type ordersController interface {
	Create(ctx context.Context, newOrder *entities.NewOrder) (*entities.Order, error)
}

func NewServer(ctl ordersController) *Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/order/create", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		order, err := ctl.Create(request.Context(), &entities.NewOrder{
			CustomerId: "aaa",
		})

		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			request.Body.Close()
			return
		}

		enc := json.NewEncoder(writer)

		err = enc.Encode(order)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		return
	})

	return &Server{
		mux: mux,
		ctl: ctl,
	}
}

func (s *Server) Run(ctx context.Context) error {
	return http.ListenAndServe(":8882", s.mux)
}
