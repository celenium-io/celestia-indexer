package celestianodeapi

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/dipdup-io/celestia-indexer/pkg/node/types"
	"github.com/goccy/go-json"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/time/rate"
)

type Node struct {
	rateLimit      *rate.Limiter
	client         *http.Client
	host           string
	jsonRpcVersion string
	token          string
	id             *atomic.Int64
	log            zerolog.Logger
}

func New(baseUrl string) *Node {
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
		id:             new(atomic.Int64),
		log:            log.With().Str("module", "celestia_node_api").Logger(),
	}
}

func (node *Node) WithRateLimit(requestPerSecond int) *Node {
	if requestPerSecond > 0 {
		node.rateLimit = rate.NewLimiter(rate.Every(time.Second/time.Duration(requestPerSecond)), requestPerSecond)
	}
	return node
}

func (node *Node) WithStartId(id int64) *Node {
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
func (node *Node) post(ctx context.Context, method string, params []any, output any) error {
	query := types.Request{
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

	start := time.Now()

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer closeWithLogError(response.Body, node.log)

	node.log.Trace().
		Int64("ms", time.Since(start).Milliseconds()).
		Str("method", query.Method).
		Int64("request_id", query.Id).
		Msg("request")

	if response.StatusCode != http.StatusOK {
		return errors.Errorf("invalid status: %d", response.StatusCode)
	}

	err = json.NewDecoder(response.Body).DecodeContext(ctx, output)
	return err
}

func closeWithLogError(stream io.ReadCloser, log zerolog.Logger) {
	if _, err := io.Copy(io.Discard, stream); err != nil {
		log.Err(err).Msg("api copy GET body response to discard")
	}
	if err := stream.Close(); err != nil {
		log.Err(err).Msg("api close GET body request")
	}
}
