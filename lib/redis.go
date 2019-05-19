package lib

import (
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"math/rand"
	"time"
)

func RedisConnFactory(name string) (redis.Conn, error) {
	for confName, cfg := range ConfRedisMap.List {
		if name == confName {
			randHost := cfg.ProxyList[rand.Intn(len(cfg.ProxyList))]
			return redis.Dial(
				"tcp",
				randHost,
				redis.DialConnectTimeout(50*time.Millisecond),
				redis.DialReadTimeout(100*time.Millisecond),
				redis.DialWriteTimeout(100*time.Millisecond))
		}
	}
	return nil, errors.New("create redis conn fail")
}

func RedisLogDo(trace *TraceContext, c redis.Conn, commandName string, args ...interface{}) (interface{}, error) {
	startExecTime := time.Now()
	reply, err := c.Do(commandName, args...)
	endExecTime := time.Now()
	if err != nil {
		Log.TagError(trace, "_com_redis_failure", map[string]interface{}{
			"method":    commandName,
			"err":       err,
			"bind":      args,
			"proc_time": fmt.Sprintf("%fms", endExecTime.Sub(startExecTime).Seconds()),
		})
	} else {
		Log.TagInfo(trace, "_com_redis_success", map[string]interface{}{
			"method":    commandName,
			"bind":      args,
			"proc_time": fmt.Sprintf("%fms", endExecTime.Sub(startExecTime).Seconds()),
		})
	}
	return reply, err
}
