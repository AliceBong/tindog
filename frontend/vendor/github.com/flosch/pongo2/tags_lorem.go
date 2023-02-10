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
				par := tagLoremParagraphs[rand.Intn(len(tagLoremParagraphs))]
				buffer.WriteString(par)
			}
		} else {
			for i := 0; i < node.count; i++ {
				if i > 0 {
					buffer.WriteString("\n")
				}
				par := tagLoremParagraphs[i%len(tagLoremParagraphs)]
				buffer.WriteString(par)
			}
		}
	case "w":
		if node.random {
			for i := 0; i < node.count; i++ {
				if i > 0 {
					buffer.WriteString(" ")
				}
				word := tagLoremWords[rand.Intn(len(tagLoremWords))]
				buffer.WriteString(word)
			}
		} else {
			for i := 0; i < node.count; i++ {
				if i > 0 {
					buffer.WriteString(" ")
				}
				word := tagLoremWords[i%len(tagLoremWords)]
				buffer.WriteString(word)
			}
		}
	case "p":
		if node.random {
			for i := 0; i < node.count; i++ {
				if i > 0 {
					buffer.WriteString("\n")
				}
				buffer.WriteString("<p>")
				par := tagLoremParagraphs[rand.Intn(len(tagLoremParagraphs))]
				buffer.WriteString(par)
				buffer.WriteString("</p>")
			}
		} else {
			for i := 0; i < node.count; i++ {
				if i > 0 {
					buffer.WriteString("\n")
				}
				buffer.WriteString("<p>")
				par := tagLoremParagraphs[i%len(tagLoremParagraphs)]
				buffer.WriteString(par)
				buffer.WriteString("</p>")

			}
		}
	default:
		panic("unsupported method")
	}

	return nil
}

func tagLoremParser(doc *Parser, start *Token, arguments *Parser) (INodeTag, *Error) {
	lorem_node := &tagLoremNode{
		position: start,
		count:    1,
		method:   "b",
	}

	if count_token := arguments.MatchType(TokenNumber); count_token != nil {
		lorem_node.count = AsValue(count_token.Val).Integer()
	}

	if method_token := arg