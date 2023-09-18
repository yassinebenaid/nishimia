package object

func NewEnvirement() *Environment {
	return &Environment{Store: make(map[string]Object)}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvirement()
	env.outer = outer
	return env
}

type Environment struct {
	Store map[string]Object
	outer *Environment
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.Store[name]

	if !ok && e.outer != nil {
		return e.outer.Get(name)
	}

	return obj, ok
}

func (e *Environment) Has(name string) bool {
	_, ok := e.Store[name]

	return ok
}

func (e *Environment) Set(name string, value Object) Object {
	e.Store[name] = value
	return value
}
