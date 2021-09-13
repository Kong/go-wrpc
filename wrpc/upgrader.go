package wrpc

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type Upgrader struct{}

const (
	secWebsocketProtoHeaderKey = "Sec-WebSocket-Protocol"
	wrpcProtoHeader            = "wrpc.konghq.com"

	upgradeHandshakeTimeout = 60 * time.Second
)

func hasWrpcProtoHeader(protos []string) bool {
	for _, proto := range protos {
		if proto == wrpcProtoHeader {
			return true
		}
	}
	return false
}

func (u *Upgrader) Upgrade(w http.ResponseWriter, r *http.Request) (*Conn, error) {
	clientProtos := websocket.Subprotocols(r)
	if !hasWrpcProtoHeader(clientProtos) {
		w.Header().Set(secWebsocketProtoHeaderKey, wrpcProtoHeader)
		w.WriteHeader(http.StatusBadRequest)
	}

	wsUpgrader := &websocket.Upgrader{
		HandshakeTimeout: upgradeHandshakeTimeout,
		Subprotocols:     []string{wrpcProtoHeader},
	}

	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, fmt.Errorf("wrpc upgrade failed: %w", err)
	}
	return &Conn{
		wbConn: conn,
	}, nil
}
