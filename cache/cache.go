package cache

import (
	"github.com/patrickmn/go-cache"
	"time"
)

var Storage  *cache.Cache

func New() *cache.Cache {
	return cache.New(5 * time.Minute, 10 * time.Minute)
}


