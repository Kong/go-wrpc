package wrpc

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	defaultHandshakeTimeout = 45 * time.Second
)

type Dialer struct {
	websocket.Dialer
}

func cloneHeader(h http.Header) http.Header {
	if h == nil {
		return http.Header{}
	}
	return h.Clone()
}

// Dial creates a new wRPC connection by calling DialContext with a
// background context.
func (d *Dialer) Dial(ctx context.Context, urlStr string,
	requestHeader http.Header) (*Conn, *http.Response, error) {
	requestHeader = cloneHeader(requestHeader)
	requestHeader.Add(secWebsocketProtoHeaderKey, wrpcProtoHeader)

	conn, resp, err := d.DialContext(ctx, urlStr,
		requestHeader)
	if err != nil {
		return nil, nil, err
	}
	if resp.StatusCode != http.StatusSwitchingProtocols {
		return nil, nil, fmt.Errorf("wRPC connection failed: unexpected HTTP response code: %v",
			resp.StatusCode)
	}
	return NewConn(conn), resp, nil
}

// DefaultDialer is a wRPC Dialer with all fields set to the default values.
var DefaultDialer = &Dialer{
	Dialer: websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: defaultHandshakeTimeout,
	},
}
