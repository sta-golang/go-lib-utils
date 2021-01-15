package codec

import "gopkg.in/yaml.v2"

type yamlCodec struct {
	helper YamlHelper
}

var yamlAPI yamlCodec

type YamlHelper interface {
	Marshal(v interface{}) ([]byte, error)
	UnMarshal(data []byte, v interface{}) error
}

func (yc *yamlCodec) Marshal(v interface{}) ([]byte, error) {
	return yc.helper.Marshal(v)
}

func (yc *yamlCodec) UnMarshal(data []byte, v interface{}) error {
	return yc.helper.UnMarshal(data, v)
}

func (yc *yamlCodec) ReplaceYamlCodec(helper YamlHelper) {
	if helper == nil {
		return
	}
	if _, ok := helper.(*yamlCodec); ok {
		return
	}
	yc.helper = helper
}

type defaultYaml struct {
}

func (dy *defaultYaml) UnMarshal(data []byte, v interface{}) error {
	return yaml.Unmarshal(data, v)
}

func (dy *defaultYaml) Marshal(v interface{}) ([]byte, error) {
	return yaml.Marshal(v)
}
