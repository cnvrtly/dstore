package gae

import (
	"context"
	"google.golang.org/appengine"
	"github.com/cnvrtly/dstore"
	"fmt"
)

func WithNamespaceId(ctx context.Context, namespaceId string) (context.Context, error) {
	ctx, err := appengine.Namespace(ctx, namespaceId)
	if err != nil {
		return nil, err
	}
	return dstore.WithNamespaceId(ctx, namespaceId), nil
}

func setStorableEntIdent(findableEnt_p interface{}, keyOpt dstore.KeyOptions) error {
	if findableEnt_p != nil {
		switch findableEnt_p.(type) {
		case dstore.Findable:
			var err error
			if numId := keyOpt.NumberId(0); numId != 0 {
				//findableEnt_p.( dstore.Findable).FindBy(strconv.FormatInt(numId, 64))
				//panic(fmt.Sprintf("VVVV=%v e=%v", numId,findableEnt_p))
				_, err = findableEnt_p.( dstore.Findable).FindBy(fmt.Sprint(numId))
			} else {
				_, err = findableEnt_p.( dstore.Findable).FindBy(keyOpt.StringId(""))
			}
			return err
		}
	}
	return nil
}
