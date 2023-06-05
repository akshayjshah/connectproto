package connectproto

import (
	"github.com/bufbuild/connect-go"
	"google.golang.org/protobuf/proto"
)

// WithBinary customizes a connect-go Client or Handler's binary protobuf handling.
func WithBinary(marshal proto.MarshalOptions, unmarshal proto.UnmarshalOptions) connect.Option {
	return connect.WithCodec(newBinaryCodec(marshal, unmarshal))
}

type binaryCodec struct {
	name      string
	marshal   proto.MarshalOptions
	stable    proto.MarshalOptions
	unmarshal proto.UnmarshalOptions
}

func newBinaryCodec(marshal proto.MarshalOptions, unmarshal proto.UnmarshalOptions) *binaryCodec {
	stable := marshal
	stable.Deterministic = true
	return &binaryCodec{
		name:      "proto",
		marshal:   marshal,
		stable:    stable,
		unmarshal: unmarshal,
	}
}

func (b *binaryCodec) Name() string { return b.name }

func (b *binaryCodec) IsBinary() bool { return true }

func (b *binaryCodec) Unmarshal(binary []byte, msg any) error {
	pm, ok := msg.(proto.Message)
	if !ok {
		return errNotProto(msg)
	}
	return b.unmarshal.Unmarshal(binary, pm)
}

func (b *binaryCodec) Marshal(msg any) ([]byte, error) {
	return b.marshalBinary(msg, false /* stable */)
}

func (b *binaryCodec) MarshalStable(msg any) ([]byte, error) {
	return b.marshalBinary(msg, true /* stable */)
}

func (b *binaryCodec) marshalBinary(msg any, stable bool) ([]byte, error) {
	pm, ok := msg.(proto.Message)
	if !ok {
		return nil, errNotProto(msg)
	}
	if stable {
		return b.stable.Marshal(pm)
	}
	return b.marshal.Marshal(pm)
}
