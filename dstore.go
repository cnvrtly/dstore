package dstore

import (
	"context"
	"errors"
	"time"
)

//TODO create type for ctx value ident
const GlobalStoreNamespaceId string = "_globalAppNS"
const CtxKeyNamespaceId ctxNamespaceIdType = ctxNamespaceIdType("store:namespaceId")

type ctxNamespaceIdType string

var ErrorNotFound error = errors.New("store:not-found-error")
var ErrorKeyTooLong error = errors.New("store:key-too-long")
var ErrorNotImplemented error = errors.New("store:method-not-implemented")

func NamespaceIdFromContext(ctx context.Context) (string, error) {
	if nsId := ctx.Value(CtxKeyNamespaceId); nsId != nil && nsId.(string) != "" {
		return nsId.(string), nil
	}
	return "", errors.New("namespaceId not defined on context")
}

func WithNamespaceId(ctx context.Context, namespaceId string) context.Context {
	return context.WithValue(ctx, CtxKeyNamespaceId, namespaceId)
}

type KeyOptions interface {
	NamespaceId(val string) string
	EntityId(val string) string
	StringId(val string) string
	NumberId(val int64) int64
	ExpiresAt(val int64) int64
	ExpiresIn(inSeconds int64) time.Duration
	Ctx(val context.Context) context.Context
	GenerateKey() (interface{}, error)
}

type SaverRetriever interface {
	Save(keyOptions KeyOptions, value interface{}, options interface{}) (interface{}, error)
	Load(keyOptions KeyOptions, setValueOnPointer interface{}) (interface{}, error)
	Delete(keyOptions KeyOptions) (error)
	GetAll(ctx context.Context, namespaceId string, entityId string, fillSlice interface{}, queryModifiers func(interface{})(interface{}) ) (interface{}, error)
	QueryIterator(ctx context.Context, namespaceId string, entityId string, queryModifiers func(interface{})(interface{}) ) (interface{}, error)
	CreateKeyOptions(ctx context.Context, namespaceId string, entityId string, stringId string, numberId int64) (KeyOptions, error)
}

type Findable interface {
	FindBy(val string) (string, error)
}

type FindableEnt struct {
	ident string
}

func (se *FindableEnt) String() (string) {
	return se.ident
}
func (se *FindableEnt) FindBy(idValue string) (string, error) {
	if se == nil {
		return "", errors.New("FindableEnt is nul - property was not set when creating instance")
	}
	if idValue != "" {
		se.ident = idValue
	}
	return se.ident, nil
}

type TooLongKeyError struct {
	Msg      string
	DeltaVal int
}

func (tlk TooLongKeyError) Delta() int {
	return tlk.DeltaVal
}

func (tlk TooLongKeyError) Error() string {
	return tlk.Msg
}

/////////////
