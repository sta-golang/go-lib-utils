package log

import (
	"bytes"
	"fmt"
	"github.com/sta-golang/go-lib-utils/err"
	"github.com/sta-golang/go-lib-utils/file"
	"github.com/sta-golang/go-lib-utils/source"
	"github.com/sta-golang/go-lib-utils/str"
	tm "github.com/sta-golang/go-lib-utils/time"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	reOpenFileTime = 10
	logSuffix      = "log"
	mb             = 1024 * 1024
	defMaxSize     = 16 * mb
	allLevel       = "ALL"
	defFileDir     = "./log"
	minFileSize    = 10 * mb
	fileLogName    = "sta:file_log"

	syncFlagInit    = 0
	syncFlagProcess = 1
)

var syncFlag = &asyncNode{}
var reSizeFlag = true

func init() {
	if runtime.GOOS == "windows" {
		reSizeFlag = false
	}
}

type fileLog struct {
	writerHelper *fileLogWriter
	asyncWriter  *asyncHelper
	level        Level
	prefix       string
	syncFlag     int32
	skip         int
}

func (fl *fileLog) setSkip(skip int) {
	fl.skip = skip
}

func (fl *fileLog) SetLevel(level Level) {
	if level < DEBUG || level > FATAL {
		return
	}
	fl.level = level
}

func (fl *fileLog) Name() string {
	return fileLogName
}

func (fl *fileLog) Sync() error {
	fl.asyncWriter.ch <- syncFlag
	for atomic.LoadInt32(&fl.syncFlag) != syncFlagProcess {
		time.Sleep(time.Millisecond * 50)
	}
	atomic.StoreInt32(&fl.syncFlag, syncFlagInit)
	return nil
}

// GetLevel 获取输出端日志级别
func (fl *fileLog) GetLevel() string {
	return levelFlages[fl.level]
}

func (fl *fileLog) Debugf(format string, args ...interface{}) {
	fl.print(DEBUG, format, args...)
}

func (fl *fileLog) Warnf(format string, args ...interface{}) {
	fl.print(WARNING, format, args...)
}

func (fl *fileLog) Infof(format string, args ...interface{}) {
	fl.print(INFO, format, args...)
}

func (fl *fileLog) Errorf(format string, args ...interface{}) {
	fl.print(ERROR, format, args...)
}

func (fl *fileLog) Fatalf(format string, args ...interface{}) {
	fl.print(FATAL, format, args...)
}

func (fl *fileLog) Debug(args ...interface{}) {
	fl.println(DEBUG, args...)
}

func (fl *fileLog) Warn(args ...interface{}) {
	fl.println(WARNING, args...)
}

func (fl *fileLog) Info(args ...interface{}) {
	fl.println(INFO, args...)
}

func (fl *fileLog) Error(args ...interface{}) {
	fl.println(ERROR, args...)
}

func (fl *fileLog) Fatal(args ...interface{}) {
	fl.println(FATAL, args...)
}

func (fl *fileLog) print(level Level, format string, args ...interface{}) {
	if level < fl.level {
		return
	}

	logFmt := "%s %s [%s] %s ==> %s\n"
	_, transFile, transLine, _ := runtime.Caller(fl.skip)
	data := fmt.Sprintf(logFmt, fl.prefix, tm.GetNowDateTimeStr(), levelFlages[level],
		fmt.Sprintf("%s:%d", transFile, transLine), fmt.Sprintf(format, args...))
	if fl.asyncWriterFunc(level, &data) {
		return
	}
	fl.writerHelper.writer(level, str.StringToBytes(&data))
}

func (fl *fileLog) println(level Level, args ...interface{}) {
	if level < fl.level {
		return
	}
	logFmt := "%s %s [%s] %s ==> "
	_, transFile, transLine, _ := runtime.Caller(fl.skip)
	data := fmt.Sprintf("%s%s", fmt.Sprintf(logFmt, fl.prefix, tm.GetNowDateTimeStr(), levelFlages[level],
		fmt.Sprintf("%s:%d", transFile, transLine)),
		fmt.Sprintln(args...))
	if fl.asyncWriterFunc(level, &data) {
		return
	}
	fl.writerHelper.writer(level, str.StringToBytes(&data))
}

func (fl *fileLog) asyncWriterFunc(level Level, data *string) bool {
	if fl.asyncWriter != nil && level <= fl.asyncWriter.level {
		if node, ok := fl.asyncWriter.pool.Get().(*asyncNode); ok {
			node.level = level
			node.bys = str.StringToBytes(data)

			fl.asyncWriter.ch <- node
			return true
		}
	}
	return false
}

