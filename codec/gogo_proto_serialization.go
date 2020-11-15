package codec

import (
	"errors"

	"github.com/gogo/protobuf/proto"
)

// GOGOPBSerialization 序列化protobuf包体
// 此序列化方式大概是谷歌的两倍左右 @link:https://github.com/gogo/protobuf
type GOGOPBSerialization struct{}

// Unmarshal 反序列protobuf
func (s *GOGOPBSerialization) Unmarshal(in []byte, body interface{}) error {
	msg, ok := body.(proto.Message)
	if !ok {
		return errors.New("unmarshal fail: body not protobuf message")
	}
	return proto.Unmarshal(in, msg)
}

// Marshal 序列化protobuf
func (s *GOGOPBSerialization) Marshal(body interface{}) ([]byte, error) {
	msg, ok := body.(proto.Message)
	if !ok {
		return nil, errors.New("marshal fail: body not protobuf message")
	}
	return proto.Marshal(msg), nil
}
