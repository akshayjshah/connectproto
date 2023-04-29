package connectproto

import (
	"testing"

	"github.com/akshayjshah/attest"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestJSONUnmarshal(t *testing.T) {
	codec := &codec{
		name:      "json",
		unmarshal: protojson.UnmarshalOptions{DiscardUnknown: true},
	}
	t.Run("unknown field", func(t *testing.T) {
		var enum descriptorpb.EnumDescriptorProto
		err := codec.Unmarshal([]byte(`{"name": "foo", "unknown_": "bar"}`), &enum)
		attest.Ok(t, err)
		attest.Equal(t, *enum.Name, "foo")
	})
	t.Run("empty input", func(t *testing.T) {
		var empty emptypb.Empty
		err := codec.Unmarshal([]byte{}, &empty)
		attest.Error(t, err)
		attest.Subsequence(t, err.Error(), "valid JSON")
	})
	t.Run("not protobuf", func(t *testing.T) {
		err := codec.Unmarshal([]byte(`{}`), &struct{}{})
		attest.Error(t, err)
	})
}

func TestJSONMarshal(t *testing.T) {
	codec := &codec{
		name:    "json",
		marshal: protojson.MarshalOptions{Multiline: true, Indent: "\t"},
	}
	enum := &descriptorpb.EnumDescriptorProto{Name: ptr("Foo")}
	t.Run("unstable", func(t *testing.T) {
		out, err := codec.Marshal(enum)
		attest.Ok(t, err)
		attest.Subsequence(t, string(out), "\n\t")
	})
	t.Run("stable", func(t *testing.T) {
		out, err := codec.MarshalStable(enum)
		attest.Ok(t, err)
		attest.Equal(t, string(out), `{"name":"Foo"}`)
	})
	t.Run("not protobuf", func(t *testing.T) {
		_, err := codec.Marshal(struct{}{})
		attest.Error(t, err)
		_, err = codec.MarshalStable(struct{}{})
		attest.Error(t, err)
	})
}

func TestJSONMetadata(t *testing.T) {
	codec := &codec{name: "json"}
	attest.Equal(t, codec.Name(), codec.name)
	attest.False(t, codec.IsBinary())
}

func ptr[T any](val T) *T {
	return &val
}
