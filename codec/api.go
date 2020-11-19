package codec

type codeC struct {
	JsonAPI   jsonCodec
	ProtoAPI  protoCodec
	CryptoAPI cryptoCodec
}

var API = &codeC{}

func init() {
	jsonAPI = jsonCodec{helper: &defaultJson{}}
	protoAPI = protoCodec{helper: &goGOPBSerialization{}}
	cryptoAPI = cryptoCodec{helper: &md5Crypto{}}
	API.JsonAPI = jsonAPI
	API.ProtoAPI = protoAPI
	API.CryptoAPI = cryptoAPI
}
