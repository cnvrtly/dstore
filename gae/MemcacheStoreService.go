package gae

import (
	"fmt"
	"github.com/cnvrtly/dstore"
	"context"
	"google.golang.org/appengine/memcache"
)

type MemcacheStoreService struct {
}

func (memSer *MemcacheStoreService) CreateKeyOptions(ctx context.Context, namespaceId string, entityId string, stringId string, numId int64) (dstore.KeyOptions,error) {

	if l:=len(stringId);l> 250 {
		return nil, dstore.TooLongKeyError{Msg:"store:key-too-long", DeltaVal:l-250}
	}
	ko:=&KeyOptions_Memcache{}
	ko.Value= map[string]interface{}{}
	ko.Ctx(ctx)
	ko.EntityId(entityId)
	ko.StringId(stringId)
	ko.NumberId(numId)
	ko.NamespaceId(namespaceId)
	return ko,nil
}

func (memSer *MemcacheStoreService) Save(keyOptions dstore.KeyOptions, value interface{}, options interface{}) (interface{}, error) {
	//fmt.Println("CACHE save")
	key, err:=keyOptions.GenerateKey()
	if err!=nil {
		return nil, err
	}
	if len(key.(string))> 250 {
		return nil, fmt.Errorf("memcache key must be len<250 len=%v key=%v", len(key.(string)), key.(string))
	}
	itm:=&memcache.Item{Key:key.(string), Value:[]byte(fmt.Sprint(value))}
	if exp:=keyOptions.ExpiresIn(0); exp!=0 {
		itm.Expiration=exp
	}

	err=memcache.Set(keyOptions.Ctx(nil), itm)
	if err!= nil {
		return nil, fmt.Errorf("memcacheStore.Save err=%v value=%v expires=%v keyLen=%v key=%v", err,value, keyOptions.ExpiresIn(0),len(key.(string)), key)
	}
	return value, nil
}

func (memSer *MemcacheStoreService) Delete(keyOptions dstore.KeyOptions) (error) {
	key, err:=keyOptions.GenerateKey()
	if err!= nil {
		return err
	}
	err=memcache.Delete(keyOptions.Ctx(nil), key.(string))
	if err==memcache.ErrCacheMiss {
		return dstore.ErrorNotFound
	}
	if err!= nil {
		return err
	}
	return nil
}

func (memSer *MemcacheStoreService) Load(keyOptions dstore.KeyOptions, setValueOnPointer interface{}) (interface{}, error) {
	//fmt.Println("mem CACHE get")
	key, err:=keyOptions.GenerateKey()
	if err!= nil {
		return nil, err
	}
	itm_p, err:=memcache.Get(keyOptions.Ctx(nil), key.(string))
	if err==memcache.ErrCacheMiss {
		return nil, dstore.ErrorNotFound
	}
	if err!= nil {
		return nil, err
	}
	res := string(itm_p.Value)
	switch setValueOnPointer.(type) {
	case *string:
		*setValueOnPointer.(*string)=res
	}
	return res, nil
}

func (memSer *MemcacheStoreService) GetAll(ctx context.Context, namespaceId string, entityId string, fillSlice interface{}, queryModifiers func(interface{})( interface{})) (interface{}, error) {
	return nil, dstore.ErrorNotImplemented
}

func (memSer *MemcacheStoreService) QueryIterator(ctx context.Context, namespaceId string, entityId string, queryModifiers func(interface{})( interface{})) (interface{}, error) {
	return nil, dstore.ErrorNotImplemented
}
