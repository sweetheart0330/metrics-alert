package router

import (
	"fmt"
	"net/http"

	"github.com/sweetheart0330/metrics-alert/internal/handler"
	models "github.com/sweetheart0330/metrics-alert/internal/model"
)

func NewRouter(h handler.Handler) *http.ServeMux {
	r := http.NewServeMux()

	r.HandleFunc(
		fmt.Sprintf("/update/{%s}/{%s}/{%s}", models.TypeParam, models.NameParam, models.ValueParam),
		h.UpdateMetric,
	)

	return r
}
