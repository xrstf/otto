// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package types

import (
	"go.xrstf.de/rudi/pkg/coalescing"
	"go.xrstf.de/rudi/pkg/lang/ast"
)

type Document struct {
	data any
}

func NewDocument(data any) (Document, error) {
	return Document{
		data: data,
	}, nil
}

func (d *Document) Data() any {
	return d.data
}

func (d *Document) Set(wrappedData any) {
	d.data = wrappedData
}

type Context struct {
	document  *Document
	funcs     Functions
	variables Variables
	coalescer coalescing.Coalescer
}

func NewContext(doc Document, variables Variables, funcs Functions, coalescer coalescing.Coalescer) Context {
	if funcs == nil {
		funcs = NewFunctions()
	}

	if variables == nil {
		variables = NewVariables()
	}

	if coalescer == nil {
		coalescer = coalescing.NewStrict()
	}

	return Context{
		document:  &doc,
		funcs:     funcs,
		variables: variables,
		coalescer: coalescer,
	}
}

// Coalesce is named this way to make the frequent calls read fluently
// (for example "ctx.Coalesce().ToBool(...)").
func (c Context) Coalesce() coalescing.Coalescer {
	return c.coalescer
}

func (c Context) GetDocument() *Document {
	return c.document
}

func (c Context) GetVariable(name string) (any, bool) {
	return c.variables.Get(name)
}

func (c Context) GetFunction(name string) (Function, bool) {
	return c.funcs.Get(name)
}

func (c Context) WithVariable(name string, val any) Context {
	return Context{
		document:  c.document,
		funcs:     c.funcs,
		variables: c.variables.With(name, val),
		coalescer: c.coalescer,
	}
}

func (c Context) WithCoalescer(coalescer coalescing.Coalescer) Context {
	return Context{
		document:  c.document,
		funcs:     c.funcs,
		variables: c.variables,
		coalescer: coalescer,
	}
}

type Function interface {
	Evaluate(ctx Context, args []ast.Expression) (any, error)

	// Description returns a short, one-line description of the function; markdown
	// can be used to highlight other function names, like "behaves similar
	// to `foo`, but …".
	Description() string
}

type TupleFunction func(ctx Context, args []ast.Expression) (any, error)

type basicFunc struct {
	f    TupleFunction
	desc string
}

func BasicFunction(f TupleFunction, description string) Function {
	return basicFunc{
		f:    f,
		desc: description,
	}
}

var _ Function = basicFunc{}

func (f basicFunc) Evaluate(ctx Context, args []ast.Expression) (any, error) {
	return f.f(ctx, args)
}

func (f basicFunc) Description() string {
	return f.desc
}

type Functions map[string]Function

func NewFunctions() Functions {
	return Functions{}
}

func (f Functions) Get(name string) (Function, bool) {
	variable, exists := f[name]
	return variable, exists
}

// Set sets/replaces the function in the current set (in-place).
// The function returns the same Functions to allow fluent access.
func (f Functions) Set(name string, fun Function) Functions {
	f[name] = fun
	return f
}

// Set removes a function from the set.
// The function returns the same Functions to allow fluent access.
func (f Functions) Delete(name string) Functions {
	delete(f, name)
	return f
}

// Add adds all functions from other to the current set.
// The function returns the same Functions to allow fluent access.
func (f Functions) Add(other Functions) Functions {
	for name, fun := range other {
		f[name] = fun
	}
	return f
}

// Remove removes all functions from this set that are part of the other set,
// to enable constructs like AllFunctions.Remove(MathFunctions)
// The function returns the same Functions to allow fluent access.
func (f Functions) Remove(other Functions) Functions {
	for name := range other {
		f.Delete(name)
	}
	return f
}

func (f Functions) DeepCopy() Functions {
	result := NewFunctions()
	for key, val := range f {
		result[key] = val
	}
	return result
}

type Variables map[string]any

func NewVariables() Variables {
	return Variables{}
}

func (v Variables) Get(name string) (any, bool) {
	variable, exists := v[name]
	return variable, exists
}

// Set sets/replaces the variable value in the current set (in-place).
// The function returns the same variables to allow fluent access.
func (v Variables) Set(name string, val any) Variables {
	v[name] = val
	return v
}

// With returns a copy of the variables, with the new variable being added to it.
func (v Variables) With(name string, val any) Variables {
	return v.DeepCopy().Set(name, val)
}

func (v Variables) DeepCopy() Variables {
	result := NewVariables()
	for key, val := range v {
		result[key] = val
	}
	return result
}

func MakeShim(val any) ast.Shim {
	return ast.Shim{Value: val}
}
