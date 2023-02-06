package pongo2

import (
	"bytes"
)

type tagIfNode struct {
	conditions []IEvaluator
	wrappers   []*NodeWrapper
}

func (node *tagIfNode) Execute(ctx *ExecutionContext, buffer *bytes.Buffer) *Error {
	for i, condition := range node.conditions {
		result, err := condition.Evaluate(ctx)
		if err != nil {
			return err
		}

		if result.IsTrue() {
			return node.wrappers[i].Execute(ctx, buffer)
		} else {
			// Last condition?
			if len(node.conditions) == i+1 && len(node.wrappers) > i+1 {
				return node.wrappers[i+1].Execute(ctx, buffer)
			}
		}
	}
	return nil
}

func tagIfParser(doc *Parser, start *Token, arguments *Parser) (INodeTag, *Error) {
	if_node := &tagIfNode{}

	// Parse first and main IF condition
	condition, err := arguments.ParseExpression()
	if err != nil {
		return nil, err
	}
	if_node.cond