type FileLogConfig struct {
	FileDir     string   `yaml:"file_dir"`     // 文件目录 默认为./log
	FileName    string   `yaml:"file_name"`    // 文件前缀名 默认为sta
	DayAge      int      `yaml:"day_age"`      // 文件保留日期 默认为7天
	LogLevel    Level    `yaml:"log_level"`    // 日志等级 默认为INFO 这个可以之后设置
	Prefix      string   `yaml:"pre_fix"`      // 日志输出前缀 默认为 FOUR-SEASONS: STA
	MaxSize     int64    `yaml:"max_size"`     // 最大大小 设置0或者辅助 默认失效。 最小为10mb 如果小于10mb则变成16mb
	AloneWriter []string `yaml:"alone_writer"` // 单独数组的等级，设置后没有出现的向等级低的方向靠
}

var DefaultFileLogConfig = func() *FileLogConfig {
	return &FileLogConfig{
		FileDir:     "./log",
		FileName:    "sta",
		DayAge:      7,
		LogLevel:    INFO,
		MaxSize:     0,
		AloneWriter: nil,
		Prefix:      PREFIX,
	}
}

var DefaultFileLogConfigForAloneWriter = func(alone []string) *FileLogConfig {
	return &FileLogConfig{
		FileDir:     "./log",
		FileName:    "sta",
		DayAge:      7,
		LogLevel:    INFO,
		MaxSize:     0,
		AloneWriter: alone,
		Prefix:      PREFIX,
	}
}

func NewFileLogAndAsync(conf *FileLogConfig, syncTime time.Duration) Logger {
	return newFileLog(conf, true, syncTime)
}

func NewFileLog(conf *FileLogConfig) Logger {
	return newFileLog(conf, false, time.Millisecond)
}

func newFileLog(conf *FileLogConfig, async bool, asyncTime time.Duration) Logger {
	conf.fixConfig()
	table := make(map[string]Level)
	for _, key := range conf.AloneWriter {
		if index, ok := levelIndexs[key]; ok {
			table[key] = index
		}
	}
	if conf.MaxSize > 0 && conf.MaxSize < minFileSize {
		conf.MaxSize = defMaxSize
	}
	ret := &fileLog{
		writerHelper: &fileLogWriter{
			helpers:  nil,
			fileDir:  conf.FileDir,
			closeCh:  make(chan *os.File, 16),
			dayAge:   conf.DayAge,
			fileName: conf.FileName,
			maxSize:  conf.MaxSize,
		},
		level:    conf.LogLevel,
		skip:     dfsStep,
		prefix:   conf.Prefix,
		syncFlag: syncFlagInit,
	}
	helpers := make([]*writerHelper, len(levelFlages))
	if len(table) == len(levelFlages) || len(table) == len(levelFlages)-1 {
		for i := 0; i < len(helpers); i++ {
			helpers[i] = &writerHelper{
				level:      levelFlages[i],
				openFile:   atomic.Value{},
				openDate:   "",
				openTime:   0,
				writerSize: 0,
				lock:       sync.Mutex{},
				target:     ret.writerHelper,
			}
		}
	} else if len(table) == 0 {
		oneHelper := &writerHelper{
			level:      allLevel,
			openFile:   atomic.Value{},
			openDate:   "",
			openTime:   0,
			writerSize: 0,
			lock:       sync.Mutex{},
			target:     ret.writerHelper,
		}
		for i := 0; i < len(levelFlages); i++ {
			helpers[i] = oneHelper
		}
	} else {
		levelArr := make([]string, len(levelFlages))
		for i := len(levelFlages) - 1; i >= 0; i-- {
			levelArr[i] = allLevel
			if _, ok := table[levelFlages[i]]; ok {
				levelArr[i] = levelFlages[i]
				for j := i + 1; j < len(levelFlages) && levelArr[j] == allLevel; j++ {
					levelArr[j] = levelFlages[i]
				}
			}
		}
		helpers[0] = &writerHelper{
			level:      levelArr[0],
			openFile:   atomic.Value{},
			openDate:   "",
			openTime:   0,
			writerSize: 0,
			lock:       sync.Mutex{},
			target:     ret.writerHelper,
		}
		for i := 1; i < len(levelFlages); i++ {
			if levelArr[i] == helpers[i-1].level {
				helpers[i] = helpers[i-1]
			} else {
				helpers[i] = &writerHelper{
					level:      levelArr[i],
					openFile:   atomic.Value{},
					openDate:   "",
					openTime:   0,
					writerSize: 0,
					lock:       sync.Mutex{},
					target:     ret.writerHelper,
				}
			}
		}

	}

	ret.writerHelper.helpers = helpers
	if async {
		le := WARNING
		if _, ok := table[levelFlages[WARNING]]; ok {
			le = INFO
		}
		ret.asyncWriter = &asyncHelper{
			target:   ret,
			ch:       make(chan *asyncNode, 128),
			syncTime: asyncTime,
			pool: sync.Pool{
				New: func() interface{} {
					return &asyncNode{}
				},
			},
			level: le,
		}
		go ret.asyncWriter.worker()
	}
	source.Monitoring(ret)
	go ret.writerHelper.asyncCloseFiles()
	return ret
}

