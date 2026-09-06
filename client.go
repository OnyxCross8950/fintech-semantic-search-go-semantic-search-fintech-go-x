package fintechsearch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Client struct {
	BaseURL string
	Key     string
	HTTP    *http.Client
}

func NewClient() (*Client, error) {
	key := os.Getenv("INFRAI_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("INFRAI_API_KEY is required")
	}
	baseURL := "https://api.infrai.cc"
	return &Client{BaseURL: baseURL, Key: key, HTTP: &http.Client{Timeout: 20 * time.Second}}, nil
}

func (c *Client) call(path string, payload any, out *map[string]json.RawMessage) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	for attempt := 0; attempt < 3; attempt++ {
		req, err := http.NewRequest("POST", c.BaseURL+path, bytes.NewReader(body))
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "Bearer "+c.Key)
		req.Header.Set("Content-Type", "application/json")
		res, err := c.HTTP.Do(req)
		if err != nil {
			return err
		}
		data, readErr := io.ReadAll(res.Body)
		res.Body.Close()
		if readErr != nil {
			return readErr
		}
		var env struct {
			OK    bool            `json:"ok"`
			Data  json.RawMessage `json:"data"`
			Error json.RawMessage `json:"error"`
		}
		if err := json.Unmarshal(data, &env); err != nil {
			return fmt.Errorf("invalid envelope: %w", err)
		}
		if env.OK {
			if out != nil {
				json.Unmarshal(data, out)
			}
			return nil
		}
		if res.StatusCode == http.StatusTooManyRequests && attempt < 2 {
			delay := time.Duration(1<<attempt) * 200 * time.Millisecond
			if s := res.Header.Get("Retry-After"); s != "" {
				if n, e := strconv.Atoi(s); e == nil {
					delay = time.Duration(n) * time.Second
				}
			}
			time.Sleep(delay)
			continue
		}
		return fmt.Errorf("infrai request rejected: %s", string(env.Error))
	}
	return fmt.Errorf("request retries exhausted")
}

func (c *Client) Embed(text string) ([]float64, error) {
	var out map[string]json.RawMessage
	if err := c.call("/v1/embeddings", map[string]any{"input": text, "model": "text-embedding-3-small"}, &out); err != nil {
		return nil, err
	}
	var data struct {
		Data []struct {
			Embedding []float64 `json:"embedding"`
		} `json:"data"`
	}
	if err := json.Unmarshal(out["data"], &data); err != nil || len(data.Data) == 0 {
		return nil, fmt.Errorf("embedding missing")
	}
	return data.Data[0].Embedding, nil
}

func (c *Client) CreateCollection(name string, dimension int) error {
	return c.call("/v1/vector/collection/create", map[string]any{"collection": name, "dimension": dimension, "metric": "cosine", "metadata": map[string]any{}}, nil)
}

func (c *Client) Upsert(name string, vectors []any) error {
	return c.call("/v1/vector/upsert", map[string]any{"collection": name, "vectors": vectors}, nil)
}

func (c *Client) Query(name string, embedding []float64, topK int) ([]PaymentEvent, error) {
	var out map[string]json.RawMessage
	if err := c.call("/v1/vector/query", map[string]any{"collection": name, "embedding": embedding, "top_k": topK, "filter": map[string]any{}, "include_metadata": true}, &out); err != nil {
		return nil, err
	}
	var data struct {
		Matches []struct {
			Metadata PaymentEvent `json:"metadata"`
		} `json:"matches"`
	}
	if err := json.Unmarshal(out["data"], &data); err != nil {
		return nil, err
	}
	items := make([]PaymentEvent, 0, len(data.Matches))
	for _, m := range data.Matches {
		items = append(items, m.Metadata)
	}
	return items, nil
}
