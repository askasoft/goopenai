package openai

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/askasoft/pango/iox"
	"github.com/askasoft/pango/log"
	"github.com/askasoft/pango/log/httplog"
	"github.com/askasoft/pango/ret"
)

// alias for Client
type OpenAI = Client

type Client struct {
	Domain string
	Apikey string

	Transport http.RoundTripper
	Timeout   time.Duration
	Logger    log.Logger

	MaxRetries  int
	RetryAfter  time.Duration
	ShouldRetry func(error) bool // default retry on not canceled error or (status = 429 || (status >= 500 && status <= 599))
}

func (c *Client) endpoint(format string, args ...any) string {
	return "https://" + c.Domain + "/v1" + fmt.Sprintf(format, args...)
}

func (c *Client) shouldRetry(err error) bool {
	sr := c.ShouldRetry
	if sr == nil {
		sr = shouldRetry
	}
	return sr(err)
}

func (c *Client) call(req *http.Request) (res *http.Response, err error) {
	client := &http.Client{
		Transport: c.Transport,
		Timeout:   c.Timeout,
	}

	res, err = httplog.TraceClientDo(c.Logger, client, req)
	if err != nil {
		if c.shouldRetry(err) {
			err = ret.NewRetryError(err, c.RetryAfter)
		}
	}
	return
}

func (c *Client) RetryForError(ctx context.Context, api func() error) (err error) {
	return ret.RetryForError(ctx, api, c.MaxRetries, c.Logger)
}

func (c *Client) authenticate(req *http.Request) {
	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", contentTypeJSON)
	}

	req.Header.Set("Authorization", "Bearer "+c.Apikey)
}

func (c *Client) doCall(req *http.Request, result any) error {
	c.authenticate(req)

	res, err := c.call(req)
	if err != nil {
		return err
	}
	defer iox.DrainAndClose(res.Body)

	decoder := json.NewDecoder(res.Body)
	if res.StatusCode == http.StatusOK {
		if result != nil {
			return decoder.Decode(result)
		}
		return nil
	}

	re := newResultError(res)
	_ = decoder.Decode(re)

	if c.shouldRetry(re) {
		re.RetryAfter = c.RetryAfter
	}
	return re
}

func (c *Client) DoPost(ctx context.Context, url string, source, result any) error {
	return c.RetryForError(ctx, func() error {
		return c.doPost(ctx, url, source, result)
	})
}

func (c *Client) doPost(ctx context.Context, url string, source, result any) error {
	buf, ct, err := buildJsonRequest(source)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, buf)
	if err != nil {
		return err
	}
	if ct != "" {
		req.Header.Set("Content-Type", ct)
	}

	return c.doCall(req, result)
}

// https://platform.openai.com/docs/api-reference/chat/create
func (c *Client) CreateChatCompletion(ctx context.Context, req *ChatCompletionRequest) (*ChatCompletionResponse, error) {
	url := c.endpoint("/chat/completions")

	res := &ChatCompletionResponse{}
	err := c.DoPost(ctx, url, req, res)
	return res, err
}

// https://platform.openai.com/docs/api-reference/embeddings/create
func (c *Client) CreateTextEmbeddings(ctx context.Context, req *TextEmbeddingsRequest) (*TextEmbeddingsResponse, error) {
	url := c.endpoint("/embeddings")

	res := &TextEmbeddingsResponse{}
	err := c.DoPost(ctx, url, req, res)
	return res, err
}
