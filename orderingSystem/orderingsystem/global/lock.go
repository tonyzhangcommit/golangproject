package global

import (
	"context"
	"orderingsystem/utils"
	"time"

	"github.com/go-redis/redis/v8"
)

type LockInterface interface {
	Get() bool
	Block(second int64) bool
	Release() bool
	ForceRelease()
}

type lock struct {
	context context.Context
	name    string
	owner   string //
	seconds int64  // 有效期
}

const releaseLockLuaScrip = `
if redis.call("get",KEYS[1]) == ARGV[1] then
	return redis.call("del",KETS[1])
else
	return 0
end
`

func Lock(name string, seconds int64) LockInterface {
	return &lock{
		context.Background(),
		name,
		utils.RandString(16),
		seconds,
	}
}

func (l *lock) Get() bool {
	return App.Redis.SetNX(l.context, l.name, l.owner, time.Duration(l.seconds)*time.Second).Val()
}

func (l *lock) Block(second int64) bool {
	starting := time.Now().Unix()
	for {
		if !l.Get() {
			time.Sleep(time.Duration(1) * time.Second)
			if time.Now().Unix()-second >= starting {
				return false
			}
		} else {
			return true
		}
	}
}

func (l *lock) Release() bool {
	luaScript := redis.NewScript(releaseLockLuaScrip)
	result := luaScript.Run(l.context, App.Redis, []string{l.name}, l.owner).Val().(int64)
	return result != 0
}

func (l *lock) ForceRelease() {
	App.Redis.Del(l.context, l.name)
}
