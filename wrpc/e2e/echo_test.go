package e2e

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/kong/go-wrpc/wrpc"
	"github.com/kong/go-wrpc/wrpc/e2e/testdata/echo"
)

type echoImpl struct {
	prefix string
}

func (e echoImpl) ProtocolError(_ context.Context,
	peer *wrpc.Peer, _ *echo.ProtocolErrorRequest,
) (*echo.ProtocolErrorResponse, error) {
	return nil, fmt.Errorf("to err is code")
}

func (e echoImpl) Sleep(_ context.Context,
	peer *wrpc.Peer,
	request *echo.SleepRequest,
) (*echo.SleepResponse, error) {
	time.Sleep(time.Duration(request.Duration) * time.Second)
	return &echo.SleepResponse{}, nil
}

func (e echoImpl) Echo(_ context.Context,
	peer *wrpc.Peer, req *echo.EchoRPCRequest,
) (*echo.EchoRPCResponse, error) {
	return &echo.EchoRPCResponse{
		S: e.prefix + req.S,
	}, nil
}

func TestEcho(t *testing.T) {
	t.Run("basic sanity test for happy path", func(t *testing.T) {
		// setup server
		echoServer := &echo.EchoServer{
			Echo: echoImpl{"foo-"},
		}
		serverPeer := &wrpc.Peer{}
		err := serverPeer.Register(echoServer)
		if err != nil {
			t.Fatalf("server failed to register Echo Service:%v", err)
		}
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			serverPeer.Upgrade(w, r)
		}))
		defer server.Close()

		// setup client
		clientPeer := &wrpc.Peer{}
		err = clientPeer.Register(&echo.EchoServer{})
		if err != nil {
			t.Fatalf("client failed to register Echo Service: %v", err)
		}

		c, _, err := wrpc.DefaultDialer.Dial(context.Background(),
			wsURLFromHTTPURL(server.URL), nil)
		if err != nil {
			t.Fatalf("client failed to connect to server: %v", err)
		}
		defer c.Close()
		clientPeer.AddConn(c)

		echoClient := &echo.EchoClient{Peer: clientPeer}
		resp, err := echoClient.Echo(context.Background(),
			&echo.EchoRPCRequest{S: "bar"})
		if err != nil {
			t.Fatalf("request failed: %v", err)
		}
		expectedResponse := "foo-bar"
		if resp.S != expectedResponse {
			t.Fatalf("expected response to be '%s' but go '%s'", expectedResponse,
				resp.S)
		}
	})
	t.Run("unregistered service in server returns an error",
		func(t *testing.T) {
			// setup server
			serverPeer := &wrpc.Peer{}
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				serverPeer.Upgrade(w, r)
			}))
			defer server.Close()

			// setup client
			clientPeer := &wrpc.Peer{}
			err := clientPeer.Register(&echo.EchoServer{})
			if err != nil {
				t.Fatalf("client failed to register Echo Service: %v", err)
			}

			c, _, err := wrpc.DefaultDialer.Dial(context.Background(),
				wsURLFromHTTPURL(server.URL), nil)
			if err != nil {
				t.Fatalf("client failed to connect to server: %v", err)
			}
			defer c.Close()
			clientPeer.AddConn(c)

			echoClient := &echo.EchoClient{Peer: clientPeer}
			_, err = echoClient.Echo(context.Background(),
				&echo.EchoRPCRequest{S: "bar"})
			if err == nil {
				t.Fatalf("expected an error")
			}
		})
	t.Run("unregistered service in client returns an error",
		func(t *testing.T) {
			// setup server
			echoServer := &echo.EchoServer{
				Echo: echoImpl{"foo-"},
			}
			serverPeer := &wrpc.Peer{}
			err := serverPeer.Register(echoServer)
			if err != nil {
				t.Fatalf("server failed to register Echo Service:%v", err)
			}
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				serverPeer.Upgrade(w, r)
			}))
			defer server.Close()

			// setup client
			clientPeer := &wrpc.Peer{}

			c, _, err := wrpc.DefaultDialer.Dial(context.Background(),
				wsURLFromHTTPURL(server.URL), nil)
			if err != nil {
				t.Fatalf("client failed to connect to server: %v", err)
			}
			defer c.Close()
			clientPeer.AddConn(c)

			echoClient := &echo.EchoClient{Peer: clientPeer}
			_, err = echoClient.Echo(context.Background(),
				&echo.EchoRPCRequest{S: "bar"})
			if err == nil {
				t.Fatalf("expected an error")
			}
		})
	t.Run("timeout is correctly processed in client",
		func(t *testing.T) {
			// setup server
			echoServer := &echo.EchoServer{
				Echo: echoImpl{"foo-"},
			}
			serverPeer := &wrpc.Peer{}
			err := serverPeer.Register(echoServer)
			if err != nil {
				t.Fatalf("server failed to register Echo Service:%v", err)
			}
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				serverPeer.Upgrade(w, r)
			}))
			defer server.Close()

			// setup client
			clientPeer := &wrpc.Peer{}
			err = clientPeer.Register(&echo.EchoServer{})
			if err != nil {
				t.Fatalf("client failed to register Echo Service: %v", err)
			}

			c, _, err := wrpc.DefaultDialer.Dial(context.Background(),
				wsURLFromHTTPURL(server.URL), nil)
			if err != nil {
				t.Fatalf("client failed to connect to server: %v", err)
			}
			defer c.Close()
			clientPeer.AddConn(c)

			echoClient := &echo.EchoClient{Peer: clientPeer}
			ctx, cancel := context.WithDeadline(context.Background(),
				time.Now().Add(100*time.Millisecond))
			defer cancel()
			_, err = echoClient.Sleep(ctx, &echo.SleepRequest{Duration: 100})
			if err == nil {
				t.Fatalf("error is nil but expected context deadline")
			}
		})
	t.Run("error from RPC is correctly propagated", func(t *testing.T) {
		// setup server
		echoServer := &echo.EchoServer{
			Echo: echoImpl{"foo-"},
		}
		serverPeer := &wrpc.Peer{}
		err := serverPeer.Register(echoServer)
		if err != nil {
			t.Fatalf("server failed to register Echo Service:%v", err)
		}
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			serverPeer.Upgrade(w, r)
		}))
		defer server.Close()

		// setup client
		clientPeer := &wrpc.Peer{}
		err = clientPeer.Register(&echo.EchoServer{})
		if err != nil {
			t.Fatalf("client failed to register Echo Service: %v", err)
		}

		c, _, err := wrpc.DefaultDialer.Dial(context.Background(),
			wsURLFromHTTPURL(server.URL), nil)
		if err != nil {
			t.Fatalf("client failed to connect to server: %v", err)
		}
		defer c.Close()
		clientPeer.AddConn(c)

		echoClient := &echo.EchoClient{Peer: clientPeer}
		resp, err := echoClient.ProtocolError(context.Background(), nil)
		if resp != nil {
			t.Errorf("didn't expect a response but got one")
		}
		if err == nil {
			t.Errorf("expeted an error but didn't get one")
		}
	})
}

