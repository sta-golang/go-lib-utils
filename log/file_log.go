package log

import (
	"errors"
	"fmt"
	"github.com/xy63237777/go-lib-utils/err"
	tm "github.com/xy63237777/go-lib-utils/time"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	reOpenFileTime = 10
	logSuffix      = "log"
	separator      = "."
)

type ExistenceType int

type fileLog struct {
	helpers []*writerHelper
	fileDir string
	closeCh chan *os.File
	dayAge  int64
}

type writerHelper struct {
	level        string
	openFile     atomic.Value // 文件句柄
	openTime     int64
	openFileName string
	writerSize   int64
	lock         sync.Mutex
	closeFlag    bool
	target       *fileLog
	maxSize      int64
}

func (fl *fileLog) asyncCloseFiles() {
	for file := range fl.closeCh {
		time.Sleep(time.Millisecond * 30)
		e := file.Close()
		if e != nil {
			FrameworkLogger.Error(err.NewError(err.LogErrCode+FileCloseError,
				fmt.Errorf("%s file clouse Err", file.Name())))
		}
	}
}

func (fl *fileLog) cleanFiles() {
	names, er := fl.getNeedCleanFileNames()
	if er != nil {
		FrameworkLogger.Error(err.NewError(err.LogErrCode+CurrencyError, fmt.Errorf("%v", er)))
	}

}

func (fl *fileLog) writer(level Level, data []byte) *err.Error {
	return fl.helpers[level].doWriter(data)
}

func (wh *writerHelper) doWriter(data []byte) *err.Error {
	openFile := wh.getOpenFile()
	if openFile == nil || time.Now().Unix()-atomic.LoadInt64(&wh.openTime) >= reOpenFileTime {

	}
	if openFile == nil {
		return err.NewError(err.LogErrCode+FileCreateError, fmt.Errorf("open %s File err"))
	}
	n, err := openFile.Write(data)
	atomic.AddInt64(&wh.writerSize, int64(n))
	if err != nil {

	}
}

func (wh *writerHelper) getOpenFile() *os.File {
	if file, ok := wh.openFile.Load().(*os.File); ok {
		return file
	}
	return nil
}

func (fl *fileLog) getNeedCleanFileNames() ([]string, *err.Error) {
	files, er := ioutil.ReadDir(fl.fileDir)
	if er != nil {
		return nil, err.NewError(err.LogErrCode+FileReadDirError,
			fmt.Errorf("filedir %v readDir hasErr %v", fl.fileDir, er))
	}
	ret := make([]string, 0)
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		fileTime, er := getLogFileTime(f.Name())
		if er != nil {
			FrameworkLogger.Error(err.NewError(err.LogErrCode+CurrencyError,
				fmt.Errorf("%v", er)))
			continue
		}
		if tm.GetNowTime().Sub(fileTime).Hours() > float64(fl.dayAge*24) {
			ret = append(ret, f.Name())
		}
	}
	return ret, nil
}

func getLogFileTime(name string) (time.Time, error) {
	index := -1
	for i := 4; i < len(name); i++ {
		if name[i] == '.' && name[i-1] == 'g' && name[i-2] == 'o' && name[i-3] == 'l' && name[i-4] == '.' {
			index = i
			break
		}
	}
	if index == -1 {
		return time.Time{}, fmt.Errorf("name %s fmt Err", name)
	}
	timeStr := name[index+1 : len(tm.LayoutDate)+1]
	return tm.ParseDate(timeStr)
}
