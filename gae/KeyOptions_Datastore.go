package gae

import (
	"google.golang.org/appengine/datastore"
	"time"
)

type KeyOptions_Datastore struct {
	KeyOptions
	generatedKey *datastore.Key
}

func (koDs *KeyOptions_Datastore) GenerateKey() (interface{}, error) {
	if koDs.generatedKey == nil {
		koDs.generatedKey = datastore.NewKey(koDs.Ctx(nil), koDs.EntityId(""), koDs.StringId(""), koDs.NumberId(0), nil)
	}
	return koDs.generatedKey, nil
}

/* can not override - create separate key options*/
func (koDs *KeyOptions_Datastore) ExpiresAt(unixTimestamp int64) int64 {
	return 0
}
func (koDs *KeyOptions_Datastore) ExpiresIn(inSeconds int64) time.Duration {
	return time.Second * 0
}
