# 本地缓存

#初始化

```go
var (
	localCache *cache.LocalCache
)

//初始化，5分钟为缓存有效时间，每隔30秒清空过期缓存
func init() {
	localCache = cache.New(5*time.Minute, 30*time.Second)
}

func GetLocalCache() *cache.LocalCache {
	return localCache
}
```

