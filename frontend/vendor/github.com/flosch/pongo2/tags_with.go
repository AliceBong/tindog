package pongo2

import (
	"bytes"
)

type tagWithNode struct {
	with_pairs map[string]IEvaluator
	wrapper    *NodeWrapper
}

func (node *tagWithNode) Execute(ctx *ExecutionContext, buf