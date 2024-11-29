package interp

import "toterich/golox/ast"

// Contains the current state of the interpreter.
// environment provides a stack-like interface to push and pop sub-envs, which are used
// to implement Scoping. Every sub-env inherits all state from the envs below it on the stack,
// so that each scope can access the state of all parent scopes.
// Shadowing is supported, meaning that if an identifier is redeclared inside a nested scope, accessing
// this identifier will return the value of the nested scope until that scope is popped. Then, accesses
// to the identifier will return the value from the surrounding scope.
type environment struct {
	scopes []map[string]ast.LoxValue
}

func newEnvironment() environment {
	// We always have at least the global scope
	return environment{scopes: []map[string]ast.LoxValue{}}
}

// Query the value of an identifier, starting with the current scope and moving up the stack.
// If the identifier does not exist in any scope, the second return parameter is false.
func (env environment) getVar(ident string) (ast.LoxValue, bool) {
	// Iterate backwards through the scopes so the most deeply nested ones are queried first
	for i := len(env.scopes) - 1; i >= 0; i -= 1 {
		val, ok := env.scopes[i][ident]
		if ok {
			return val, ok
		}
	}

	return ast.NewNilValue(), false
}

// Set the value of an identifier in the current scope.
func (env *environment) setVar(ident string, value ast.LoxValue) {
	env.scopes[len(env.scopes)-1][ident] = value
}

// Push a new scope
func (env *environment) push() {
	env.scopes = append(env.scopes, map[string]ast.LoxValue{})
}

// Pop the most recent scope. All variables declared in this scope are discarded.
func (env *environment) pop() {
	env.scopes = env.scopes[:len(env.scopes)-1]
}
