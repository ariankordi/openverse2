package config

import (
  //"github.com/go-redis/redis"
  "github.com/patrickmn/go-cache"

  "time"
)

/*func initCache() *redis.Client {
  client := redis.NewClient(&redis.Options{
    Addr: ":6379",
    DB: 2,
  })
  client.FlushDB()
  return client
}*/

//var CacheStore = initCache()
var CacheStore = cache.New(30 * time.Minute, 60 * time.Minute)
