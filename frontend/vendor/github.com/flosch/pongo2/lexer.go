
package pongo2

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

const (
	TokenError = iota
	EOF

	TokenHTML

	TokenKeyword
	TokenIdentifier
	TokenString
	TokenNumber
	TokenSymbol
)

var (
	tokenSpaceChars                = " \n\r\t"
	tokenIdentifierChars           = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_"
	tokenIdentifierCharsWithDigits = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_0123456789"
	tokenDigits                    = "0123456789"

	// Available symbols in pongo2 (within filters/tag)
	TokenSymbols = []string{
		// 3-Char symbols

		// 2-Char symbols
		"==", ">=", "<=", "&&", "||", "{{", "}}", "{%", "%}", "!=", "<>",

		// 1-Char symbol
		"(", ")", "+", "-", "*", "<", ">", "/", "^", ",", ".", "!", "|", ":", "=", "%",
	}

	// Available keywords in pongo2
	TokenKeywords = []string{"in", "and", "or", "not", "true", "false", "as", "export"}
)