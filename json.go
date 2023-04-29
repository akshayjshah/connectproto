package connectproto

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/bufbuild/connect-go"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/runtime/protoiface"
)

// WithJSON customizes a connect-go Client or Handler's JSON handling.
func WithJSON(marshal protojson.MarshalOptions, unmarshal protojson.UnmarshalOptions) connect.Option {
	return connect.WithOptions(
		connect.WithCodec(&codec{name: "json", marshal: marshal, unmarshal: unmarshal}),
		connect.WithCodec(&codec{name: "json; charset=utf-8", marshal: marshal, unmarshal: unmarshal}),
	)
}

type codec struct {
	name      string
	marshal   protojson.MarshalOptions
	unmarshal protojson.UnmarshalOptions
}

func (c *codec) Name() string { return c.name }

func (c *codec) IsBinary() bool { return false }

func (c *codec) Unmarshal(binary []byte, msg any) error {
	pm, ok := msg.(proto.Message)
	if !ok {
		return errNotProto(msg)
	}
	if len(binary) == 0 {
		return errors.New("zero-length payload is not a valid JSON object")
	}
	return c.unmarshal.Unmarshal(binary, pm)
}

func (c *codec) Marshal(msg any) ([]byte, error) {
	pm, ok := msg.(proto.Message)
	if !ok {
		return nil, errNotProto(msg)
	}
	return c.marshal.Marshal(pm)
}

func (c *codec) MarshalStable(message any) ([]byte, error) {
	// protojson doesn't offer deterministic output. It does order fields by
	// number, but it deliberately introduce inconsistent whitespace (see
	// https://github.com/golang/protobuf/issues/1373). To make the output as
	// consistent as possible, we'll need to normalize.
	uncompacted, err := c.Marshal(message)
	if err != nil {
		return nil, err
	}
	compacted := bytes.NewBuffer(uncompacted[:0])
	if err = json.Compact(compacted, uncompacted); err != nil {
		return nil, err
	}
	return compacted.Bytes(), nil
}

func errNotProto(msg any) error {
	if _, ok := msg.(protoiface.MessageV1); ok {
		return fmt.Errorf("%T uses github.com/golang/protobuf, but connect-go only supports google.golang.org/protobuf: see https://go.dev/blog/protobuf-apiv2", msg)
	}
	return fmt.Errorf("%T doesn't implement proto.Message", msg)
}
