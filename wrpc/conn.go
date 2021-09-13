package wrpc

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	. "github.com/kong/go-wrpc/wrpc/internal/wrpc"
	"google.golang.org/protobuf/proto"
)

const (
	defaultTimeout = 60 * time.Second
)

// Conn is a wRPC connection.
// A wRPC connection consists of an HTTP Websocket connection and some
// additional state around it.
type Conn struct {
	wbConn    *websocket.Conn
	readLock  sync.Mutex
	writeLock sync.Mutex
	seq       uint32
	seqLock   sync.Mutex
	inflight  sync.Map

	handler wRPCServerHandler

	errLogger func(error)
}

type wRPCServerHandler interface {
	handle(ctx context.Context, request Request) (Response, error)
}

// NewConn creates a wRCP connection using a websocket connection. It is assumed
// that the websocket connection has already been upgraded to wRPC during the
// handshake.
func NewConn(conn *websocket.Conn) *Conn {
	return &Conn{
		wbConn: conn,
	}
}

func (c *Conn) nextSeq() uint32 {
	c.seqLock.Lock()
	c.seq++
	r := c.seq
	c.seqLock.Unlock()
	return r
}

// Close closes the underlying connection.
// An attempt is made to close the connection gracefully and if that fails,
// then it is closed forcefully.
func (c *Conn) Close() error {
	err := c.writeToSocket(websocket.CloseMessage, nil)
	if err != nil {
		// TODO(hbagdi): error is being lost here, handle it
		return c.wbConn.Close()
	}
	return nil
}

// DoRPC invokes a request and returns the response.
// The returned error signals a protocol error and not an application-layer
// error.
func (c *Conn) DoRPC(ctx context.Context,
	req Request) (Response, error) {

	m := createRPCMessage(rpcMessageOpts{
		svcID:    req.svcID,
		rpcID:    req.rpcID,
		encoding: req.encoding,
		payload:  req.payload,
	})

	// initiate the RPC
	resultChan, cleanup, err := c.sendRPCMessage(ctx, m)
	if err != nil {
		return Response{}, fmt.Errorf("initiating rpc: %v", err)
	}
	defer cleanup()

	// wait for Response
	var res *WebsocketPayload
	select {
	case <-ctx.Done():
		return Response{}, ctx.Err()
	case res = <-resultChan:
	}

	switch res.Payload.Mtype {
	case MessageType_MESSAGE_TYPE_ERROR:
		return Response{}, fmt.Errorf("%v", res.Payload.Error.String())
	case MessageType_MESSAGE_TYPE_RPC:
		if len(res.Payload.Payloads) > 0 {
			return Response{
				encoding: res.Payload.PayloadEncoding,
				payload:  res.Payload.Payloads[0],
			}, nil
		}
		return Response{}, nil
	}
	return Response{}, nil
}

func (c *Conn) sendRPCMessage(ctx context.Context,
	m *WebsocketPayload) (chan *WebsocketPayload, func(), error) {
	var err error
	select {
	case <-ctx.Done():
		return nil, nil, fmt.Errorf("timeout")
	default:
	}

	seq := c.nextSeq()
	m.Payload.Seq = seq
	m.Payload.Deadline = deadlineFromCtx(ctx)

	resultChan := make(chan *WebsocketPayload, 1)
	c.inflight.Store(seq, resultChan)
	cleanup := func() { c.inflight.Delete(seq) }

	err = validateMessage(m)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid wRPC message: %v", err)
	}

	err = c.writeMessage(ctx, m)
	if err != nil {
		return nil, nil, fmt.Errorf("write message to socket: %v", err)
	}
	return resultChan, cleanup, nil
}

func (c *Conn) readFromSocket() (int, []byte, error) {
	c.readLock.Lock()
	defer c.readLock.Unlock()
	return c.wbConn.ReadMessage()
}

