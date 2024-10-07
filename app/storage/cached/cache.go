package cached

import (
	"context"
	"errors"
	"github.com/patrickmn/go-cache"
	redis2 "github.com/redis/go-redis/v9"
	"log"
	"math/rand/v2"
	"reflect"
	"star/app/storage/mysql"
	"star/app/storage/redis"
	"star/constant/str"
	"star/models"
	"sync"
	"time"
)

//只适用于那些读多写少的微小型数据的储存，不能储存大规模数据

// 随机写入redis的时间范围
const redisRandom = 60

var cacheMap = make(map[string]*cache.Cache)
var mu sync.RWMutex

// 创建或获取cache
func getOrCreateCache(name string) *cache.Cache {
	c, ok := cacheMap[name]
	if !ok {
		mu.Lock()
		defer mu.Unlock()
		c, ok = cacheMap[name]
		if !ok {
			c = cache.New(5*time.Minute, 10*time.Minute)
			cacheMap[name] = c
			return c
		}
	}
	return c
}

// ScanGetUser  从缓存中读取数据，如果没有逐步往redis,mysql获取，传入user中
func ScanGetUser(ctx context.Context, key string, obj *models.User) (bool, error) {

	//先在缓存里找
	c := getOrCreateCache(key)
	x, found := c.Get(key)
	if found {
		*obj = *x.(*models.User)
		return true, nil
	}

	//再在redis里找
	if err := redis.Client.HGetAll(ctx, key).Scan(obj); err != nil {
		if !errors.Is(err, redis2.Nil) {
			log.Println("redis hGetAll obj error:", err)
			return false, str.ErrServiceBusy
		}
	}
	//wrappedObj := obj.(cacheItem)
	if obj.IsDirty() {
		c.Set(key, reflect.ValueOf(obj).Elem(), cache.DefaultExpiration)
		return true, nil
	}

	//再往mysql里找
	err := mysql.QueryUserInfo(obj, obj.GetID())
	if err != nil {
		return false, err
	}

	//将查询的值储存到redis和缓存里
	if err := redis.Client.HSet(ctx, key, obj).Err(); err != nil {
		log.Println("set obj error", err)
		return true, str.ErrServiceBusy
	}
	c.Set(key, obj, cache.DefaultExpiration)

	return true, nil
}

// ScanDeleteUser 将user从缓存里删除，下次读取回到mysql里读取
func ScanDeleteUser(ctx context.Context, key string) {

	c := getOrCreateCache(key)
	c.Delete(key)
	redis.Client.HDel(ctx, key)

}

// Get 查询字符串，从缓存中读取数据，读取成功返回true,失败返回false,不存在也返回false
func Get(ctx context.Context, key string) (string, bool, error) {
	c := getOrCreateCache("strings")
	//先从cache里查询
	if obj, found := c.Get(key); found {
		return obj.(string), true, nil
	}
	//如果没查到往redis查
	result := redis.Client.Get(ctx, key)
	if result.Err() != nil && !errors.Is(result.Err(), redis2.Nil) {
		log.Println("get cache error:", result.Err())
		return "", false, str.ErrServiceBusy
	}
	value, err := result.Result()
	switch {
	case errors.Is(err, redis2.Nil):
		return "", false, nil
	default:
		c.Set(key, value, cache.DefaultExpiration)
		return value, true, nil
	}
}

// GetWithFunc 获取字符串缓存，如果不存在，则调用传入函数
func GetWithFunc(ctx context.Context, key string, f func(string) (string, error)) (string, error) {
	value, ok, err := Get(ctx, key)
	if err != nil {
		return "", err
	}
	if ok {
		return value, nil
	}
	//如果不存在，则调用传入函数获取
	value, err = f(key)
	if err != nil {
		return "", err
	}
	Write(ctx, key, value, true)
	return value, nil
}

// Write 写入字符串缓存，如果state为true，则写入redis缓存
func Write(ctx context.Context, key string, value string, state bool) {
	c := getOrCreateCache("strings")
	c.Set(key, value, cache.DefaultExpiration)
	if state {
		redis.Client.Set(ctx, key, value, 72*time.Hour+time.Duration(rand.IntN(redisRandom))*time.Minute)
	}
}

// Delete 删除字符串缓存
func Delete(ctx context.Context, key string) {
	c := getOrCreateCache("strings")
	c.Delete(key)
	redis.Client.Del(ctx, key)
}

// WriteWithOvertime 根据overtime写入redis缓存
func WriteWithOvertime(ctx context.Context, key string, value string, overtime time.Duration) {
	c := getOrCreateCache("strings")
	c.Set(key, value, cache.DefaultExpiration)
	redis.Client.Set(ctx, key, value, overtime)
}
