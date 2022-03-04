package evaluator

import "github.com/smalldevshima/go-monkey/object"

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: B_LEN,
	},
}

var B_LEN object.BuiltinFunction = func(args ...object.Object) object.Object {
	argc := len(args)
	if argc != 1 {
		return newError(ERR_ARG_COUNT_MISMATCH, "len", 1, argc)
	}

	switch arg := args[0].(type) {
	case *object.String:
		return &object.Integer{Value: int64(len(arg.Value))}
	default:
		return newError(ERR_BUILTIN_TYPE_ERROR, 0, "len", object.O_STRING, arg.Type())
	}
}
