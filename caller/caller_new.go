package caller

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/gontainer/grouperror"
	"github.com/gontainer/reflectpro/caller/internal/caller"
	intReflect "github.com/gontainer/reflectpro/internal/reflect"
)

// NewCallMethod replace will replace [CallMethod] and [ForceCallMethod].
//
//nolint:wrapcheck
func NewCallMethod(object any, method string, args []any, convertArgs bool) (_ []any, err error) {
	defer func() {
		if err != nil {
			err = grouperror.Prefix(fmt.Sprintf("cannot call method (%T).%+q: ", object, method), err)
		}
	}()

	fn, err := caller.Method(object, method)
	if err != nil {
		if errors.Is(err, caller.ErrInvalidMethod) {
			chain, chainErr := intReflect.ValueToKindChain(reflect.ValueOf(object))
			if chainErr == nil && chain.Prefixed(reflect.Ptr) {
				return caller.ValidateAndForceCallMethod(object, method, args, convertArgs, caller.DontValidate)
			}
		}

		return nil, err
	}

	return caller.CallFunc(fn, args, convertArgs)
}
