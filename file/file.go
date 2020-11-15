package file

import (
	"os"
	"path/filepath"
	"strings"
)

// FileSize 获取文件大小
// @Returns 1: 文件大小 2: 文件是否存在  3: error
func FileSize(filePath string) (int64, bool, error) {
	file, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, false, nil
		}
		return 0, false, err
	}
	return file.Size(), true, nil
}

// GetExt 获取文件的小写扩展名,不包括点"." .
func GetExt(fpath string) string {
	suffix := filepath.Ext(fpath)
	if suffix != "" {
		return strings.ToLower(suffix[1:])
	}
	return suffix
}