package object

/// Functions

func NewEnvironment() *Environment {
	return &Environment{store: make(map[string]Object), outer: nil}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

/// Types

type Environment struct {
	store map[string]Object
	outer *Environment
}

func (e *Environment) Get(name string) (obj Object, ok bool) {
	obj, ok = e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}
