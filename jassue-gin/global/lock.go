package global

import (
	"context"
	"jassue-gin/utils"
	"time"

	"github.com/go-redis/redis/v8"
)

type Interface interface {
	Get() bool
	Block(second int64) bool
	Release() bool
	ForceRelease()
}

type lock struct {
	context context.Context
	name    string
	owner   string
	seconds int64
}

const releaseLockScript = `
if redis.call("get",KEYS(1)) == ARGV(1) then
	return redis.call("del",KEYS(1))
else
	return 0
end
`

func Lock(name string, seconds int64) Interface {
	return &lock{
		context.Background(),
		name,
		utils.RandString(16),
		seconds,
	}
}

func (l *lock) Get() bool {
	// 如果不存在创建，返回true,如果存在，则不进行操作，返回false
	return App.Redis.SetNX(l.context, l.name, l.owner, time.Duration(l.seconds)*time.Second).Val()
}

// 阻塞一段时间，尝试获取锁
func (l *lock) Block(seconds int64) bool {
	starting := time.Now().Unix()
	for {
		if !l.Get() {
			time.Sleep(time.Duration(1) * time.Second)
			if time.Now().Unix()-seconds >= starting {
				return false
			}
		} else {
			return true
		}
	}
}

func (l *lock) Release() bool {
	luaScript := redis.NewScript(releaseLockScript)
	result := luaScript.Run(l.context, App.Redis, []string{l.name}, l.owner).Val().(int64)
	return result != 0
}

func (l *lock) ForceRelease() {
	App.Redis.Del(l.context, l.name).Val()
}
