package lcache

import (
	"bufio"
	cache "github.com/patrickmn/go-cache"
	"log"
	"os"
	"time"
)

var LCache *cache.Cache

func init() {
	LCache = cache.New(5*time.Second, 60*time.Second)
}

func SetCahce(k string, x interface{}, d time.Duration) {
	LCache.Set(k, x, d)
}

func GetCache(k string) (interface{}, bool) {
	return LCache.Get(k)
}

// 设置cache 无时间参数
func SetDefaultCahce(k string, x interface{}) {
	LCache.SetDefault(k, x)
}

// 删除 cache
func DeleteCache(k string) {
	LCache.Delete(k)
}

// Add() 加入缓存
func AddCache(k string, x interface{}, d time.Duration) {
	LCache.Add(k, x, d)
}

// IncrementInt() 对已存在的key 值自增n
func IncrementIntCahce(k string, n int) (num int, err error) {
	return LCache.IncrementInt(k, n)
}

func Permanent() {
	logf, err := os.Open("/var/cache/des")
	if err != nil {
		log.Fatalln(err)
	}
	writer := bufio.NewWriter(logf)
	LCache.Save(writer)
}

func Recover() {
	logf, err := os.Open("/var/cache/des")
	if err != nil {
		log.Fatalln(err)
	}
	reader := bufio.NewReader(logf)
	LCache.Load(reader)
}
