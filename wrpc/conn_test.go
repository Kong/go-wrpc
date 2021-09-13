package wrpc

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
)

func TestConnReadThread(t *testing.T) {
	t.Run("connection is closed on a non-binary message",
		func(t *testing.T) {
			_, c, cleanup, err := setup()
			if err != nil {
				t.Errorf("expected nil but got err: %v", err)
			}
			defer cleanup()
			err = c.writeToSocket(websocket.TextMessage,
				[]byte("bad-message"))
			if err != nil {
				t.Error("failed to write text message")
			}

			_, _, err = c.readFromSocket()
			if err == nil {
				t.Error("expected an error but got none")
			}
			if !websocket.IsCloseError(err, websocket.CloseUnsupportedData) {
				t.Error("expected error code 1003 but didn't get one")
			}
		})
	t.Run("connection is closed on a malformed binary message",
		func(t *testing.T) {
			_, c, cleanup, err := setup()
			if err != nil {
				t.Errorf("expected nil but got err: %v", err)
			}
			defer cleanup()
			err = c.writeToSocket(websocket.BinaryMessage,
				[]byte("bad-binary-message"))
			if err != nil {
				t.Error("failed to write text message")
			}

			_, _, err = c.readFromSocket()
			if err == nil {
				t.Error("expected an error but got none")
			}
			if !websocket.IsCloseError(err,
				websocket.CloseInvalidFramePayloadData) {
				t.Error("expected error code 1007 but didn't get one")
			}
		})
	t.Run("closing connection aborts read thread",
		func(t *testing.T) {
			t.Skip()
			// setup server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,
				r *http.Request) {
				var u Upgrader
				u.Upgrade(w, r)
			}))
			defer server.Close()

			clientConn, _, err := DefaultDialer.Dial(context.Background(),
				wsURLFromHTTPURL(server.URL), nil)
			if err != nil {
				t.Errorf("client failed to connect to server: %v", err)
			}

			errChan := make(chan error)
			go func() {
				errChan <- clientConn.readThread()
			}()
			err = <-errChan
			if err != nil {
				t.Errorf("expected no error but got: %v", err)
			}
		})
}

func TestLocalAndRemoteAddr(t *testing.T) {
	_, c, cleanup, err := setup()
	if err != nil {
		t.Errorf("expected no error but got: %v", err)
	}
	defer cleanup()
	remote := c.RemoteAddr()
	if remote.Network() != "tcp" {
		t.Errorf("remote: expected 'tcp' but got: %v", remote.Network())
	}
	local := c.RemoteAddr()
	if local.Network() != "tcp" {
		t.Errorf("local: expected 'tcp' but got: %v", remote.Network())
	}
	if !strings.HasPrefix(remote.String(), "127.0.0.1") {
		t.Errorf("remote: expected 127.0.0.1 as suffix but got: %v",
			remote.String())
	}
	if !strings.HasPrefix(local.String(), "127.0.0.1") {
		t.Errorf("local: expected 127.0.0.1 as suffix but got: %v",
			local.String())
	}
}

//nolint:unparam
func setup() (*Peer, *Conn, func(), error) {
	// setup server
	serverPeer := &Peer{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,
		r *http.Request) {
		serverPeer.Upgrade(w, r)
	}))
	defer server.Close()

	c, _, err := DefaultDialer.Dial(context.Background(),
		wsURLFromHTTPURL(server.URL), nil)
	if err != nil {
		return nil, nil, nil, err
	}
	return serverPeer, c, func() { server.Close() }, nil
}
