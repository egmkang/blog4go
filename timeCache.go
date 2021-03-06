// Copyright (c) 2015, huangjunwei <huangjunwei@youmi.net>. All rights reserved.

package blog4go

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

const (
	// PrefixTimeFormat const time format prefix
	PrefixTimeFormat = "[2006/01/02 15:04:05"

	// DateFormat date format
	DateFormat = "2006-01-02"
)

// timeFormatCacheType is a time formated cache
type timeFormatCacheType struct {
	// current time
	now time.Time
	// current date
	date string
	// current formated date
	format []byte
	// yesterdate
	dateYesterday string

	// lock for read && write
	lock *sync.RWMutex

	//millisceonds cache
	seconds      int64
	formatCache  []byte
	milliSeconds [][]byte
}

// global time cache instance used for every log writer
var timeCache = timeFormatCacheType{}

func init() {
	timeCache.lock = new(sync.RWMutex)
	timeCache.now = time.Now()
	timeCache.date = timeCache.now.Format(DateFormat)
	timeCache.format = []byte(timeCache.now.Format(PrefixTimeFormat))
	timeCache.dateYesterday = timeCache.now.Add(-24 * time.Hour).Format(DateFormat)
	initMilliSeconds()

	// update timeCache every seconds
	go func() {
		// tick every seconds
		t := time.Tick(1 * time.Second)

		//UpdateTimeCacheLoop:
		for {
			select {
			case <-t:
				timeCache.fresh()
			}
		}
	}()
}

func initMilliSeconds() {
	timeCache.milliSeconds = make([][]byte, 1024)
	var index = 0
	for {
		if index >= 1024 {
			break
		}
		timeCache.milliSeconds[index] = []byte(fmt.Sprintf(".%03d]", index))
		index++
	}
}

// Now now
func (timeCache *timeFormatCacheType) Now() time.Time {
	timeCache.lock.RLock()
	defer timeCache.lock.RUnlock()
	return timeCache.now
}

// Date date
func (timeCache *timeFormatCacheType) Date() string {
	timeCache.lock.RLock()
	defer timeCache.lock.RUnlock()
	return timeCache.date
}

// DateYesterday date
func (timeCache *timeFormatCacheType) DateYesterday() string {
	timeCache.lock.RLock()
	defer timeCache.lock.RUnlock()
	return timeCache.dateYesterday
}

// Format format
func (timeCache *timeFormatCacheType) Format() ([]byte, []byte) {
	now := time.Now()
	oldValue := atomic.LoadInt64(&timeCache.seconds)
	newValue := now.Unix()
	milliSeconds := now.Nanosecond() / 1000 / 1000
	milliSecondsFormat := timeCache.milliSeconds[milliSeconds%1024]

	if oldValue != newValue {
		format := []byte(now.Format(PrefixTimeFormat))
		if atomic.CompareAndSwapInt64(&timeCache.seconds, oldValue, newValue) {
			timeCache.formatCache = format
			atomic.StoreInt64(&timeCache.seconds, newValue)
		}
	}
	return timeCache.formatCache, milliSecondsFormat
}

// fresh data in timeCache
func (timeCache *timeFormatCacheType) fresh() {
	timeCache.lock.Lock()
	defer timeCache.lock.Unlock()

	// get current time and update timeCache
	now := time.Now()
	timeCache.now = now
	timeCache.format = []byte(now.Format(PrefixTimeFormat))
	date := now.Format(DateFormat)
	if date != timeCache.date {
		timeCache.dateYesterday = timeCache.date
		timeCache.date = now.Format(DateFormat)
	}
}
