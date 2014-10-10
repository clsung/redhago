package redhago

import (
	"time"

	"github.com/golang/glog"
)

func RedisSetWaitGet(pool *Pool, key string, value interface{}, wait, expire int) ([]byte, error) {
	var err error
	if expire == 0 {
		err = pool.Set(key, value)
	} else {
		err = pool.Setex(key, value, expire)
	}
	if err != nil {
		glog.Errorf("Set key %s with %v failed: %v", key, value, err)
		return nil, err
	}
	time.Sleep(time.Duration(wait) * time.Second)

	return pool.Get(key)
}
