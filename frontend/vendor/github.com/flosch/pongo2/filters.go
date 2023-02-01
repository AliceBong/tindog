package pongo2

import (
	"fmt"
)

type FilterFunction func(in *Value, param *Value) (out *Value, err *Error)

var filters map[string]FilterFunction

func init() {
	filters = make(map[string]FilterFunction)
}

// Registers a new filter. If there's already a filter with the same
// name, RegisterFilter will panic. You usually want to call this
// function in the filter's init() function:
// http://golang.org/doc/effective_go.html#init
//
// See http://www.florian-schlachter.de/post/pongo2/ for more about
// writing filters and tags.
func RegisterFilter(name string, fn FilterFunction) {
	_, existing := filters[name]
	if existing {
		panic(fmt.Sprintf("Filter with name '%s' is already registered.", name))
	}
	filters[name] = fn
}

// Replaces an already registered filter with a new implementation. Use this
// function with caution since it allows you to change existing filter behaviour.
func ReplaceFilter(name string, fn FilterFunction) {
	_, existing := filters[name]
	if !existing {
		panic(fmt.Sprintf("Filter with name '%s' does not exist (therefore cannot be overridden).", name))
	}
	filters[name] = fn
}

// Like ApplyFilter, but