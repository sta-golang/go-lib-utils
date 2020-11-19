package codec

import (
	"crypto/md5"
	"encoding/hex"
	str "github.com/xy63237777/go-lib-utils/str"
	"strings"
)

type cryptoCodec struct {
	helper cryptoHelper
}

var cryptoAPI cryptoCodec

type cryptoHelper interface {
	Encode(data string) (string, error)
	Check(content, encrypted string) bool
}

func (cc *cryptoCodec) Encode(data string) (string, error) {
	return cc.helper.Encode(data)
}

func (cc *cryptoCodec) Check(content, encrypted string) bool {
	return cc.helper.Check(content, encrypted)
}

func (cc *cryptoCodec) ReplaceCryptoCodec(helper cryptoHelper) {
	if helper == nil {
		return
	}
	if _, ok := helper.(*cryptoCodec); ok {
		return
	}
	cc.helper = helper
}

type md5Crypto struct {
}

func (mc *md5Crypto) Encode(data string) (string, error) {
	h := md5.New()
	_, err := h.Write(str.StringToBytes(data))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func (mc *md5Crypto) Check(content, encrypted string) bool {
	target, err := mc.Encode(content)
	if err != nil {
		return false
	}
	return strings.EqualFold(target, encrypted)
}
