package main

import (
	"fmt"
	"goAdvancedTpl/internal/server/handlers"
	"goAdvancedTpl/internal/server/storage"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	metStorage := storage.NewMetricStorage()
	h := handlers.NewAPIHandler(metStorage)
	r := chi.NewRouter()

	r.Route("/update", func(r chi.Router) {
		r.Post("/", h.WriteWholeMetric)
		r.Post("/{metricType}/{metricName}/{metricValue}", h.WriteMetric)
	})
	r.Route("/value", func(r chi.Router) {
		r.Post("/", h.GetWholeMetric)
		r.Get("/{metricType}/{metricName}", h.GetMetric)
	})
	r.Get("/", h.AllMetrics)
	er := http.ListenAndServe(":8080", r)
	if er != nil {
		fmt.Println(er.Error())
	}
}