type fileLogWriter struct {
	helpers  []*writerHelper
	fileDir  string
	closeCh  chan *os.File
	dayAge   int
	fileName string
	maxSize  int64
	initFlag bool
}

// writerHelper
// 文件名为 {target.fileName}.log.{openDate}.{level}.{numberSuffix}
// 例如 sta.log.2020-11-25.INFO.1
type writerHelper struct {
	level      string
	openFile   atomic.Value // 文件句柄
	openDate   string
	openTime   int64
	writerSize int64
	lock       sync.Mutex
	target     *fileLogWriter
}

func (fl *fileLogWriter) asyncCloseFiles() {
	for fi := range fl.closeCh {
		time.Sleep(time.Millisecond * 30)
		e := fi.Close()
		if e != nil {
			ConsoleLogger.Error(err.NewError(err.LogErrCode+FileCloseError,
				fmt.Errorf("%s file clouse Err", fi.Name())))
		}

	}
}

func (fl *fileLogWriter) cleanFiles() {
	names, er := fl.getNeedCleanFileNames()
	if er != nil {
		ConsoleLogger.Error(err.NewError(err.LogErrCode+CurrencyError, fmt.Errorf("%v", er)))
	}
	for i := range names {
		fileN := fmt.Sprintf("%s/%s", fl.fileDir, names[i])
		_ = os.Remove(fileN)
	}
}

func (fl *fileLogWriter) writer(level Level, data []byte) {
	er := fl.helpers[level].doWriter(data)
	if er != nil {
		ConsoleLogger.Errorf("level : %s writer Err %v", levelFlages[level], er)
	}
}

func (wh *writerHelper) doWriter(data []byte) *err.Error {
	openFile := wh.getOpenFile()

	if openFile == nil || tm.GetNowTime().Unix()-atomic.LoadInt64(&wh.openTime) >= reOpenFileTime {
		wh.lock.Lock()
		wh.reLoadFile()
		wh.lock.Unlock()
	}
	openFile = wh.getOpenFile()
	if openFile == nil {
		return err.NewError(err.LogErrCode+FileCreateError, fmt.Errorf("open %s File err", wh.getFilePath()))
	}
	n, er := openFile.Write(data)

	atomic.AddInt64(&wh.writerSize, int64(n))
	if er != nil {
		return err.NewError(err.LogErrCode+FileWriterError, fmt.Errorf("writer err %v", er))
	}
	if reSizeFlag && wh.target.maxSize > 0 && atomic.LoadInt64(&wh.writerSize) > wh.target.maxSize {
		wh.lock.Lock()
		wh.reCreateFile()
		wh.lock.Unlock()
	}
	return nil
}

func (wh *writerHelper) getFilePath() string {
	return fmt.Sprintf("%s/%s.%s.%s.%s", wh.target.fileDir, wh.target.fileName,
		logSuffix, wh.openDate, wh.level)

}

func (wh *writerHelper) reCreateFile() {
	if wh.target.maxSize > 0 && atomic.LoadInt64(&wh.writerSize) >= wh.target.maxSize {

		atomic.StoreInt64(&wh.target.maxSize, 0)
		// 旧文件改名
		oldPath := wh.getFilePath()
		newName := fmt.Sprintf("%s.%s", oldPath, tm.GetNowTimeStr())
		// 这里适用于linux 不适用于windows
		if _, e := os.Stat(oldPath); !os.IsNotExist(e) {
			_ = os.Rename(oldPath, newName)
		}

		wh.reLoadFile()
	}
}

