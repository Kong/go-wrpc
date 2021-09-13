package wrpc

import (
	"context"
	"net/http"
	"net/http/httptest"
	url2 "net/url"
	"testing"

	"github.com/gorilla/websocket"
)

func TestUpgraderAndDialer(t *testing.T) {
	t.Run("missing sec-websocket-protocol header results in an error ", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,
			r *http.Request) {
			u := Upgrader{}
			u.Upgrade(w, r)
		}))
		defer server.Close()

		conn, resp, err := websocket.DefaultDialer.DialContext(context.Background(),
			wsURLFromHTTPURL(server.URL), nil)

		if err == nil {
			t.Errorf("expected handshake error but got: %v", err)
		}
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("expected HTTP 400 due to missing header but got %v",
				resp.StatusCode)
		}
		if resp.Header.Get(secWebsocketProtoHeaderKey) != wrpcProtoHeader {
			t.Errorf("expected %v header field to contain value '%v' but"+
				" got '%v'", secWebsocketProtoHeaderKey, wrpcProtoHeader,
				resp.Header.Get(secWebsocketProtoHeaderKey))
		}
		if conn != nil {
			t.Errorf("expected conn to be nil")
		}
	})
	t.Run("incorrect sec-websocket-protocol header results in an error ",
		func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,
				r *http.Request) {
				u := Upgrader{}
				u.Upgrade(w, r)
			}))
			defer server.Close()

			conn, resp, err := websocket.DefaultDialer.DialContext(context.Background(),
				wsURLFromHTTPURL(server.URL), http.Header{
					secWebsocketProtoHeaderKey: []string{"foo.konghq.com"},
				})

			if err == nil {
				t.Errorf("expected handshake error but got: %v", err)
			}
			if resp.StatusCode != http.StatusBadRequest {
				t.Errorf("expected HTTP 400 due to missing header but got %v",
					resp.StatusCode)
			}
			if resp.Header.Get(secWebsocketProtoHeaderKey) != wrpcProtoHeader {
				t.Errorf("expected %v header field to contain value '%v' but"+
					" got '%v'", secWebsocketProtoHeaderKey, wrpcProtoHeader,
					resp.Header.Get(secWebsocketProtoHeaderKey))
			}
			if conn != nil {
				t.Errorf("expected conn to be nil")
			}
		})
	t.Run("correct sec-websocket-protocol header results in a successful"+
		" upgrade", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,
			r *http.Request) {
			u := Upgrader{}
			u.Upgrade(w, r)
		}))
		defer server.Close()

		conn, resp, err := DefaultDialer.Dial(context.Background(),
			wsURLFromHTTPURL(server.URL), nil)
		if err != nil {
			t.Errorf("expected error to be nil but got %v", err)
		}
		if resp.StatusCode != http.StatusSwitchingProtocols {
			t.Errorf("expected HTTP 101 as a result of successful handshake")
		}
		if resp.Header.Get(secWebsocketProtoHeaderKey) != wrpcProtoHeader {
			t.Errorf("expected %v header field to contain value '%v' but"+
				" got '%v'", secWebsocketProtoHeaderKey, wrpcProtoHeader,
				resp.Header.Get(secWebsocketProtoHeaderKey))
		}
		if conn == nil {
			t.Errorf("expected conn to not be nil")
		}
	})
}

func wsURLFromHTTPURL(urlStr string) string {
	url, _ := url2.Parse(urlStr)
	url.Scheme = "ws"
	return url.String()
}
