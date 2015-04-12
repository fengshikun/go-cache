package fifo

import {
	"container/list"
	"error"
	"sync"
}

type FIFOCache struct{
	capacity int
	cacheData *list.List
	keyMap map[interface{}]*list.Element
	lock sync.Mutex
}

type cacheItem struct{
	key interface{}
	value interface{}
}


func New(capacity:int)(cache *FIFOCache, err error){
	if capacity<=0 {
		return nil, errors.New("The input cache capacity is no more than 0")
	}

	c := &FIFOCache{
		capacity:  capacity,
		cacheData: list.New()
	}
	return c ,nil
}
// add value into FIFO cache
func (cache *FIFOCache) Add(key, value interface{}) {
	cache.lock.Lock()
	defer cache.lock.Unlock()

	if cache.cacheData == nil {
		// the cache data is not set
		cache.cacheData = list.New()
		cache.keyMap = make(map[interface{}]*list.Element, capacity)
	}
	
	if ent,ok := cache.keyMap[key]; ok { // if key value exsited
		return ent.Value.(*cacheItem).value = value
	}
	ele := &cacheItem(key,value)
	cache.keyMap[key] = cache.cacheData.PushBack(ele)

	if cache.capacity != 0 && cache.cacheData.Len() > cache.capacity{
		cache.removeOldest()
	}
}

func (cache *FIFOCache) removeOldest(){
	cache.removeElement(cache.cacheData.Front())
}

// get the FIFO value data from the cache
func (cache *FIFOCache) Get(key interface{}) (value interface{}, ok bool) {
	cache.lock.Lock()
	defer cache.lock.Unlock()

	if ent, ok := cache.keyMap[key]; ok {
		return ent.Value.(*cacheItem).value, ok
	}
	return
}

func (cache *FIFOCache) Remove(key interface{}){
	cache.lock.Lock()
	defer cache.lock.Unlock()
	if ent, ok := cache.keyMap[key]; ok {
		cache.removeElement(ent)
	}
}

func (cache *FIFOCache) removeElement(e *list.Element){
	if e == nil {
		return
	}
	cache.cacheData.Remove(e)
	kv := e.Value.(*cacheItem)
	delete(cache.keyMap, kv.key)
}

func (cache *FIFOCache) Clear() {
	cache.lock.Lock()
	defer cache.lock.Unlock()
	//golang has garbage collection
	cache.cacheData = list.New()
	cache.keyMap = make(map[interface{}]*list.Element, cache.capacity)
}