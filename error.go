package connectproto

import (
	"fmt"

	"google.golang.org/protobuf/runtime/protoiface"
)

func errNotProto(msg any) error {
	if _, ok := msg.(protoiface.MessageV1); ok {
		return fmt.Errorf("%T uses github.com/golang/protobuf, but connectproto only supports google.golang.org/protobuf: see https://go.dev/blog/protobuf-apiv2", msg)
	}
	return fmt.Errorf("%T doesn't implement proto.Message", msg)
}
