package blob

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/goccy/go-json"
	"github.com/pkg/errors"
	"golang.org/x/time/rate"
)

type jsonRpcRequest struct {
	Method  string `json:"method"`
	Params  []any  `json:"params"`
	Id      uint64 `json:"id"`
	JsonRpc string `json:"jsonrpc"`
}

type jsonRpcResponse[T any] struct {
	Id      uint64 `json:"id"`
	JsonRpc string `json:"jsonrpc"`
	Error   *Error `json:"error,omitempty"`
	Result  T      `json:"result"`
}

// Error -
type Error struct {
	Code    int64           `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

// Error -
func (e Error) Error() string {
	return fmt.Sprintf("code=%d message=%s data=%s", e.Code, e.Message, string(e.Data))
}

// errors
var (
	ErrRequest = errors.New("request error")
)

type Node struct {
	rateLimit      *rate.Limiter
	client         *http.Client
	host           string
	jsonRpcVersion string
	token          string
	id             *atomic.Uint64
}

func NewNode(baseUrl string) *Node {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 10
	t.MaxConnsPerHost = 10
	t.MaxIdleConnsPerHost = 10

	return &Node{
		host: baseUrl,
		client: &http.Client{
			Transport: t,
		},
		jsonRpcVersion: "2.0",
		id:             new(atomic.Uint64),
	}
}

func (node *Node) WithRateLimit(requestPerSecond int) *Node {
	if requestPerSecond > 0 {
		node.rateLimit = rate.NewLimiter(rate.Every(time.Second/time.Duration(requestPerSecond)), requestPerSecond)
	}
	return node
}

func (node *Node) WithStartId(id uint64) *Node {
	if id > 0 {
		node.id.Store(id)
	}
	return node
}

func (node *Node) WithJsonRpcVersion(version string) *Node {
	if version != "" {
		node.jsonRpcVersion = version
	}
	return node
}

func (node *Node) WithAuthToken(token string) *Node {
	if token != "" {
		node.token = token
	}
	return node
}

func (node *Node) Blobs(ctx context.Context, height uint64, hash ...string) ([]Blob, error) {
	if len(hash) == 0 {
		return nil, nil
	}

	var response jsonRpcResponse[[]Blob]
	if err := node.post(ctx, "blob.GetAll", []any{height, hash}, &response); err != nil {
		return nil, err
	}

	if response.Error != nil {
		return nil, errors.Wrapf(ErrRequest, "request %d error: %s", response.Id, response.Error.Error())
	}
	return response.Result, nil
}

func (node *Node) post(ctx context.Context, method string, params []any, output any) error {
	query := jsonRpcRequest{
		JsonRpc: node.jsonRpcVersion,
		Id:      node.id.Add(1),
		Method:  method,
		Params:  params,
	}

	body := new(bytes.Buffer)
	if err := json.NewEncoder(body).Encode(query); err != nil {
		return err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, node.host, body)
	if err != nil {
		return err
	}
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", node.token))

	if node.rateLimit != nil {
		if err := node.rateLimit.Wait(ctx); err != nil {
			return err
		}
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	buffer := new(bytes.Buffer)
	if _, err := io.Copy(buffer, response.Body); err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return errors.Errorf("invalid status: %d %s", response.StatusCode, buffer.String())
	}

	if err := json.NewDecoder(buffer).DecodeContext(ctx, output); err != nil {
		return err
	}
	return nil
}
