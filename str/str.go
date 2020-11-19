package str

import (
	"hash/crc32"
	"strings"
	"unsafe"
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
