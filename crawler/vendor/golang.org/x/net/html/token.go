// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package html

import (
	"bytes"
	"errors"
	"io"
	"strconv"
	"strings"

	"golang.org/x/net/html/atom"
)

// A TokenType is the type of a Token.
type TokenType uint32

const (
	// ErrorToken means that an error occurred during tokenization.
	ErrorToken TokenType = iota
	// TextToken means a text node.
	TextToken
	// A StartTagToken looks like <a>.
	StartTagToken
	// An EndTagToken looks like </a>.
	EndTagToken
	// A SelfClosingTagToken tag looks like <br/>.
	SelfClosingTagToken
	// A CommentToken looks like <!--x-->.
	CommentToken
	// A DoctypeToken looks like <!DOCTYPE x>
	DoctypeToken
)

// ErrBufferExceeded means that the buffering limit was exceeded.
var ErrBufferExceeded = errors.New("max buffer exceeded")

// String returns a string representation of the TokenType.
func (t TokenType) String() string {
	switch t {
	case ErrorToken:
		return "Error"
	case TextToken:
		return "Text"
	case StartTagToken:
		return "StartTag"
	case EndTagToken:
		return "EndTag"
	case SelfClosingTagToken:
		return "SelfClosingTag"
	case CommentToken:
		return "Comment"
	case DoctypeToken:
		return "Doctype"
	}
	return "Invalid(" + strconv.Itoa(int(t)) + ")"
}

// An Attribute is an attribute namespace-key-value triple. Namespace is
// non-empty for foreign attributes like xlink, Key is alphabetic (and hence
// does not contain escapable characters like '&', '<' or '>'), and Val is
// unescaped (it looks like "a<b" rather than "a&lt;b").
//
// Namespace is only used by the parser, not the tokenizer.
type Attribute struct {
	Namespace, Key, Val string
}

// A Token consists of a TokenType and some Data (tag name for start and end
// tags, content for text, comments and doctypes). A tag Token may also contain
// a slice of Attributes. Data is unescaped for all Tokens (it looks like "a<b"
// rather than "a&lt;b"). For tag Tokens, DataAtom is the atom for Data, or
// zero if Data is not a known tag name.
type Token struct {
	Type     TokenType
	DataAtom atom.Atom
	Data     string
	Attr     []Attribute
}

// tagString returns a string representation of a tag Token's Data and Attr.
func (t Token) tagString() string {
	if len(t.Attr) == 0 {
		return t.Data
	}
	buf := bytes.NewBufferString(t.Data)
	for _, a := range t.Attr {
		buf.WriteByte(' ')
		buf.WriteString(a.Key)
		buf.WriteString(`="`)
		escape(buf, a.Val)
		buf.WriteByte('"')
	}
	return buf.String()
}

// String returns a string representation of the Token.
func (t Token) String() string {
	switch t.Type {
	case ErrorToken:
		return ""
	case TextToken:
		return EscapeString(t.Data)
	case StartTagToken:
		return "<" + t.tagString() + ">"
	case EndTagToken:
		return "</" + t.tagString() + ">"
	case SelfClosingTagToken:
		return "<" + t.tagString() + "/>"
	case CommentToken:
		return "<!--" + t.Data + "-->"
	case DoctypeToken:
		return "<!DOCTYPE " + t.Data + ">"
	}
	return "Invalid(" + strconv.Itoa(int(t.Type)) + ")"
}

// span is a range of bytes in a Tokenizer's buffer. The start is inclusive,
// the end is exclusive.
type span struct {
	start, end int
}

