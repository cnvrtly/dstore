package gae

import (
	"fmt"
	"sync"
	"github.com/cnvrtly/dstore"
	"context"
	"encoding/json"
	"errors"
)

var dsValues map[string]string = make(map[string]string)
var dsMtx sync.Mutex

type MockDatastoreStoreService struct {
}

func (memSer *MockDatastoreStoreService) CreateKeyOptions(ctx context.Context, namespaceId string, entityId string, stringId string, numId int64) (dstore.KeyOptions,error) {

	if l:=len(stringId);l> 500 {
		return nil, dstore.TooLongKeyError{Msg:"store:key-too-long", DeltaVal:l-500}//dstore.ErrorKeyTooLong
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
func (memSer *MockDatastoreStoreService) Save(keyOptions dstore.KeyOptions, value interface{}, options interface{}) (interface{}, error) {
	fmt.Println("DATASTORE save")
	if keyOptions.StringId("")=="" && keyOptions.NumberId(0)== 0 {
		keyOptions.GenerateKey()
		//memSer.setId(keyOptions.NumberId(0), value)
	}
	valueJson, err:=json.Marshal(value)
	if err!= nil {
		return nil, err
	}
	key,_:=keyOptions.GenerateKey()
	dsMtx.Lock()
	dsValues[key.(string)] = string(valueJson)
	dsMtx.Unlock()

	err=setStorableEntIdent(value, keyOptions)
	if err!= nil {
		return nil, err
	}

	return value, nil
}

func (memSer *MockDatastoreStoreService) Delete(keyOptions KeyOptions) (error) {
	return errors.New("not yet implemented")
}
func (memSer *MockDatastoreStoreService) Load(keyOptions dstore.KeyOptions, setValueOnPointer interface{}) (interface{}, error) {
	fmt.Println("DATASTORE get")
	key,_:=keyOptions.GenerateKey()
	dsMtx.Lock()
	res := dsValues[key.(string)]
	dsMtx.Unlock()
	if res== "" {
		return nil, dstore.ErrorNotFound
	}
	if setValueOnPointer != nil && res!=""{
		json.Unmarshal([]byte(res), setValueOnPointer)
	}

	err:=setStorableEntIdent(setValueOnPointer, keyOptions)
	if err!= nil {
		return nil, err
	}

	return res, nil
}
/*

func (memSer *MockDatastoreStoreService) setId(idVal int64, setValueOnPointer interface{}) (interface{}, error) {

	setValueOnPointer.MockId=idVal
	*/
/*elem := reflect.TypeOf(setValueOnPointer).Elem()
	fieldLen:=elem.NumField()
	for i:=0 ; i<fieldLen; i++{
		field:=elem.FieldByIndex([]int{i})
		if tag := string(field.Tag); tag!= "" {

		}
	}*//*



}*/
