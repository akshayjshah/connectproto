package connectproto

import (
	"bytes"
	"testing"
	"time"

	"go.akshayshah.org/attest"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
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
	codec := newBinaryCodec(proto.MarshalOptions{}, proto.UnmarshalOptions{})
	attest.Equal(t, codec.Name(), "proto")
	attest.True(t, codec.IsBinary())
}

func TestVTUnmarshal(t *testing.T) {
	codec := newBinaryVTCodec()
	pb := timestamppb.Now()
	bin, err := proto.Marshal(pb)
	attest.Ok(t, err)
	var vt timestampVT
	attest.Ok(t, codec.Unmarshal(bin, &vt))
	attest.Equal(t, vt.t, pb.AsTime())
	// also ensure codec can unmarshal to non-VT structs
	attest.Ok(t, codec.Unmarshal(bin, &timestamppb.Timestamp{}))
}

func TestVTMarshal(t *testing.T) {
	codec := newBinaryVTCodec()
	vt := &timestampVT{time.Now()}
	bin, err := codec.Marshal(vt)
	attest.Ok(t, err)
	var pb timestamppb.Timestamp
	attest.Ok(t, proto.Unmarshal(bin, &pb))
	attest.Equal(t, pb.AsTime(), vt.t)
	// also ensure codec can marshal non-VT structs
	_, err = codec.Marshal(&timestamppb.Timestamp{})
	attest.Ok(t, err)
}

func TestVTMetadata(t *testing.T) {
	codec := newBinaryVTCodec()
	attest.Equal(t, codec.Name(), "proto")
	attest.True(t, codec.IsBinary())
}

type timestampVT struct {
	t time.Time
}

func (t *timestampVT) MarshalVT() ([]byte, error) {
	msg := timestamppb.New(t.t)
	return proto.Marshal(msg)
}

func (t *timestampVT) UnmarshalVT(binary []byte) error {
	var msg timestamppb.Timestamp
	if err := proto.Unmarshal(binary, &msg); err != nil {
		return err
	}
	t.t = msg.AsTime()
	return nil
}
