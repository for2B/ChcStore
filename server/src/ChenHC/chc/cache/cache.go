package cache

import (
	"log"

	"github.com/asdine/storm"
	"fmt"
)

type CacheManager struct {
	stormDB     *storm.DB
	syncEventID chan int
}

func NewCacheManager(stormDBpath string) *CacheManager {
	cm := &CacheManager{}
	// fmt.Println(stormDBpath)
	cm.InitStormDB(stormDBpath)
	fmt.Println("init CacheManager success")
	return cm
}

func (cm *CacheManager) InitStormDB(stormDBpath string) {
	stormDB, err := storm.Open(stormDBpath)
	if err != nil {
		log.Panicln("storm.Open fail:", err)
	}
	cm.stormDB = stormDB
	cm.syncEventID = make(chan int, 1000)

}

func (cm *CacheManager) Set(buckeyName string, key interface{}, value interface{}) error {
	return cm.stormDB.Set(buckeyName, key, value)
}

func (cm *CacheManager) Get(buckeyName string, key interface{}, to interface{}) error {
	return cm.stormDB.Get(buckeyName, key, to)
}

func (cm *CacheManager) Del(buckeyName string, key interface{}) error {
	return cm.stormDB.Delete(buckeyName, key)
}
