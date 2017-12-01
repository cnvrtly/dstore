package gae

import (
	"context"
)

const mc_separator string="~§~"
type KeyOptions struct {
	Value map[string]interface{}
}
func (mko *KeyOptions) NamespaceId(val string) string{
	if val!= "" {
		mko.Value["namespaceId"]=val
		if ctx:=mko.Ctx(nil); ctx!=nil {
			mko.Ctx(ctx)
		}
	}
	//TODO replace with: if val, ok := mko.Value["namespaceId"]; ok && val!=""
	if mko.Value["namespaceId"] != nil {
		return mko.Value["namespaceId"].(string)
	}
	return ""
}
func (mko *KeyOptions) EntityId(val string) string{
	if val!= "" {
		mko.Value["entityNameId"]=val
	}
	if mko.Value["entityNameId"] != nil {
		return mko.Value["entityNameId"].(string)
	}
	return ""
}
func (mko *KeyOptions) StringId(val string) string{
	if val!= "" {
		mko.Value["stringId"]=val
	}
	if mko.Value["stringId"] != nil {
		return mko.Value["stringId"].(string)
	}
	return ""
}
func (mko *KeyOptions) NumberId(val int64) int64{
	if val!=0 {
		mko.Value["numberId"]=val
	}
	if mko.Value["numberId"] != nil {
		return mko.Value["numberId"].(int64)
	}
	return 0
}
func (mko *KeyOptions) Ctx(val context.Context) context.Context{
	if val!= nil {
		mko.Value["ctx"]=val
		if ns:=mko.NamespaceId(""); ns!= "" {
			nsCtx, err:=WithNamespaceId(val, ns)
			if err!= nil {
				panic(err)
			}
			mko.Value["nsCtx"]=nsCtx
		}
	}
	if mko.Value["ctx"]== nil && mko.Value["nsCtx"]== nil {
		return nil
	}
	if mko.Value["nsCtx"]!= nil {
		return mko.Value["nsCtx"].(context.Context)
	}
	return mko.Value["ctx"].(context.Context)
}
/* should be implemented for specific store type
func(mko *KeyOptions) GenerateKey() (interface{}, error) {
	var idVal string = mko.stringId
	if idVal == "" {
		if mko.numberId== 0 {
			mko.numberId=rand.Int63()
		}
		idVal=fmt.Sprint(mko.numberId)
	}
	return strings.Join([]string{mko.entityNameId, idVal}, mc_separator)
}*/


