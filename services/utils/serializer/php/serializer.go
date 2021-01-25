package php

import (
	"github.com/luoxiaojun1992/go-php-serialize/phpserialize"
	"github.com/luoxiaojun1992/go-skeleton/services/helper"
	"sync"
)

var cacheLock sync.RWMutex
var decodeCache map[string]interface{}

func Setup() {
	decodeCache = make(map[string]interface{})
}

func Decode(value string) (result interface{}, err error) {
	cacheLock.RLock()
	cache, hasCache := decodeCache[value]
	if hasCache {
		cacheLock.RUnlock()
		return cache, nil
	}
	cacheLock.RUnlock()

	res, resErr := phpserialize.Decode(value)
	if !helper.CheckErr(resErr) {
		cacheLock.Lock()
		if len(decodeCache) < 20 {
			decodeCache[value] = res
		}
		cacheLock.Unlock()
	}

	return res, resErr
}
