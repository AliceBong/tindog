package pongo2

import (
	"bytes"
	"fmt"
	"math"
)

type tagWidthratioNode struct {
	position     *Token
	current, max IEvaluator
	width        IEvaluator
	ctx_name     string
}

func (node *tagWidthratioNode) Execute(ctx *ExecutionContext, buffer *bytes.Buffer) *Error {
	current, err := node