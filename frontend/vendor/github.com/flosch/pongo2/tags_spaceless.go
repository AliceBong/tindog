package pongo2

import (
	"bytes"
	"regexp"
)

type tagSpacelessNode struct {
	wrapper *NodeWrapper
}

var tagSpacelessRegexp = regexp.MustCompile(`(?U:(<.*>))([\t\n\v\f\r ]+)(?U:(<.*>))`)

func (node *tagSpacelessNode) Execute(ctx *ExecutionContext, buffer *bytes.Buffer) *Error {
	b := bytes.NewBuffer(make([]byte, 0, 1024)) // 1 KiB

	err := node.wrapper.Execute(ctx, b)
	if err != nil {
		return err
	}

	s := b.String()
	// Repeat this recursively
	changed := true
	for changed {
		s2 := tagSpacelessRegexp.ReplaceAllString(s, "$1$3")
		changed = s != s2
		s = s2
	}

	buffer.WriteString(s)

	return nil
}

func tagSpacelessParser(doc *Parser, start *Token, arguments *Parser) (INodeTag, *Error) {
	spaceless_node := &tagSpacelessNode{}

	wrapper, _, err := doc.WrapUntilTag("endspaceless")
	if err != nil {
		return nil, err
	}
	spaceless_node.wrapper = wrapper

	if arguments.Remaining() > 0 {
		return nil, arguments.Error("Malformed spaceless-