func (wh *writerHelper) doReLoadFile() {
	if !wh.target.initFlag {
		er := os.MkdirAll(wh.target.fileDir, 0666)
		if er != nil {
			ConsoleLogger.Fatalf("create Dir : %s fatal err %v", wh.target.fileDir, er)
			panic(er)
		}
		wh.target.initFlag = true
	}
	lastFile := wh.getOpenFile()
	lastDate := wh.openDate
	atomic.StoreInt64(&wh.openTime, tm.GetNowTime().Unix())
	wh.openDate = tm.GetNowDateStr()

	if wh.target.dayAge > 0 && lastDate != wh.openDate {
		go wh.target.cleanFiles()
	}

	path := wh.getFilePath()
	of, er := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if er != nil {
		ConsoleLogger.Error(err.NewError(err.LogErrCode+FileWriterError,
			fmt.Errorf("open file : %v err %v", path, er)))
		return
	}
	wh.openFile.Store(of)
	if lastFile != nil {
		wh.target.closeCh <- lastFile
	}

	size, flag, _ := file.FileSize(path)
	if flag {
		atomic.StoreInt64(&wh.writerSize, size)
	}
}

func (wh *writerHelper) reLoadFile() {
	if wh.getOpenFile() == nil || tm.GetNowTime().Unix()-atomic.LoadInt64(&wh.openTime) >= reOpenFileTime {
		wh.doReLoadFile()
	}
}

func (wh *writerHelper) getOpenFile() *os.File {
	if fi, ok := wh.openFile.Load().(*os.File); ok {
		return fi
	}
	return nil
}

func (fl *fileLogWriter) getNeedCleanFileNames() ([]string, *err.Error) {
	files, er := ioutil.ReadDir(fl.fileDir)
	if er != nil {
		return nil, err.NewError(err.LogErrCode+FileReadDirError,
			fmt.Errorf("filedir %v readDir hasErr %v", fl.fileDir, er))
	}
	ret := make([]string, 0)
	for _, f := range files {
		if f.IsDir() || !strings.HasPrefix(f.Name(), fl.fileName) {
			continue
		}
		fileTime, er := getLogFileTime(f.Name())
		if er != nil {
			ConsoleLogger.Error(err.NewError(err.LogErrCode+CurrencyError,
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
	timeStr := name[index+1 : index+len(tm.LayoutDate)+1]
	return tm.ParseDate(timeStr)
}

type asyncHelper struct {
	target   *fileLog
	ch       chan *asyncNode
	syncTime time.Duration
	level    Level
	pool     sync.Pool
}

type asyncNode struct {
	level Level
	bys   []byte
}

func (ah *asyncHelper) worker() {
	tick := time.Tick(ah.syncTime)
	buffers := make([]*bytes.Buffer, ah.level+1)
	buffers[0] = &bytes.Buffer{}
	for i := 1; i < len(buffers); i++ {
		if ah.target.writerHelper.helpers[i] == ah.target.writerHelper.helpers[i-1] {
			buffers[i] = buffers[i-1]
		} else {
			buffers[i] = &bytes.Buffer{}
		}
	}
	syncDoWriter := func() {
		for i := range buffers {
			if buffers[i].Len() > 0 {
				ah.target.writerHelper.writer(Level(i), buffers[i].Bytes())
				buffers[i].Reset()
			}
		}
	}
	for {
		select {
		case <-tick:
			syncDoWriter()
		case ele := <-ah.ch:
			if ele == syncFlag {
				syncDoWriter()
				atomic.StoreInt32(&ah.target.syncFlag, syncFlagProcess)
			} else {
				_, _ = buffers[ele.level].Write(ele.bys)
				if buffers[ele.level].Len() > (1024<<2)-512 {
					syncDoWriter()
				}
				ah.pool.Put(ele)
			}
		}
	}
}

func (flc *FileLogConfig) fixConfig() {
	if flc.MaxSize < 0 {
		flc.MaxSize = 0
	}
	if flc.MaxSize > 0 && flc.MaxSize < minFileSize {
		flc.MaxSize = defMaxSize
	}
	if flc.LogLevel < 0 || int(flc.LogLevel) > len(levelFlages) {
		flc.LogLevel = INFO
	}
	if flc.DayAge < 0 || flc.DayAge > 360 {
		flc.DayAge = 0
	}
	if flc.Prefix == "" {
		flc.Prefix = PREFIX
	}
	if flc.FileDir == "" {
		flc.FileDir = "./log"
	}
	if flc.FileName == "" {
		flc.FileName = "sta"
	}
}