// A Tokenizer returns a stream of HTML Tokens.
type Tokenizer struct {
	// r is the source of the HTML text.
	r io.Reader
	// tt is the TokenType of the current token.
	tt TokenType
	// err is the first error encountered during tokenization. It is possible
	// for tt != Error && err != nil to hold: this means that Next returned a
	// valid token but the subsequent Next call will return an error token.
	// For example, if the HTML text input was just "plain", then the first
	// Next call would set z.err to io.EOF but return a TextToken, and all
	// subsequent Next calls would return an ErrorToken.
	// err is never reset. Once it becomes non-nil, it stays non-nil.
	err error
	// readErr is the error returned by the io.Reader r. It is separate from
	// err because it is valid for an io.Reader to return (n int, err1 error)
	// such that n > 0 && err1 != nil, and callers should always process the
	// n > 0 bytes before considering the error err1.
	readErr error
	// buf[raw.start:raw.end] holds the raw bytes of the current token.
	// buf[raw.end:] is buffered input that will yield future tokens.
	raw span
	buf []byte
	// maxBuf limits the data buffered in buf. A value of 0 means unlimited.
	maxBuf int
	// buf[data.start:data.end] holds the raw bytes of the current token's data:
	// a text token's text, a tag token's tag name, etc.
	data span
	// pendingAttr is the attribute key and value currently being tokenized.
	// When complete, pendingAttr is pushed onto attr. nAttrReturned is
	// incremented on each call to TagAttr.
	pendingAttr   [2]span
	attr          [][2]span
	nAttrReturned int
	// rawTag is the "script" in "</script>" that closes the next token. If
	// non-empty, the subsequent call to Next will return a raw or RCDATA text
	// token: one that treats "<p>" as text instead of an element.
	// rawTag's contents are lower-cased.
	rawTag string
	// textIsRaw is whether the current text token's data is not escaped.
	textIsRaw bool
	// convertNUL is whether NUL bytes in the current token's data should
	// be converted into \ufffd replacement characters.
	convertNUL bool
	// allowCDATA is whether CDATA sections are allowed in the current context.
	allowCDATA bool
}

// AllowCDATA sets whether or not the tokenizer recognizes <![CDATA[foo]]> as
// the text "foo". The default value is false, which means to recognize it as
// a bogus comment "<!-- [CDATA[foo]] -->" instead.
//
// Strictly speaking, an HTML5 compliant tokenizer should allow CDATA if and
// only if tokenizing foreign content, such as MathML and SVG. However,
// tracking foreign-contentness is difficult to do purely in the tokenizer,
// as opposed to the parser, due to HTML integration points: an <svg> element
// can contain a <foreignObject> that is foreign-to-SVG but not foreign-to-
// HTML. For strict compliance with the HTML5 tokenization algorithm, it is the
// responsibility of the user of a tokenizer to call AllowCDATA as appropriate.
// In practice, if using the tokenizer without caring whether MathML or SVG
// CDATA is text or comments, such as tokenizing HTML to find all the anchor
// text, it is acceptable to ignore this responsibility.
func (z *Tokenizer) AllowCDATA(allowCDATA bool) {
	z.allowCDATA = allowCDATA
}

// NextIsNotRawText instructs the tokenizer that the next token should not be
// considered as 'raw text'. Some elements, such as script and title elements,
// normally require the next token after the opening tag to be 'raw text' that
// has no child elements. For example, tokenizing "<title>a<b>c</b>d</title>"
// yields a start tag token for "<title>", a text token for "a<b>c</b>d", and
// an end tag token for "</title>". There are no distinct start tag or end tag
// tokens for the "<b>" and "</b>".
//
// This tokenizer implementation will generally look for raw text at the right
// times. Strictly speaking, an HTML5 compliant tokenizer should not look for
// raw text if in foreign content: <title> generally needs raw text, but a
// <title> inside an <svg> does not. Another example is that a <textarea>
// generally needs raw text, but a <textarea> is not allowed as an immediate
// child of a <select>; in normal parsing, a <textarea> implies </select>, but
// one cannot close the implicit element when parsing a <select>'s InnerHTML.
// Similarly to AllowCDATA, tracking the correct moment to override raw-text-
// ness is difficult to do purely in the tokenizer, as opposed to the parser.
// For strict compliance with the HTML5 tokenization algorithm, it is the
// responsibility of the user of a tokenizer to call NextIsNotRawText as
// appropriate. In practice, like AllowCDATA, it is acceptable to ignore this
// responsibility for basic usage.
//
// Note that this 'raw text' concept is different from the one offered by the
// Tokenizer.Raw method.
func (z *Tokenizer) NextIsNotRawText() {
	z.rawTag = ""
}

// Err returns the error associated with the most recent ErrorToken token.
// This is typically io.EOF, meaning the end of tokenization.
func (z *Tokenizer) Err() error {
	if z.tt != ErrorToken {
		return nil
	}
	return z.err
}

