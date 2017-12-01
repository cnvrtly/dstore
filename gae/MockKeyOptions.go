package gae

import (
	"strings"
	"fmt"
	"math/rand"
	"context"
	"time"
)

const mockMCseparator string="~||~"
type mockKeyOptions struct {
	KeyOptions
}

func(mko *mockKeyOptions) GenerateKey() (interface{}, error){
	var idVal string = mko.StringId("")
	if idVal == "" {
		if mko.NumberId(0)== 0 {
			mko.NumberId(rand.Int63())
		}
		idVal=fmt.Sprint(mko.NumberId(0))
	}
	return strings.Join([]string{mko.EntityId(""), idVal}, mockMCseparator), nil
}

func (mko *mockKeyOptions) ExpiresAt(unixTimestamp int64) int64{
	if unixTimestamp !=0 {
		mko.Value["expiresAt"]= unixTimestamp
		delete(mko.Value, "expiresIn")
	}else if unixTimestamp ==-1 {
		delete(mko.Value,"expiresAt")
		delete(mko.Value, "expiresIn")
	}
	if mko.Value["expiresAt"] != nil {
		return mko.Value["expiresAt"].(int64)
	}
	if mko.Value["expiresIn"] != nil {
		// this is different/recalc on every invocation - should it set  mko.Value["expiresAt"]?
		return time.Now().Unix() + ( mko.Value["expiresIn"]).(int64)
	}
	return 0
}
func (mko *mockKeyOptions) ExpiresIn(inSeconds int64) time.Duration{
	if inSeconds !=0 {
		mko.Value["expiresIn"]= inSeconds
		delete(mko.Value, "expiresAt")
	}else if inSeconds ==-1 {
		delete(mko.Value,"expiresIn")
		delete(mko.Value, "expiresAt")
	}
	if mko.Value["expiresIn"] != nil {
		return time.Duration(mko.Value["expiresIn"].(int64))*time.Second
	}
	if mko.Value["expiresAt"] != nil {
		return time.Unix( mko.Value["expiresAt"].(int64),0).Sub(time.Now())
	}

	return 0
}


func (mko *mockKeyOptions) Ctx(val context.Context) context.Context{
	if val!= nil {
		mko.Value["ctx"]=val
	}
	if mko.Value["ctx"]== nil {
		return nil
	}
	return mko.Value["ctx"].(context.Context)
}