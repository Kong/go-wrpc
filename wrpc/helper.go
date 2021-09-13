package wrpc

import (
	"context"
	"fmt"
	"time"

	. "github.com/kong/go-wrpc/wrpc/internal/wrpc"
)

func deadlineFromCtx(ctx context.Context) uint32 {
	if ctx == nil {
		deadline := time.Now().Add(defaultTimeout)
		return uint32(deadline.Unix())
	}
	deadline, ok := ctx.Deadline()
	if !ok {
		deadline = time.Now().Add(defaultTimeout)
	}
	return uint32(deadline.Unix())
}

type rpcMessageOpts struct {
	svcID, rpcID ID
	ack, seq     uint32
	encoding     Encoding
	payload      []byte
}

type errorMessageOpts struct {
	svcID, rpcID, ack, seq uint32
	error                  error
}

func createErrorMessage(opts errorMessageOpts) *WebsocketPayload {
	return &WebsocketPayload{
		Version: 1,
		Payload: &PayloadV1{
			Mtype: MessageType_MESSAGE_TYPE_ERROR,
			Error: &Error{
				// XXX
				Etype:       ErrorType_ERROR_TYPE_UNSPECIFIED,
				Description: opts.error.Error(),
			},
			SvcId: opts.svcID,
			RpcId: opts.rpcID,
			Seq:   opts.seq,
			Ack:   opts.ack,
		},
	}
}

func createRPCMessage(opts rpcMessageOpts) *WebsocketPayload {
	return &WebsocketPayload{
		Version: 1,
		Payload: &PayloadV1{
			Mtype:           MessageType_MESSAGE_TYPE_RPC,
			SvcId:           uint32(opts.svcID),
			RpcId:           uint32(opts.rpcID),
			Seq:             opts.seq,
			Ack:             opts.ack,
			PayloadEncoding: opts.encoding,
			Payloads:        [][]byte{opts.payload},
		},
	}
}

func validateMessage(m *WebsocketPayload) error {
	if m.Version != 1 {
		return fmt.Errorf("wrpc: invalid version: %v", m.Version)
	}
	if m.Payload == nil {
		return fmt.Errorf("no payload")
	}
	if m.Payload.Seq == 0 {
		return fmt.Errorf("invalid seq(0)")
	}
	switch m.Payload.Mtype {
	case MessageType_MESSAGE_TYPE_ERROR:
		return validateErrorMessage(m.Payload)
	case MessageType_MESSAGE_TYPE_RPC:
		return validateRPCMessage(m.Payload)
	case MessageType_MESSAGE_TYPE_UNSPECIFIED:
	default:
		return fmt.Errorf("invalid message type: %d", m.Payload.Mtype)
	}
	return nil
}

func validateErrorMessage(m *PayloadV1) error {
	if m.Error == nil {
		return fmt.Errorf("error message without any error")
	}
	return nil
}

func validateRPCMessage(m *PayloadV1) error {
	if m.SvcId == 0 {
		return fmt.Errorf("invalid svc_id(0)")
	}
	if m.RpcId == 0 {
		return fmt.Errorf("invalid rpc_id(0)")
	}
	if m.Ack == 0 && m.Deadline == 0 {
		return fmt.Errorf("invalid deadline(0) for request")
	}
	if m.Ack != 0 && m.Deadline != 0 {
		return fmt.Errorf("invalid deadline(%v) for response", m.Deadline)
	}
	numPayloads := len(m.Payloads)
	if numPayloads > 1 {
		return fmt.Errorf("unexpected number of payloads(%v) in a message",
			numPayloads)
	}
	if numPayloads > 0 {
		err := validateEncoding(m.PayloadEncoding)
		if err != nil {
			return err
		}
	}
	return nil
}

func validateEncoding(e Encoding) error {
	switch e {
	case Encoding_ENCODING_PROTO3:
		return nil
	case Encoding_ENCODING_UNSPECIFIED:
	default:
		return fmt.Errorf("invalid encoding(%v)", e)
	}
	return fmt.Errorf("invalid encoding(%v)", e)
}
