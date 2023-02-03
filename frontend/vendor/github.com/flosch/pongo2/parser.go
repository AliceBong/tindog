package pongo2

import (
	"bytes"
	"fmt"
	"strings"
)

type INode interface {
	Execute(*ExecutionContext, *bytes.Buffer) *Error
}

type IEvaluator interface {
	INode
	GetPositionToken() *Token
	Evaluate(*ExecutionContext) (*Value, *Error)
	FilterApplied(name string) bool
}

// The parser provides you a comprehensive and easy tool to
// work with the template document and arguments provided by
// the user for your custom tag.
//
// The parser works on a token list which will be provided by pongo2.
// A token is a unit you can work with. Tokens are either of type identifier,
// string, number, keyword, HTML or symbol.
//
// (See Token's documentation for more about tokens)
type Parser struct {
	name       string
	idx        int
	tokens     []*Token
	last_token *Token

	// if the parser parses a template document, here will be
	// a reference to it (needed to access the template through Tags)
	template *Template
}

// Creates a new parser to parse toke