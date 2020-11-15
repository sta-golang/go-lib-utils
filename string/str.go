package string

import (
	"bytes"
	"hash/crc32"
	"io/ioutil"
	"strings"
	"unsafe"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"golang.org/x/text/width"
)

// BytesToString 0copy 字节数组转化成String
// ZeroCopy 不产生内存拷贝
func BytesToString(bys []byte) string {
	if bys == nil || len(bys) <= 0 {
		return ""
	}
	return *(*string)(unsafe.Pointer(&bys))
}

// StringToBytes 0copy将String转换成bytes数组
// ZeroCopy 不产生内存拷贝
func StringToBytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s)) // 获取s的起始地址开始后的两个 uintptr 指针
	h := [3]uintptr{x[0], x[1], x[1]}      // 构造三个指针数组
	return *(*[]byte)(unsafe.Pointer(&h))
}

// CRC32 计算一个字符串的 crc32 多项式.
func CRC32(str string) uint32 {
	return crc32.ChecksumIEEE(StringToBytes(str))
}

// SBC2DBC 全角转半角
func SBC2DBC(str string) string {
	return width.Narrow.String(str)
}

// DBC2SBC 半角转全角.
func DBC2SBC(str string) string {
	return width.Widen.String(str)
}

// Utf8ToGbk UTF-8转GBK编码.
func Utf8ToGbk(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e := ioutil.ReadAll(reader)
	return d, e
}

// GbkToUtf8 GBK转UTF-8编码.
func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	return d, e
}

// IsEmpty 字符串是否为空(包括空格).
func IsEmpty(str string) bool {
	if str == "" || len(Trim(str)) == 0 {
		return true
	}

	return false
}

// Trim 去除字符串首尾处的空白字符（或者其他字符）.
func Trim(str string, characterMask ...string) string {
	mask := getTrimMask(characterMask)
	return strings.Trim(str, mask)
}

// Ltrim 删除字符串开头的空白字符（或其他字符）.
func Ltrim(str string, characterMask ...string) string {
	mask := getTrimMask(characterMask)
	return strings.TrimLeft(str, mask)
}

// Rtrim 删除字符串末端的空白字符（或者其他字符）.
func Rtrim(str string, characterMask ...string) string {
	mask := getTrimMask(characterMask)
	return strings.TrimRight(str, mask)
}

// getTrimMask 去除mask字符.
func getTrimMask(characterMask []string) string {
	var mask string
	if len(characterMask) == 0 {
		mask = " \t\n\r\v\f\x00　"
	} else {
		mask = strings.Join(characterMask, "")
	}
	return mask
}
