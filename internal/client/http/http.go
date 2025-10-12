package http

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	models "github.com/sweetheart0330/metrics-alert/internal/model"
)

const (
	updMetricPath   = "/update/"
	updMetricsBatch = "/updates/"
)

type Config struct {
	Host string `env:"ADDRESS"`
}
type Client struct {
	cfg Config
	cl  *http.Client
}

func NewClient(cfg Config) *Client {
	return &Client{cl: &http.Client{}, cfg: cfg}
}

func (c Client) SendGaugeMetric(m models.Metrics) error {
	if m.Value == nil {
		return fmt.Errorf("no value in metric")
	}

	//strVal := strconv.FormatFloat(*m.Value, 'f', -1, 64)

	resp, err := c.sendJSONRequest(m, updMetricPath)
	if err != nil {
		return fmt.Errorf("failed to send gauge metric, err: %w", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, err := getDecompressedBody(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read resp body body, err: %w", err)
		}

		return fmt.Errorf("bad response from server, status code: %d, error: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (c Client) SendCounterMetric(m models.Metrics) error {
	if m.Delta == nil {
		return fmt.Errorf("no delta in metric")
	}

	//strVal := strconv.FormatInt(*m.Delta, 10)

	resp, err := c.sendJSONRequest(m, updMetricPath)
	if err != nil {
		return fmt.Errorf("failed to send counter metric, err: %w", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, err := getDecompressedBody(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read resp body body, err: %w", err)
		}

		return fmt.Errorf("bad response from server, status code: %d, error: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (c Client) SendMetricsBatch(metrics []models.Metrics) error {
	resp, err := c.sendJSONRequest(metrics, updMetricsBatch)
	if err != nil {
		return fmt.Errorf("failed to send counter metric, err: %w", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, err := getDecompressedBody(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read resp body body, err: %w", err)
		}

		return fmt.Errorf("bad response from server, status code: %d, error: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (c Client) sendJSONRequest(data interface{}, method string) (*http.Response, error) {
	reqURL := formJSONURL(c.cfg.Host, method)
	jsonMetric, err := json.Marshal(&data)
	if err != nil {
		return nil, fmt.Errorf("failed to create json body, err: %w", err)
	}

	body, err := compressBody(jsonMetric)
	if err != nil {
		return nil, fmt.Errorf("failed to compress json body, err: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, reqURL, body)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}

	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("Content-Type", "application/json")
	//LogIncomingRequest(req)
	resp, err := c.cl.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not send request: %w", err)
	}
	//	LogOutgoingResponse(resp)
	return resp, nil
}

func (c Client) sendRequest(m models.Metrics, strVal string) (*http.Response, error) {
	reqURL := formURL(c.cfg.Host, m, strVal)
	req, err := http.NewRequest(http.MethodPost, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}

	req.Header.Set("Content-Type", "text/plain")

	resp, err := c.cl.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not send request: %w", err)
	}

	return resp, nil
}

func formJSONURL(url string, method string) string {
	builder := strings.Builder{}
	builder.WriteString(url)
	builder.WriteString(method)

	return builder.String()
}

func formURL(url string, m models.Metrics, val string) string {
	builder := strings.Builder{}
	builder.WriteString(url)
	builder.WriteString(updMetricPath)
	builder.WriteString("/")
	builder.WriteString(m.MType)
	builder.WriteString("/")
	builder.WriteString(m.ID)
	builder.WriteString("/")
	builder.WriteString(val)

	return builder.String()
}

func compressBody(body []byte) (io.Reader, error) {
	var buf bytes.Buffer
	gzWriter := gzip.NewWriter(&buf)
	_, err := gzWriter.Write(body)
	if err != nil {
		return nil, fmt.Errorf("failed to compress json body, err: %w", err)
	}
	defer gzWriter.Close()

	return &buf, nil
}

func getDecompressedBody(body io.ReadCloser) ([]byte, error) {
	gz, err := gzip.NewReader(body)
	if err != nil {
		return nil, fmt.Errorf("failed to decompress response body, err: %w", err)
	}

	body = struct {
		io.Reader
		io.Closer
	}{gz, body}
	defer body.Close()

	respBody, err := io.ReadAll(body)
	if err != nil {
		return nil, fmt.Errorf("failed to read resp body body, err: %w", err)
	}

	return respBody, nil
}
