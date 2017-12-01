package gae

import (
	"strings"
	"errors"
	"time"
	"strconv"
)

type KeyOptions_Memcache struct {
	KeyOptions
}

func(mko *KeyOptions_Memcache) GenerateKey() (interface{}, error) {
	var idVal string = mko.StringId("")
	if idVal == "" {
		if mko.NumberId(0)== 0 {
			return nil, errors.New("no identifier present")
		}
		idVal=strconv.FormatInt(mko.NumberId(0), 64)
	}
	return strings.Join([]string{mko.EntityId(""), idVal}, mc_separator), nil
}

func (mko *KeyOptions_Memcache) ExpiresAt(unixTimestamp int64) int64{
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
func (mko *KeyOptions_Memcache) ExpiresIn(inSeconds int64) time.Duration{
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