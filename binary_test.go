package connectproto

import (
	"bytes"
	"testing"

	"go.akshayshah.org/attest"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestBinaryUnmarshal(t *testing.T) {
	codec := &binaryCodec{
		name:      "proto",
		unmarshal: proto.UnmarshalOptions{Merge: true},
	}
	t.Run("merge", func(t *testing.T) {
		enum := descriptorpb.EnumDescriptorProto{Name: ptr("foo")}
		err := codec.Unmarshal([]byte{}, &enum)
		attest.Ok(t, err)
		attest.Equal(t, *enum.Name, "foo")
	})
	t.Run("empty input", func(t *testing.T) {
		var empty emptypb.Empty
		err := codec.Unmarshal([]byte{}, &empty)
		attest.Ok(t, err)
	})
	t.Run("not protobuf", func(t *testing.T) {
		err := codec.Unmarshal([]byte{}, &struct{}{})
		attest.Error(t, err)
	})
}

func TestBinaryMarshal(t *testing.T) {
	codec := newBinaryCodec(proto.MarshalOptions{}, proto.UnmarshalOptions{})
	dict, err := structpb.NewStruct(map[string]any{
		"foo": "bar",
		"baz": "quux",
	})
	attest.Ok(t, err)
	t.Run("unstable", func(t *testing.T) {
		out, err := codec.Marshal(dict)
		attest.Ok(t, err)
		var isUnstable bool
		for i := 0; i < 100; i++ {
			again, err := codec.MarshalStable(dict)
			attest.Ok(t, err)
			if !bytes.Equal(out, again) {
				isUnstable = true
			}
		}
		attest.True(t, isUnstable, attest.Sprint("Marshal produced stable output"))
	})
	t.Run("stable", func(t *testing.T) {
		out, err := codec.MarshalStable(dict)
		attest.Ok(t, err)
		for i := 0; i < 100; i++ {
			again, err := codec.MarshalStable(dict)
			attest.Ok(t, err)
			attest.True(t, bytes.Equal(out, again), attest.Sprint("MarshalStable produced unstable output"))
		}
	})
	t.Run("not protobuf", func(t *testing.T) {
		_, err := codec.Marshal(struct{}{})
		attest.Error(t, err)
		_, err = codec.MarshalStable(struct{}{})
		attest.Error(t, err)
	})
}

func TestBinaryMetadata(t *testing.T) {
	codec := &binaryCodec{name: "proto"}
	attest.Equal(t, codec.Name(), codec.name)
	attest.True(t, codec.IsBinary())
}
