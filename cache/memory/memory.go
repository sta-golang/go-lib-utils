package memory

import (
	"github.com/xy63237777/go-lib-utils/cache"
	"github.com/xy63237777/go-lib-utils/log"
	"github.com/xy63237777/go-lib-utils/str"
	"hash/crc32"
	"sync"
	"sync/atomic"
	"time"
)

// 内部chanSize定义
const (
	defChanSize = 4096

	defShards     = 4
	defMaxKey     = 100000
	defGCInterval = 60
	noExpire      = -1
)

// memory 内存变量缓存组件，只能设置带过期时间的内存变量，
// 只有过期策略，没有淘汰策略，当保存的变量个数超过最大限制，
// 就不在缓存新的变量，直到有变量过期删除
type Memory struct {
	cacheMap      []map[string]*memoryItem
	mapRWMutex    []sync.RWMutex
	usedKey       int32
	cfg           Config
	setChan       chan *passNode // 写channel
	gcTree        GCHelper       // gc树
	cacheItemPool sync.Pool      // 缓存对象池
	passNodePool  sync.Pool      // 异步过期管理对象池
}

type Config struct {
	Shards      int `yaml:"shards"`
	GCInterval  int `yaml:"gc_interval"`
	MaxKey      int `yaml:"max_Key"`
	SetChanSize int `yaml:"set_Chan_Size"`
}

func NewConfig(shards, gcInter, maxKey, setChanSize int) Config {
	cfg := Config{
		Shards:      shards,
		GCInterval:  gcInter,
		MaxKey:      maxKey,
		SetChanSize: setChanSize,
	}
	if cfg.Shards <= 0 {
		cfg.Shards = defShards
	}
	if cfg.GCInterval <= 0 {
		cfg.GCInterval = defGCInterval
	} else if cfg.GCInterval <= 10 {
		cfg.GCInterval = defGCInterval >> 1
	}
	if cfg.MaxKey <= 0 {
		cfg.MaxKey = defMaxKey
	} else if cfg.MaxKey <= 1000 {
		cfg.MaxKey = defMaxKey / 10
	}
	if cfg.SetChanSize <= 0 {
		cfg.SetChanSize = defChanSize
	} else if cfg.SetChanSize <= 1024 {
		cfg.SetChanSize = defChanSize >> 1
	}
	return cfg
}

func NewDefaultConfig() Config {
	return Config{
		Shards:      defShards,
		GCInterval:  defGCInterval,
		MaxKey:      defMaxKey,
		SetChanSize: defChanSize,
	}
}

type memoryItem struct {
	key    string
	val    interface{}
	expire int64
}

// passNode chan异步管理key移除时间
type passNode struct {
	remove int64
	key    string
}

// New 创建缓存对象
var New = func(cfg Config) cache.Cache {
	return NewWithGCHelper(cfg, newTree())
}

// NewWithGCHelper 创建缓存对象
var NewWithGCHelper = func(cfg Config, helper GCHelper) cache.Cache {

	cacheMap := make([]map[string]*memoryItem, cfg.Shards)
	for i := 0; i < cfg.Shards; i++ {
		cacheMap[i] = make(map[string]*memoryItem)
	}
	cache := &Memory{
		cfg:        cfg,
		cacheMap:   cacheMap,
		mapRWMutex: make([]sync.RWMutex, cfg.Shards),
		usedKey:    0,
		setChan:    make(chan *passNode, cfg.SetChanSize),
		gcTree:     helper,
		cacheItemPool: sync.Pool{
			New: func() interface{} {
				return &memoryItem{}
			},
		},
		passNodePool: sync.Pool{
			New: func() interface{} {
				return &passNode{}
			},
		},
	}
	go cache.setRoutine()
	go cache.gcRoutine()
	return cache
}

// shardsNum 获取key的分片位置
func (m *Memory) shardsNum(s string) uint32 {
	b := str.StringToBytes(&s)
	return crc32.ChecksumIEEE(b) % uint32(m.cfg.Shards)
}

// Set 写入本地缓存
func (m *Memory) Set(key string, val interface{}) bool {
	return m.SetWithRemove(key, val, noExpire)
}

// SetWithRemove 写入本地缓存
func (m *Memory) SetWithRemove(key string, val interface{}, expire int) bool {
	//是否写满
	if atomic.LoadInt32(&m.usedKey) >= int32(m.cfg.MaxKey) {
		return false
	}
	//setChan满了
	if len(m.setChan) >= m.cfg.SetChanSize {
		return false
	}
	//写入缓冲区
	curSecond := time.Now().Unix()
	item, _ := m.cacheItemPool.Get().(*memoryItem)
	item.key = key
	item.val = val
	item.expire = curSecond + int64(expire)
	if expire <= noExpire {
		item.expire = noExpire
	}
	splitNum := m.shardsNum(key)
	m.mapRWMutex[splitNum].Lock()
	m.cacheMap[splitNum][key] = item
	m.mapRWMutex[splitNum].Unlock()
	if expire <= noExpire {
		return true
	}
	//写入channel
	n, _ := m.passNodePool.Get().(*passNode)
	n.remove = item.expire
	n.key = key
	m.setChan <- n
	return true
}

// Get 从本地缓存中读数据，返回3个字段：缓存数据、是否在有效期内（在remove之前仅返回一次false）
func (m *Memory) Get(key string) (interface{}, bool) {
	splitNum := m.shardsNum(key)
	curSecond := time.Now().Unix()
	//读取数据
	m.mapRWMutex[splitNum].RLock()
	if item, ok := m.cacheMap[splitNum][key]; ok {
		m.mapRWMutex[splitNum].RUnlock() //解锁
		if item.expire <= noExpire {
			return item.val, true
		}
		useAble := true
		//过期的时候返回false，并发控制依赖业务控制
		if curSecond > item.expire {
			useAble = false
		}
		return item.val, useAble
	}
	m.mapRWMutex[splitNum].RUnlock() //解锁
	return nil, false
}

var initPassNode passNode

// setRoutine set协程
func (m *Memory) setRoutine() {
	for pNode := range m.setChan {
		m.gcTree.Add(pNode.remove, pNode.key)
		*pNode = initPassNode
		m.passNodePool.Put(pNode)
	}
}

var initCacheItem memoryItem

// gcRoutine gc协程
func (m *Memory) gcRoutine() {
	ticker := time.NewTicker(time.Second * time.Duration(m.cfg.GCInterval))
	for range ticker.C {
		curSecond := time.Now().Unix()
		keys := m.gcTree.Pruning(curSecond)
		if keys == nil {
			continue
		}
		if log.GlobalLevel == log.DEBUG {
			log.Debug("gcTree pruning keys : %v", keys)
		}
		for _, key := range keys {
			splitNum := m.shardsNum(key)
			m.mapRWMutex[splitNum].Lock()
			item, ok := m.cacheMap[splitNum][key]
			if !ok || curSecond < item.expire { //未过期
				m.mapRWMutex[splitNum].Unlock()
				continue
			}
			delete(m.cacheMap[splitNum], key)
			m.mapRWMutex[splitNum].Unlock()
			*item = initCacheItem
			m.cacheItemPool.Put(item)
			//使用计数减一
			atomic.AddInt32(&m.usedKey, -1)
		}
	}
}