func (c *Conn) readThread() error {
	for {
		var err error
		mtype, wsMessage, err := c.readFromSocket()
		if err != nil {
			// websocket connection was closed
			_, ok := err.(*websocket.CloseError)
			if ok {
				return nil
			}
			if errors.Is(err, net.ErrClosed) {
				return nil
			}
			return err
		}

		if mtype != websocket.BinaryMessage {
			data := websocket.FormatCloseMessage(websocket.CloseUnsupportedData,
				"non-binary message not permitted")
			return c.writeToSocket(websocket.CloseMessage, data)
		}

		message := &WebsocketPayload{}
		err = proto.Unmarshal(wsMessage, message)
		if err != nil {
			data := websocket.FormatCloseMessage(websocket.
				CloseInvalidFramePayloadData,
				"invalid protobuf payload")
			return c.writeToSocket(websocket.CloseMessage, data)
		}

		err = validateMessage(message)
		if err != nil {
			data := websocket.FormatCloseMessage(websocket.
				CloseInvalidFramePayloadData,
				fmt.Sprintf("invalid wRPC message: %v", err))
			return c.writeToSocket(websocket.CloseMessage, data)
		}

		go func() {
			err = c.routeMessage(message)
			if err != nil {
				var pErr protocolError
				if errors.As(err, &pErr) {
					err := c.sendError(pErr, message)
					if err != nil {
						c.errLogger(fmt.Errorf(
							"wrpc failed to send an error to peer: %v", err))
					}
				} else {
					c.errLogger(fmt.Errorf("wrpc message routing failed: %v", err))
				}
			}
		}()
	}
}

func (c *Conn) sendError(err protocolError, reqM *WebsocketPayload) error {
	opt := errorMessageOpts{error: err}
	if err.isContextual {
		opt.ack = reqM.Payload.Seq
		opt.svcID = reqM.Payload.SvcId
		opt.rpcID = reqM.Payload.RpcId
	}
	m := createErrorMessage(opt)
	ctx, cancel := context.WithDeadline(context.Background(),
		time.Now().Add(defaultTimeout))
	defer cancel()
	return c.writeMessage(ctx, m)
}

type protocolError struct {
	error        string
	isContextual bool
}

func (e protocolError) Error() string {
	return e.error
}

func (c *Conn) routeMessage(m *WebsocketPayload) error {
	ack := m.Payload.Ack
	if ack > 0 {
		// m is a Response to an RPC
		resultChan, ok := c.inflight.Load(ack)
		if !ok {
			return protocolError{
				error:        "invalid response message",
				isContextual: true,
			}
		}
		resCh := resultChan.(chan *WebsocketPayload)
		resCh <- m
	}
	if ack == 0 {
		// m is a Request
		if m.Payload.Mtype == MessageType_MESSAGE_TYPE_RPC {
			// m is an initiation of an RPC
			deadline := int64(m.Payload.Deadline)
			ctx, cancel := context.WithDeadline(context.Background(),
				time.Unix(deadline, 0))
			defer cancel()

			req := Request{
				svcID:    ID(m.Payload.SvcId),
				rpcID:    ID(m.Payload.RpcId),
				encoding: m.Payload.PayloadEncoding,
			}
			if len(m.Payload.Payloads) > 0 {
				req.payload = m.Payload.Payloads[0]
			}
			resp, err := c.handler.handle(ctx, req)
			if err != nil {
				return fmt.Errorf("handler: %w", err)
			}

			m := createRPCMessage(rpcMessageOpts{
				svcID:    ID(m.Payload.SvcId),
				rpcID:    ID(m.Payload.RpcId),
				ack:      m.Payload.Seq,
				encoding: resp.encoding,
				payload:  resp.payload,
			})

			err = c.writeMessage(ctx, m)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *Conn) writeMessage(_ context.Context, m *WebsocketPayload) error {
	c.prepareMessageBeforeWrite(m)

	data, err := proto.Marshal(m)
	if err != nil {
		return err
	}
	return c.writeToSocket(websocket.BinaryMessage, data)
}

func (c *Conn) writeToSocket(messageType int, data []byte) error {
	c.writeLock.Lock()
	defer c.writeLock.Unlock()
	// TODO handle websocket errors
	return c.wbConn.WriteMessage(messageType, data)
}

func (c *Conn) prepareMessageBeforeWrite(m *WebsocketPayload) {
	if m.Version == 0 {
		m.Version = 1
	}
	if m.Payload.Seq == 0 {
		m.Payload.Seq = c.nextSeq()
	}
}
