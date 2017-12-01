package gae

import (
	"testing"
	"context"
	"github.com/cnvrtly/dstore"
	"fmt"
	"encoding/json"
)


func TestMockMemcacheStoreService_Save(t *testing.T) {
	ctx:=context.Background()
	ctx= dstore.WithNamespaceId(ctx, "test.com")
	memc:=&MockMemcacheStoreService{}
	val,err:=memc.CreateKeyOptions(ctx,"test.com","MOCK_1","key1",0)
	if err != nil {
		t.Error( err)
	}

	keyOpt:= val
	res, _:=memc.Load(keyOpt,nil)

	memc.Save(keyOpt, "value1", nil)
	res, err :=memc.Load(keyOpt, nil)
	if err!= nil {
		t.Error(err)
	}
	if res.(string)!= "value1" {
		t.Error("should be value1")
	}

	keyOpt1,err:= memc.CreateKeyOptions(ctx,"test.com","MOCK_2", "key1", 0)
	if err != nil {
		t.Error(err)
	}

	res1, notFoundErr:=memc.Load(keyOpt1, nil)

	if res1!=nil {
		t.Error("should be empty value/string=", res1)
	}

	if notFoundErr!= dstore.ErrorNotFound {
		t.Error("should be ErrorNotFound error")
	}

	keyOpt2,err:= memc.CreateKeyOptions(ctx,"test.com","MOCK_1","key1", 0)
	if err != nil {
		t.Error(err)
	}
	res2, _:=memc.Load(keyOpt2, nil)

	if res2== "" {
		t.Error("should have value")
	}
}

func TestMockMemcacheStoreService_Load(t *testing.T) {
	ctx:=context.Background()
	ctx= dstore.WithNamespaceId(ctx, "test.com")
	memSer:=&MockMemcacheStoreService{}
	keyOpt,err:=memSer.CreateKeyOptions(ctx,"test.com", "MOCK2", "id1",0)
	if err != nil {
		t.Error(err)
	}

	memSer.Save(keyOpt,"123", nil)
	saved:="val555"
	memSer.Load(keyOpt, &saved)
	if saved!= "123" {
		t.Error("pointer string not changed on Load value")
	}else{
		fmt.Println("Pointer CNGED OK")
	}

}

type testData struct {

	TitleNoFormatting string
	URL               string
}

func TestMockDatastoreStoreService_Load(t *testing.T) {
	ctx:=context.Background()
	ctx= dstore.WithNamespaceId(ctx, "test.com")
	dsSer :=&MockDatastoreStoreService{}
	keyOpt,err:= dsSer.CreateKeyOptions(ctx,"test.com", "MOCK2", "",0)
	if err != nil {
		t.Error(err)
	}

	urlVal := "UUUU"
	val:=testData{TitleNoFormatting: "yes", URL: urlVal}

	dsSer.Save(keyOpt,val , nil)
	saved:=&testData{}

	respVal, err:=dsSer.Load(keyOpt, saved)
	if err!= nil {
		t.Error(err)
	}
	//fmt.Println("RRRR=", respVal, "SSS=", saved.URL)
	if saved.URL!=urlVal {
		t.Error("pointer not set on Load value", saved)
	}
	resTestData:=testData{}
	//TODO create from json and compare .URL to urlVal
	resDataErr:=json.Unmarshal([]byte(respVal.(string)),&resTestData)
	if resDataErr!=nil {
		t.Error(resDataErr)
	}
	if resTestData.URL!=urlVal {
		t.Error("pointer not returned on Load value", respVal)
	}

}
