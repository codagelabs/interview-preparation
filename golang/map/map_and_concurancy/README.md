# Map and Concurancy
It’s important to note that maps in go are not safe for concurrent use. 

By default, maps are not thread-safe. This means that if you try to read or write to a map from multiple goroutines simultaneously, you may encounter race conditions and unpredictable behavior.

To safely use maps concurrently, you need to implement proper synchronization mechanisms. Here are some common approaches:

1. Use a Mutex or RWMutex to lock the map when accessing it.

```go
var mu sync.Mutex

func safeMapAccess(key string) {
    mu.Lock()
    defer mu.Unlock()
}   
```

2. Use a sync.Map, which is designed for concurrent use.

```go   
var m sync.Map

func concurrentMapAccess(key string) {
    m.Load(key)
    m.Store(key, value)
}
``` 

3. Use a channel to synchronize access to the map.

```go
ch := make(chan map[string]int)

func safeMapAccess(key string) {    
    ch <- map[string]int{key: value}
}                                   
``` 
