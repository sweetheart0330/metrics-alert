package router

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/sweetheart0330/metrics-alert/internal/handler"
	models "github.com/sweetheart0330/metrics-alert/internal/model"
)

func NewRouter(h handler.Handler) *chi.Mux {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(h.MiddlewareLogger())
	mux.Use(h.CompressHandle)
	mux.Use(h.DecompressHandle)

	mux.Post("/update", h.UpdateJSONMetric)
	mux.Post("/update/", h.UpdateJSONMetric)

	mux.Post(
		fmt.Sprintf("/update/{%s}/{%s}/{%s}", models.TypeParam, models.NameParam, models.ValueParam),
		h.UpdateMetric,
	)
	mux.Post("/value", h.GetJSONMetric)
	mux.Post("/value/", h.GetJSONMetric)

	mux.Get(
		fmt.Sprintf("/value/{%s}/{%s}", models.TypeParam, models.NameParam),
		h.GetMetric,
	)

	mux.Get("/", h.GetAllMetrics)

	return mux
}
