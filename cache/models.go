package cache

import (
	"github.com/nd-r/tech-db-forum/models"
	"sync"
	"github.com/cornelk/hashmap"
	"unsafe"
)

type CachedForum struct {
	Model models.Forum
	Json  []byte
}

type CachedUser struct {
	Model models.User
	Json  []byte
}

type ForumCache struct {
	cMap hashmap.HashMap
	lock sync.RWMutex
}

type UserCache struct {
	cMap hashmap.HashMap
	lock sync.RWMutex
}

func (fc *ForumCache) Get(slug interface{}) (*CachedForum) {
	if val, ok := fc.cMap.Get(slug); ok {
		return (*CachedForum)(val)
	}

	return nil
}

func (fc *ForumCache) Push(slug interface{}, cf *CachedForum) {
	fc.cMap.Set(slug, unsafe.Pointer(cf))
}

func (fc *ForumCache) Clear() {
	fc.cMap = hashmap.HashMap{}
}

func (uc *UserCache) Get(nickname interface{}) ([]byte) {

	if val, ok := uc.cMap.Get(nickname); ok {
		return *(*[]byte)(val)
	}
	return nil
}

func (uc *UserCache) Push(nickname interface{}, cu *[]byte) {
	uc.cMap.Set(nickname, unsafe.Pointer(cu))
}

func (uc *UserCache) Clear() {
	uc.cMap = hashmap.HashMap{}
}

var TheGreatForumCache = ForumCache{cMap: hashmap.HashMap{}}
var TheGreatUserCache = UserCache{cMap: hashmap.HashMap{}}
