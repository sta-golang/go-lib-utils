package codec

import "encoding/json"

type jsonCodec struct {
	helper JsonHelper
}

var jsonAPI jsonCodec

type JsonHelper interface {
	Marshal(v interface{}) ([]byte, error)
	UnMarshal(data []byte, v interface{}) error
}

func (jc *jsonCodec) Marshal(v interface{}) ([]byte, error) {
	return jc.helper.Marshal(v)
}

func (jc *jsonCodec) UnMarshal(data []byte, v interface{}) error {
	return jc.helper.UnMarshal(data, v)
}

func (jc *jsonCodec) ReplaceJsonCodec(helper JsonHelper) {
	if helper == nil {
		return
	}
	if _, ok := helper.(*jsonCodec); ok {
		return
	}
	jc.helper = helper
}

type defaultJson struct {
}

func (dj *defaultJson) Marshal(v interface{}) ([]byte, error) {

	return json.Marshal(v)
}

func (dj *defaultJson) UnMarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
