package codec

type codeC struct {
	JsonAPI jsonCodec
	ProtoAPI protoCodec
}

var API = &codeC{}

func init() {
	jsonAPI = jsonCodec{helper:&defaultJson{}}
	protoAPI = protoCodec{helper:&goGOPBSerialization{}}
	API.JsonAPI = jsonAPI
	API.ProtoAPI = protoAPI
}