package pongo2

import (
	"bytes"
	"fmt"
	"io"
)

type Template struct {
	set *TemplateSet

	// Input
	is_tpl_string bool
	name          string
	tpl           string
	size          int

	// Calculation
	tokens []*Token
	parser *Parser

	// first come, first serve (it's important to not override existing entries in here)
	level           int
	parent          *Template
	child           *Template
	blocks          map[string]*NodeWrapper
	exported_macros map[string]*tagMacroNode

	// Output
	root *nodeDocument
}

func newTemplateString(set *TemplateSet, tpl string) (*Template, error) {
	return newTemplate(set, "<string>", true, tpl)
}

func newTemplate(set *TemplateSet, name string, is_tpl_string bool, tpl string) (*Template, error) {
	// Create the template
	t := &Template{
		set:             set,
		is_tpl_string:   is_tpl_string,
		name:            name,
		tpl:             tpl,
		size:            len(tpl),
		blocks:          make(map[string]*NodeWrapper),
		exported_macros: make(map[string]*tagMacroNode),
	}

	// Tokenize it
	tokens, err := lex(name, tpl)
	if err != nil {
		return nil