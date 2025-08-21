package handler

import (
	"html/template"
	"path/filepath"

	model "github.com/sweetheart0330/metrics-alert/internal/model"
	"github.com/sweetheart0330/metrics-alert/internal/service/contracts"
)

var metricHTML = `
<!DOCTYPE html>
<html>
<head><title>Metrics</title></head>
<body>
<h1>Metrics</h1>
<ul class="styled-list">
{{range $name, $metric := .}}
<li class="{{ $metric.MType }}">
    {{$metric.MType}} {{$metric.ID}}:
    {{if eq $metric.MType "` + model.Gauge + `"}}{{$metric.Value}}{{end}}
    {{if eq $metric.MType "` + model.Counter + `"}}{{$metric.Delta}}{{end}}
</li>
{{end}}
</ul>
</body>
</html>

<style>
.styled-list {
  list-style-type: square;
  padding-left: 20px;
}

.styled-list li {
  margin-bottom: 8px;
  padding-bottom: 8px;
  border-bottom: 1px solid #ccc;
}

.styled-list li:last-child {
  border-bottom: none;
  margin-bottom: 0;
  padding-bottom: 0;
}

.styled-list li.gauge {
  color: #2980b9;
  font-weight: bold;
}

.styled-list li.counter {
  color: #27ae60;
  font-weight: bold;
}
</style>
`

type Handler struct {
	metric   contracts.MetricService
	template *template.Template
}

func NewHandler(metric contracts.MetricService) (Handler, error) {
	tmplPath := filepath.Join("internal", "handler", "template", "metrics.html")
	tmpl, err := template.
		New("metrics.html").
		Funcs(funcMap).
		ParseFiles(tmplPath)
	if err != nil {
		return Handler{}, err
	}

	return Handler{metric: metric, template: tmpl}, nil
}
