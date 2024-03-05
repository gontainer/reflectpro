package caller

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/gontainer/grouperror"
	"github.com/gontainer/reflectpro/caller/internal/caller"
)

func isPtr(v any) bool {
	return reflect.ValueOf(v).Kind() == reflect.Ptr
}

//nolint:wrapcheck
func callMethod(
	object any,
	method string,
	args []any,
	convertArgs bool,
	validator caller.FuncValidator,
) (
	_ []any,
	err error,
) {
	defer func() {
		if err != nil {
			err = grouperror.Prefix(fmt.Sprintf("cannot call method (%T).%+q: ", object, method), err)
		}
	}()

	fn, err := caller.Method(object, method)
	if err != nil {
		if errors.Is(err, caller.ErrInvalidMethod) && isPtr(object) {
			return caller.ValidateAndForceCallMethod(object, method, args, convertArgs, validator)
		}

		return nil, err
	}

	if validator != nil {
		if err := validator.Validate(fn); err != nil {
			return nil, err
		}
	}

	return caller.CallFunc(fn, args, convertArgs)
}
