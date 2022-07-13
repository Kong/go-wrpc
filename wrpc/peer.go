package wrpc

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"

	. "github.com/kong/go-wrpc/wrpc/internal/wrpc"
	"google.golang.org/protobuf/proto"
)

// Peer represents a wRPC peer.
// A peer is capable of invoking and responding to RPCs.
type Peer struct {
	once sync.Once

	conn     *Conn
	registry *serviceRegistry

	ErrLogger func(error)
	// ClosedCallbackFunc if set, called after closing the connection.
	ClosedCallbackFunc func(*Peer)
}

func (p *Peer) init() {
	p.registry = &serviceRegistry{}
	if p.ErrLogger == nil {
		p.ErrLogger = func(error) {}
	}
}

func (p *Peer) Close() error {
	if p.conn == nil {
		return nil
	}
	return p.conn.Close()
}

func (p *Peer) RemoteAddr() net.Addr {
	if p.conn == nil {
		return nil
	}
	return p.conn.RemoteAddr()
}

// Register registers a Service for communication.
func (p *Peer) Register(s Service) error {
	p.once.Do(p.init)
	return p.registry.Add(s)
}

// AddConn adds a connection to peer and starts listenting for requests on the
// connection.
func (p *Peer) AddConn(conn *Conn) {
	p.once.Do(p.init)
	conn.handler = p
	conn.errLogger = p.ErrLogger

	p.conn = conn
	go func() {
		if err := p.conn.readThread(); err != nil {
			p.ErrLogger(fmt.Errorf("read thread: %w", err))
		}
		if err := p.conn.Close(); err != nil {
			p.ErrLogger(fmt.Errorf("closing wrpc connection: %w", err))
		}

		if p.ClosedCallbackFunc != nil {
			p.ClosedCallbackFunc(p)
		}
	}()
}

// Clone returns a new Peer with the same registry and logger, ready
// for either AddConn() or Upgrade()
func (p *Peer) Clone() *Peer {
	p.once.Do(p.init)
	p2 := &Peer{
		registry:  p.registry,
		ErrLogger: p.ErrLogger,
	}
	p2.once.Do(func() {})
	return p2
}

// Upgrade upgrades an HTTP connection to wRPC connection and starts tracking
// the connection.
func (p *Peer) Upgrade(w http.ResponseWriter, r *http.Request) error {
	p.once.Do(p.init)
	u := Upgrader{}
	c, err := u.Upgrade(w, r)
	if err != nil {
		return err
	}
	p.AddConn(c)
	return nil
}

func (p *Peer) fetchRPC(svcID, rpcID ID) (RPC, error) {
	p.once.Do(p.init)
	return p.registry.Get(svcID, rpcID)
}

func processResponse(in Response, out interface{}) error {
	if in.payload != nil {
		decode := decoderFunc(in.encoding, in.payload)
		return decode(out)
	}
	return nil
}

// VerifyRPC returns an error if the svcID/rpcID pair
// doesn't correspond to a registered RPC.
func (p *Peer) VerifyRPC(svcID, rpcID ID) error {
	p.once.Do(p.init)

	// verify the RPC is part of registered RPCs
	_, err := p.fetchRPC(svcID, rpcID)
	return err
}

// CreateRequest marshals and wraps a generic input as
// the payload in a new request
func CreateRequest(svcID, rpcID ID, input interface{}) (Request, error) {
	data, err := protoMarshal(input)
	if err != nil {
		return Request{}, err
	}

	return Request{
		svcID:    svcID,
		rpcID:    rpcID,
		encoding: Encoding_ENCODING_PROTO3,
		payload:  data,
	}, nil
}

// DoRequest sends a prepared request to the peer.
func (p *Peer) DoRequest(ctx context.Context, req Request, output proto.Message) error {
	p.once.Do(p.init)

	resp, err := p.conn.DoRPC(ctx, req)
	if err != nil {
		return fmt.Errorf("rpc: %w", err)
	}

	err = processResponse(resp, output)
	if err != nil {
		return fmt.Errorf("response: %w", err)
	}
	return nil
}

type Request struct {
	svcID, rpcID ID
	encoding     Encoding
	payload      []byte
}

type Response struct {
	encoding Encoding
	payload  []byte
}

// Handle is called by the underlying connection for every valid message.
func (p *Peer) handle(ctx context.Context, req Request) (Response, error) {
	p.once.Do(p.init)
	rpc, err := p.fetchRPC(req.svcID, req.rpcID)
	if err != nil {
		return Response{}, protocolError{
			error:        "invalid call: " + err.Error(),
			isContextual: true,
		}
	}

	resp, err := p.invokeHandler(ctx, rpc, req)
	if err != nil {
		return Response{}, protocolError{
			error:        "handler: " + err.Error(),
			isContextual: true,
		}
	}
	return resp, nil
}

func (p *Peer) invokeHandler(ctx context.Context, rpc RPC, req Request) (Response, error) {
	decodeFunc := decoderFunc(req.encoding, req.payload)

	result, err := rpc.Handler()(ctx, p, decodeFunc)
	if err != nil {
		return Response{}, err
	}

	return encodeResult(result)
}

func encodeResult(in interface{}) (Response, error) {
	encoding := Encoding_ENCODING_PROTO3
	encoder := encoderFunc(encoding)

	data, err := encoder(in)
	if err != nil {
		return Response{}, err
	}

	return Response{
		encoding: encoding,
		payload:  data,
	}, nil
}

func nilDecoder(_ interface{}) error {
	return fmt.Errorf("nil decoder")
}

func nilEncoder(_ interface{}) ([]byte, error) {
	return nil, fmt.Errorf("nil encoder")
}

func decoderFunc(encoding Encoding, data []byte) func(interface{}) error {
	switch encoding {
	case Encoding_ENCODING_PROTO3:
		return func(in interface{}) error {
			return protoUnmarshal(data, in)
		}
	case Encoding_ENCODING_UNSPECIFIED:
	default:
		return nilDecoder
	}
	return nil
}

func encoderFunc(encoding Encoding) func(interface{}) ([]byte, error) {
	switch encoding {
	case Encoding_ENCODING_PROTO3:
		return protoMarshal
	case Encoding_ENCODING_UNSPECIFIED:
	default:
		return nilEncoder
	}
	return nil
}

var errMissingProtoMessage = fmt.Errorf("input does not implement proto." +
	"Message interface")

func protoMarshal(m interface{}) ([]byte, error) {
	protoMessage, err := protoMessage(m)
	if err != nil {
		return nil, err
	}
	return proto.Marshal(protoMessage)
}

func protoUnmarshal(data []byte, m interface{}) error {
	protoMessage, err := protoMessage(m)
	if err != nil {
		return err
	}
	return proto.Unmarshal(data, protoMessage)
}

func protoMessage(in interface{}) (proto.Message, error) {
	protoMessage, ok := in.(proto.Message)
	if !ok {
		return nil, errMissingProtoMessage
	}
	return protoMessage, nil
}
