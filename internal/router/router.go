package router

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/sweetheart0330/metrics-alert/internal/handler"
	models "github.com/sweetheart0330/metrics-alert/internal/model"
)

func NewRouter(h handler.Handler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(h.MiddlewareLogger())
	r.Use(h.CompressHandle)
	r.Use(h.DecompressHandle)

	r.Route("/update", func(r chi.Router) {
		r.Post("/", h.UpdateJSONMetric)
		r.Post(
			fmt.Sprintf("/{%s}/{%s}/{%s}", models.TypeParam, models.NameParam, models.ValueParam),
			h.UpdateMetric,
		)
	})

	r.Route("/value", func(r chi.Router) {
		r.Get("/", h.GetJSONMetric)
		r.Get(
			fmt.Sprintf("/{%s}/{%s}", models.TypeParam, models.NameParam),
			h.GetMetric,
		)
	})

	r.Get("/", h.GetAllMetrics)

	return r
}
