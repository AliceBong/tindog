package pongo2

import (
	"bytes"
	"math/rand"
	"strings"
	"time"
)

var (
	tagLoremParagraphs = strings.Split(tagLoremText, "\n")
	tagLoremWords      = strings.Fields(tagLoremText)
)

type tagLoremNode struct {
	position *Token
	count    int    // number of paragraphs
	method   string // w = words, p = HTML paragraphs, b = plain-text (default is b)
	random   bool   // does not use the default paragraph "Lorem ipsum dolor sit amet, ..."
}

func (node *tagLoremNode) Execute(ctx *ExecutionContext, buffer *bytes.Buffer) *Error {
	switch node.method {
	case "b":
		if node.random {
			for i := 0; i < node.count; i++ {
				if i > 0 {
					buffer.WriteString("\n")
				}
				par := tagLoremPara