func BenchmarkEcho(b *testing.B) {
	b.Run("basic sanity test", func(b *testing.B) {
		// setup server
		echoServer := &echo.EchoServer{
			Echo: echoImpl{"foo-"},
		}
		serverPeer := &wrpc.Peer{}
		err := serverPeer.Register(echoServer)
		if err != nil {
			b.Fatalf("server failed to register Echo Service:%v", err)
		}
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			serverPeer.Upgrade(w, r)
		}))
		defer server.Close()

		// setup client
		clientPeer := &wrpc.Peer{}
		err = clientPeer.Register(&echo.EchoServer{})
		if err != nil {
			b.Fatalf("client failed to register Echo Service: %v", err)
		}

		c, _, err := wrpc.DefaultDialer.Dial(context.Background(),
			wsURLFromHTTPURL(server.URL), nil)
		if err != nil {
			b.Fatalf("client failed to connect to server: %v", err)
		}
		defer c.Close()
		clientPeer.AddConn(c)
		echoClient := &echo.EchoClient{Peer: clientPeer}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			echoClient.Echo(context.Background(),
				&echo.EchoRPCRequest{S: "bar"})
		}
	})
}

func wsURLFromHTTPURL(urlStr string) string {
	url, _ := url.Parse(urlStr)
	url.Scheme = "ws"
	return url.String()
}
