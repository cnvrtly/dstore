package gae

import (
	"context"
	"google.golang.org/appengine/datastore"
	"github.com/cnvrtly/dstore"
	"fmt"
	"encoding/json"
)

type DatastoreStoreService struct {
}

func (dsSer *DatastoreStoreService) CreateKeyOptions(ctx context.Context, namespaceId string, entityId string, stringId string, numId int64) (dstore.KeyOptions, error) {

	if l:=len(stringId);l> 500 {
		return nil, dstore.TooLongKeyError{Msg:"store:key-too-long", DeltaVal:l-500}//dstore.ErrorKeyTooLong
	}
	ko := &KeyOptions_Datastore{}
	ko.Value = map[string]interface{}{}
	ko.Ctx(ctx)
	ko.EntityId(entityId)
	ko.StringId(stringId)
	ko.NumberId(numId)
	ko.NamespaceId(namespaceId)
	return ko, nil
}

func (dsSer *DatastoreStoreService) Save(keyOptions dstore.KeyOptions, saveValueFromPointer interface{}, options interface{}) (interface{}, error) {

	key, err := keyOptions.GenerateKey()
	if err != nil {
		return nil, err
	}

	retKey, err := datastore.Put(keyOptions.Ctx(nil), key.(*datastore.Key), saveValueFromPointer)
	if err != nil {
		strVal,_:=json.Marshal(saveValueFromPointer)
		fmt.Printf("DS save err =%s ptr=%s\n",err, string(strVal))
		return nil, err
	}
	keyOptions.NumberId(retKey.IntID())
	keyOptions.StringId(retKey.StringID())
	err=setStorableEntIdent(saveValueFromPointer, keyOptions)
	if err!= nil {
		return nil, err
	}

	return retKey, nil
}

func (dsSer *DatastoreStoreService) Delete(keyOptions dstore.KeyOptions) (error) {
	key, err := keyOptions.GenerateKey()
	if err != nil {
		return err
	}
	return datastore.Delete(keyOptions.Ctx(nil), key.(*datastore.Key))
}

func (dsSer *DatastoreStoreService) Load(keyOptions dstore.KeyOptions, setValueOnPointer interface{}) (interface{}, error) {
	//fmt.Println("DS get")
	key, err := keyOptions.GenerateKey()
	if err != nil {
		return nil, err
	}
	err = datastore.Get(keyOptions.Ctx(nil), key.(*datastore.Key), setValueOnPointer)
	if err == datastore.ErrNoSuchEntity {
		return nil, dstore.ErrorNotFound
	}
	if err != nil {
		return nil, err
	}

	err=setStorableEntIdent(setValueOnPointer, keyOptions)
	if err!= nil {
		return nil, err
	}

	return setValueOnPointer, nil
}