// readByte returns the next byte from the input stream, doing a buffered read
// from z.r into z.buf if necessary. z.buf[z.raw.start:z.raw.end] remains a contiguous byte
// slice that holds all the bytes read so far for the current token.
// It sets z.err if the underlying reader returns an error.
// Pre-condition: z.err == nil.
func (z *Tokenizer) readByte() byte {
	if z.raw.end >= len(z.buf) {
		// Our buffer is exhausted and we have to read from z.r. Check if the
		// previous read resulted in an error.
		if z.readErr != nil {
			z.err = z.readErr
			return 0
		}
		// We copy z.buf[z.raw.start:z.raw.end] to the beginning of z.buf. If the length
		// z.raw.end - z.raw.start is more than half the capacity of z.buf, then we
		// allocate a new buffer before the copy.
		c := cap(z.buf)
		d := z.raw.end - z.raw.start
		var buf1 []byte
		if 2*d > c {
			buf1 = make([]byte, d, 2*c)
		} else {
			buf1 = z.buf[:d]
		}
		copy(buf1, z.buf[z.raw.start:z.raw.end])
		if x := z.raw.start; x != 0 {
			// Adjust the data/attr spans to refer to the same contents after the copy.
			z.data.start -= x
			z.data.end -= x
			z.pendingAttr[0].start -= x
			z.pendingAttr[0].end -= x
			z.pendingAttr[1].start -= x
			z.pendingAttr[1].end -= x
			for i := range z.attr {
				z.attr[i][0].start -= x
				z.attr[i][0].end -= x
				z.attr[i][1].start -= x
				z.attr[i][1].end -= x
			}
		}
		z.raw.start, z.raw.end, z.buf = 0, d, buf1[:d]
		// Now that we have copied the live bytes to the start of the buffer,
		// we read from z.r into the remainder.
		var n int
		n, z.readErr = readAtLeastOneByte(z.r, buf1[d:cap(buf1)])
		if n == 0 {
			z.err = z.readErr
			return 0
		}
		z.buf = buf1[:d+n]
	}
	x := z.buf[z.raw.end]
	z.raw.end++
	if z.maxBuf > 0 && z.raw.end-z.raw.start >= z.maxBuf {
		z.err = ErrBufferExceeded
		return 0
	}
	return x
}

// Buffered returns a slice containing data buffered but not yet tokenized.
func (z *Tokenizer) Buffered() []byte {
	return z.buf[z.raw.end:]
}

// readAtLeastOneByte wraps an io.Reader so that reading cannot return (0, nil).
// It returns io.ErrNoProgress if the underlying r.Read method returns (0, nil)
// too many times in succession.
func readAtLeastOneByte(r io.Reader, b []byte) (int, error) {
	for i := 0; i < 100; i++ {
		n, err := r.Read(b)
		if n != 0 || err != nil {
			return n, err
		}
	}
	return 0, io.ErrNoProgress
}

// skipWhiteSpace skips past any white space.
func (z *Tokenizer) skipWhiteSpace() {
	if z.err != nil {
		return
	}
	for {
		c := z.readByte()
		if z.err != nil {
			return
		}
		switch c {
		case ' ', '\n', '\r', '\t', '\f':
			// No-op.
		default:
			z.raw.end--
			return
		}
	}
}

// readRawOrRCDATA reads until the next "</foo>", where "foo" is z.rawTag and
// is typically something like "script" or "textarea".
func (z *Tokenizer) readRawOrRCDATA() {
	if z.rawTag == "script" {
		z.readScript()
		z.textIsRaw = true
		z.rawTag = ""
		return
	}
loop:
	for {
		c := z.readByte()
		if z.err != nil {
			break loop
		}
		if c != '<' {
			continue loop
		}
		c = z.readByte()
		if z.err != nil {
			break loop
		}
		if c != '/' {
			continue loop
		}
		if z.readRawEndTag() || z.err != nil {
			break loop
		}
	}
	z.data.end = z.raw.end
	// A textarea's or title's RCDATA can contain escaped entities.
	z.textIsRaw = z.rawTag != "textarea" && z.rawTag != "title"
	z.rawTag = ""
}

// readRawEndTag attempts to read a tag like "</foo>", where "foo" is z.rawTag.
// If it succeeds, it backs up the input position to reconsume the tag and
// returns true. Otherwise it returns false. The opening "</" has already been
// consumed.
func (z *Tokenizer) readRawEndTag() bool {
	for i := 0; i < len(z.rawTag); i++ {
		c := z.readByte()
		if z.err != nil {
			return false
		}
		if c != z.rawTag[i] && c != z.rawTag[i]-('a'-'A') {
			z.raw.end--
			return false
		}
	}
	c := z.readByte()
	if z.err != nil {
		return false
	}
	switch c {
	case ' ', '\n', '\r', '\t', '\f', '/', '>':
		// The 3 is 2 for the leading "</" plus 1 for the trailing character c.
		z.raw.end -= 3 + len(z.rawTag)
		return true
	}
	z.raw.end--
	return false
}

// readScript reads until the next </script> tag, following the byzantine
// rules for escaping/hiding the closing ta