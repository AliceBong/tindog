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

	if method_token := arguments.MatchType(TokenIdentifier); method_token != nil {
		if method_token.Val != "w" && method_token.Val != "p" && method_token.Val != "b" {
			return nil, arguments.Error("lorem-method must be either 'w', 'p' or 'b'.", nil)
		}

		lorem_node.method = method_token.Val
	}

	if arguments.MatchOne(TokenIdentifier, "random") != nil {
		lorem_node.random = true
	}

	if arguments.Remaining() > 0 {
		return nil, arguments.Error("Malformed lorem-tag arguments.", nil)
	}

	return lorem_node, nil
}

func init() {
	rand.Seed(time.Now().Unix())

	RegisterTag("lorem", tagLoremParser)
}

const tagLoremText = `Lorem ipsum dolor sit amet, consectetur adipisici elit, sed eiusmod tempor incidunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquid ex ea commodi consequat. Quis aute iure reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint obcaecat cupiditat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.
Duis autem vel eum iriure dolor in hendrerit in vulputate velit esse molestie consequat, vel illum dolore eu feugiat nulla facilisis at vero eros et accumsan et iusto odio dignissim qui blandit praesent luptatum zzril delenit augue duis dolore te feugait nulla facilisi. Lorem ipsum dolor sit amet, consectetuer adipiscing elit, sed diam nonummy nibh euismod tincidunt ut laoreet dolore magna aliquam erat volutpat.
Ut wisi enim ad minim veniam, quis nostrud exerci tation ullamcorper suscipit lobortis nisl ut aliquip ex ea commodo consequat. Duis autem vel eum iriure do