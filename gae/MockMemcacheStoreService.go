package gae

import (
	"fmt"
	"sync"
	"github.com/cnvrtly/dstore"
	"context"
	"errors"
)

var cacheValues map[string]string = make(map[string]string)
var mtx sync.Mutex

type MockMemcacheStoreService struct {
}

func (memSer *MockMemcacheStoreService) CreateKeyOptions(ctx context.Context, namespaceId string, entityId string, stringId string, numId int64) ( dstore.KeyOptions, error) {
	if l:=len(stringId);l> 250 {
		return nil, dstore.TooLongKeyError{Msg:"store:key-too-long", DeltaVal:l-250}
	}
	ko := &mockKeyOptions{}
	ko.Value = map[string]interface{}{}
	ko.Ctx(ctx)
	ko.EntityId(entityId)
	ko.StringId(stringId)
	ko.NumberId(numId)
	ko.NamespaceId(namespaceId)
	return ko, nil
}
func (memSer *MockMemcacheStoreService) Save(keyOptions dstore.KeyOptions, value interface{}, options interface{}) (interface{}, error) {
	fmt.Println("CACHE save")
	key,_:=keyOptions.GenerateKey()
	mtx.Lock()
	cacheValues[key.(string)] = value.(string)
	mtx.Unlock()
	return value, nil
}

func (memSer *MockMemcacheStoreService) Delete(keyOptions dstore.KeyOptions) (error) {
	return errors.New("not yet implemented")
}
func (memSer *MockMemcacheStoreService) Load(keyOptions dstore.KeyOptions, setValueOnPointer interface{}) (interface{}, error) {
	fmt.Println("CACHE get")
	key,_:=keyOptions.GenerateKey()
	mtx.Lock()
	res := cacheValues[key.(string)]
	mtx.Unlock()
	if res== "" {
		return nil, dstore.ErrorNotFound
	}
	switch setValueOnPointer.(type) {
	case *string:
		*setValueOnPointer.(*string)=res
	}
	return res, nil
}
