package eval

import "github.com/yassinebenaid/nishimia/object"

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("invalid arguments count in function call, expected 1 argumets, got %d ",
					len(args),
				)
			}

			if str, ok := args[0].(*object.String); ok {
				return &object.Integer{
					Value: int64(len(str.Value)),
				}
			}

			return newError("argument to `len` not supported, got INTEGER")
		},
	},
}
