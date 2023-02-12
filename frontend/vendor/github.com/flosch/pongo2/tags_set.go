package pongo2

import "bytes"

type tagSetNode struct {
	name       string
	expression IEvaluator
}

func (node *tagSetNode) Execute(ctx *ExecutionContext, buffer *bytes.Buffer) *Error {
	// Evaluate expression
	value, err := node.expression.Evaluate(ctx)
	if err != nil {
		return err
	}

	ctx.Private[node.name] = value
	return nil
}

func tagSetParser(doc *Parser, start *Token, arguments *Parser) (INodeTag, *Error) {
	node := &tagSetNode{}