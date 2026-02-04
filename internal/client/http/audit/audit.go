package audit

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	models "github.com/sweetheart0330/metrics-alert/internal/model"
)

type Client struct {
	cl  *http.Client
	url string
}

func NewClient(url string) *Client {
	return &Client{
		cl:  &http.Client{},
		url: url,
	}
}

func (c *Client) Consume(ctx context.Context, ev models.AuditEvent) error {
	b, err := json.Marshal(ev)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := c.cl.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("audit http status: %s", resp.Status)
	}

	return nil
}

func (c *Client) Close() error {
	return nil
}
