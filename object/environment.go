package object

/// Functions

func NewEnvironment() *Environment {
	return &Environment{store: make(map[string]Object)}
}

/// Types

type Environment struct {
	store map[string]Object
}

func (e *Environment) Get(name string) (obj Object, ok bool) {
	obj, ok = e.store[name]
	return
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}
