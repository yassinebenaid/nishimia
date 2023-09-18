package object

func NewEnvirement() *Envirement {
	return &Envirement{store: make(map[string]Object)}
}

func NewEnclosedEnvirement(outer *Envirement) *Envirement {
	env := NewEnvirement()
	env.outer = outer
	return env
}

type Envirement struct {
	store map[string]Object
	outer *Envirement
}

func (e *Envirement) Get(name string) (Object, bool) {
	obj, ok := e.store[name]

	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}

	return obj, ok
}

func (e *Envirement) Has(name string) bool {
	_, ok := e.store[name]

	return ok
}

func (e *Envirement) Set(name string, value Object) Object {
	e.store[name] = value
	return value
}
