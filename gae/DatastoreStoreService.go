package gae

import (
	"context"
	"google.golang.org/appengine/datastore"
	"github.com/cnvrtly/dstore"
	"fmt"
	"encoding/json"
	"google.golang.org/appengine"
	"reflect"
)

type DatastoreStoreService struct {
}

func (dsSer *DatastoreStoreService) CreateKeyOptions(ctx context.Context, namespaceId string, entityId string, stringId string, numId int64) (dstore.KeyOptions, error) {

	if l := len(stringId); l > 500 {
		return nil, dstore.TooLongKeyError{Msg: "store:key-too-long", DeltaVal: l - 500} //dstore.ErrorKeyTooLong
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
		strVal, _ := json.Marshal(saveValueFromPointer)
		fmt.Printf("DS save err =%s ptr=%s\n", err, string(strVal))
		return nil, err
	}
	keyOptions.NumberId(retKey.IntID())
	keyOptions.StringId(retKey.StringID())
	err = setStorableEntIdent(saveValueFromPointer, keyOptions)
	if err != nil {
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

	err = setStorableEntIdent(setValueOnPointer, keyOptions)
	if err != nil {
		return nil, err
	}

	return setValueOnPointer, nil
}

func (dsSer *DatastoreStoreService) GetAll(ctx context.Context, namespaceId string, entityId string, fillSlice interface{}, queryModifiers func(interface{}) (interface{})) (interface{}, error) {
	q, ctx, err := createDatastoreQuery(ctx, namespaceId, entityId, queryModifiers)
	if err != nil {
		return nil, err
	}

	keys, err := q.GetAll(ctx, fillSlice)
	if err != nil {
		return nil, err
	}

	if fillSlice != nil && len(keys) > 0 {

		switch reflect.TypeOf(fillSlice).Kind() {
		case reflect.Ptr:
			// fillSlice is pointer to slice of pointers
			val := reflect.ValueOf(fillSlice)
			// get slice pointer
			ptrSlice := val.Elem().Interface()

			switch reflect.TypeOf(ptrSlice).Kind() {
			case reflect.Slice:
				// get the slice
				s := reflect.ValueOf(ptrSlice)
				for i := 0; i < s.Len(); i++ {
					// if element is Findable set FindBy / ident value
					if currFindable, ok := s.Index(i).Interface().(dstore.Findable); ok {
						currFindable.FindBy(keys[i].StringID())
					}
					/*this is same but with member reflection
					findable := s.Index(i).Elem().FieldByName("FindableEnt")
					if findable.Kind() != reflect.Invalid && findable.MethodByName("FindBy").Kind()==reflect.Func{
						//type has "FindBy" / implements findable interface
						ident := keys[i].StringID()
						findable.MethodByName("FindBy").Call([]reflect.Value{reflect.ValueOf(ident)})
					}*/
				}
			}
		}
	}
	return keys, nil
}

func (dsSer *DatastoreStoreService) QueryResults(ctx context.Context, namespaceId string, entityId string, queryModifiers func(interface{}) (interface{})) (interface{}, error) {
	q, ctx, err := createDatastoreQuery(ctx, namespaceId, entityId, queryModifiers)
	if err != nil {
		return nil, err
	}

	return q.Run(ctx), nil
}

func createDatastoreQuery(ctx context.Context, namespaceId string, entityId string, queryModifiers func(interface{}) (interface{})) (*datastore.Query, context.Context, error) {
	ctx, err := appengine.Namespace(ctx, namespaceId)
	if err != nil {
		return nil, nil, err
	}
	q := datastore.NewQuery(entityId)
	if queryModifiers != nil {
		q = queryModifiers(q).(*datastore.Query)
	}

	return q, ctx, nil
}
