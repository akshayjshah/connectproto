package connectproto

import (
	"bytes"
	"encoding/json"
	"errors"

	"github.com/bufbuild/connect-go"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// WithJSON customizes a connect-go Client or Handler's JSON handling.
func WithJSON(marshal protojson.MarshalOptions, unmarshal protojson.UnmarshalOptions) connect.Option {
	return connect.WithOptions(
		connect.WithCodec(&jsonCodec{name: "json", marshal: marshal, unmarshal: unmarshal}),
		connect.WithCodec(&jsonCodec{name: "json; charset=utf-8", marshal: marshal, unmarshal: unmarshal}),
	)
}

type jsonCodec struct {
	name      string
	marshal   protojson.MarshalOptions
	unmarshal protojson.UnmarshalOptions
}

func (j *jsonCodec) Name() string { return j.name }

func (j *jsonCodec) IsBinary() bool { return false }

func (j *jsonCodec) Unmarshal(binary []byte, msg any) error {
	pm, ok := msg.(proto.Message)
	if !ok {
		return errNotProto(msg)
	}
	if len(binary) == 0 {
		return errors.New("zero-length payload is not a valid JSON object")
	}
	return j.unmarshal.Unmarshal(binary, pm)
}

func (j *jsonCodec) Marshal(msg any) ([]byte, error) {
	pm, ok := msg.(proto.Message)
	if !ok {
		return nil, errNotProto(msg)
	}
	return j.marshal.Marshal(pm)
}

func (j *jsonCodec) MarshalStable(message any) ([]byte, error) {
	// protojson doesn't offer deterministic output. It does order fields by
	// number, but it deliberately introduce inconsistent whitespace (see
	// https://github.com/golang/protobuf/issues/1373). To make the output as
	// consistent as possible, we'll need to normalize.
	uncompacted, err := j.Marshal(message)
	if err != nil {
		return nil, err
	}
	compacted := bytes.NewBuffer(uncompacted[:0])
	if err = json.Compact(compacted, uncompacted); err != nil {
		return nil, err
	}
	return compacted.Bytes(), nil
}
