package gae

import (
	"context"
)

type KeyOptions_Gae struct {
	KeyOptions
}

//TODO make gae specific implementation for appengine context dependency
func (mko *KeyOptions_Gae) Ctx(val context.Context) context.Context{
	if val!= nil {
		mko.KeyOptions.Value["ctx"]=val
	}
	if mko.KeyOptions.Value["ctx"]== nil {
		return nil
	}
	return mko.KeyOptions.Value["ctx"].(context.Context)
}