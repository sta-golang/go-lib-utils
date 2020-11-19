package codec

import (
	"errors"

	"github.com/gogo/protobuf/proto"
)

type protoCodec struct {
	helper ProtoHelper
}

var protoAPI protoCodec

type ProtoHelper interface {
	Marshal(body interface{}) ([]byte, error)
	Unmarshal(in []byte, body interface{}) error
}

func (pc *protoCodec) Unmarshal(in []byte, body interface{}) error {
	return pc.helper.Unmarshal(in, body)
}

func (pc *protoCodec) Marshal(body interface{}) ([]byte, error) {
	return pc.helper.Marshal(body)
}

func (pc *protoCodec) ReplaceProtoCodec(helper ProtoHelper) {
	if helper == nil {
		return
	}
	if _, ok := helper.(*protoCodec); ok {
		return
	}
	pc.helper = helper
}

// GOGOPBSerialization 序列化protobuf包体
// 此序列化方式大概是谷歌的两倍左右 @link:https://github.com/gogo/protobuf
type goGOPBSerialization struct{}

// Unmarshal 反序列protobuf
func (s *goGOPBSerialization) Unmarshal(in []byte, body interface{}) error {
	msg, ok := body.(proto.Message)
	if !ok {
		return errors.New("unmarshal fail: body not protobuf message")
	}
	return proto.Unmarshal(in, msg)
}

// Marshal 序列化protobuf
func (s *goGOPBSerialization) Marshal(body interface{}) ([]byte, error) {
	msg, ok := body.(proto.Message)
	if !ok {
		return nil, errors.New("marshal fail: body not protobuf message")
	}
	return proto.Marshal(msg)
}
