package pongo2

import (
	"bytes"
	"fmt"
)

type tagImportNode struct {
	position *Token
	filename string
	template *Template
	macros   map[string]*tagMacroNode // alias/name -> macro instance
}

func (node *tagImportNode) Execute(ctx *ExecutionContext, buffer *bytes.Buffer) *Error {
	for name, macro := range node.macros {
		func(name string, macro *tagMa