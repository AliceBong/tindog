package pongo2

import (
	"bytes"
)

type tagIncludeNode struct {
	tpl                *Template
	filename_evaluator IEvaluator
	lazy               bool
	only               bool
	filename           string
	with_pairs         map[string]IEvaluator
}

func (node *tagIncludeNode) Execute(ctx *ExecutionContext, buffer *bytes.Buffer) *Error {
	// Building the context for the template
	include_ctx := make(Context)

	// Fill the context with all data from the parent
	if !node.only {
		include_ctx.Update(ctx.Public)
		include_ctx.Update(ctx.Private)
	}

	// Put all custom with-pairs into the context
	for key, value := range node.with_pairs {
		val, err := value.Evaluate(ctx)
		if err != nil {
			return err
		}
		include_ctx[key] = val
	}

	// Execute the template
	if node.lazy {
		// Evaluate the filename
		filename, err := node.filename_evaluator.Evaluate(ctx)
		if err != nil {
			return err
		}

		if filename.String() == "" {
			return ctx.Error("Filename for 'include'-tag evaluated to an empty string.", nil)
		}

		// Get include-filename
		included_filename := ctx.template.set.resolveFilename(ctx.template, filename.String())

		included_tpl, err2 := ctx.template.set.FromFile(included_filename)
		if e