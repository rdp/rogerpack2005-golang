// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
	Data-driven templates for generating textual output such as
	HTML. See
		http://code.google.com/p/json-template/wiki/Reference
	for full documentation of the template language. A summary:

	Templates are executed by applying them to a data structure.
	Annotations in the template refer to elements of the data
	structure (typically a field of a struct) to control execution
	and derive values to be displayed.  The template walks the
	structure as it executes and the "cursor" @ represents the
	value at the current location in the structure.

	Data items may be values or pointers; the interface hides the
	indirection.

	Major constructs ({} are metacharacters; [] marks optional elements):

		{# comment }

	A one-line comment.

		{.section field} XXX [ {.or} YYY ] {.end}

	Set @ to the value of the field.  It may be an explicit @
	to stay at the same point in the data. If the field is nil
	or empty, execute YYY; otherwise execute XXX.

		{.repeated section field} XXX [ {.alternates with} ZZZ ] [ {.or} YYY ] {.end}

	Like .section, but field must be an array or slice.  XXX
	is executed for each element.  If the array is nil or empty,
	YYY is executed instead.  If the {.alternates with} marker
	is present, ZZZ is executed between iterations of XXX.

		{field}
		{field|formatter}

	Insert the value of the field into the output. Field is
	first looked for in the cursor, as in .section and .repeated.
	If it is not found, the search continues in outer sections
	until the top level is reached.
	
	If a formatter is specified, it must be named in the formatter
	map passed to the template set up routines or in the default
	set ("html","str","") and is used to process the data for
	output.  The formatter function has signature
		func(wr io.Write, data interface{}, formatter string)
	where wr is the destination for output, data is the field
	value, and formatter is its name at the invocation site.
*/

package template

import (
	"fmt";
	"io";
	"os";
	"reflect";
	"strings";
	"template";
	"container/vector";
)

// Errors returned during parsing. TODO: different error model for execution?

type ParseError struct {
	os.ErrorString
}

// Most of the literals are aces.
var lbrace = []byte{ '{' }
var rbrace = []byte{ '}' }
var space = []byte{ ' ' }
var tab = []byte{ '\t' }

// The various types of "tokens", which are plain text or (usually) brace-delimited descriptors
const (
	tokAlternates = iota;
	tokComment;
	tokEnd;
	tokLiteral;
	tokOr;
	tokRepeated;
	tokSection;
	tokText;
	tokVariable;
)

// FormatterMap is the type describing the mapping from formatter
// names to the functions that implement them.
type FormatterMap map[string] func(io.Write, interface{}, string)

// Built-in formatters.
var builtins = FormatterMap {
	"html" : HtmlFormatter,
	"str" : StringFormatter,
	"" : StringFormatter,
}

// The parsed state of a template is a vector of xxxElement structs.
// Sections have line numbers so errors can be reported better during execution.

// Plain text.
type textElement struct {
	text	[]byte;
}

// A literal such as .meta-left or .meta-right
type literalElement struct {
	text []byte;
}

// A variable to be evaluated
type variableElement struct {
	linenum	int;
	name	string;
	formatter	string;	// TODO(r): implement pipelines
}

// A .section block, possibly with a .or
type sectionElement struct {
	linenum int;	// of .section itself
	field	string;	// cursor field for this block
	start	int;	// first element
	or	int;	// first element of .or block
	end	int;	// one beyond last element
}

// A .repeated block, possibly with a .or and a .alternates
type repeatedElement struct {
	sectionElement;	// It has the same structure...
	altstart	int;	// ... except for alternates
	altend	int;
}

// Template is the type that represents a template definition.
// It is unchanged after parsing.
type Template struct {
	fmap	FormatterMap;	// formatters for variables
	// Used during parsing:
	ldelim, rdelim	[]byte;	// delimiters; default {}
	buf	[]byte;	// input text to process
	p	int;	// position in buf
	linenum	int;	// position in input
	errors	chan os.Error;	// for error reporting during parsing (only)
	// Parsed results:
	elems	*vector.Vector;
}

// Internal state for executing a Template.  As we evaluate the struct,
// the data item descends into the fields associated with sections, etc.
// Parent is used to walk upwards to find variables higher in the tree.
type state struct {
	parent	*state;	// parent in hierarchy
	data	reflect.Value;	// the driver data for this section etc.
	wr	io.Write;	// where to send output
	errors	chan os.Error;	// for reporting errors during execute
}

func (parent *state) clone(data reflect.Value) *state {
	return &state{parent, data, parent.wr, parent.errors}
}

// New creates a new template with the specified formatter map (which
// may be nil) to define auxiliary functions for formatting variables.
func New(fmap FormatterMap) *Template {
	t := new(Template);
	t.fmap = fmap;
	t.ldelim = lbrace;
	t.rdelim = rbrace;
	t.errors = make(chan os.Error);
	t.elems = vector.New(0);
	return t;
}

// Generic error handler, called only from execError or parseError.
func error(errors chan os.Error, line int, err string, args ...) {
	errors <- ParseError{fmt.Sprintf("line %d: %s", line, fmt.Sprintf(err, args))};
	sys.Goexit();
}

// Report error and stop executing.  The line number must  be provided explicitly.
func (t *Template) execError(st *state, line int, err string, args ...) {
	error(st.errors, line, err, args);
}

// Report error and stop parsing.  The line number comes from the template state.
func (t *Template) parseError(err string, args ...) {
	error(t.errors, t.linenum, err, args)
}

// -- Lexical analysis

// Is c a white space character?
func white(c uint8) bool {
	return c == ' ' || c == '\t' || c == '\r' || c == '\n'
}

// Safely, does s[n:n+len(t)] == t?
func equal(s []byte, n int, t []byte) bool {
	b := s[n:len(s)];
	if len(t) > len(b) {	// not enough space left for a match.
		return false
	}
	for i , c := range t {
		if c != b[i] {
			return false
		}
	}
	return true
}

// nextItem returns the next item from the input buffer.  If the returned
// item is empty, we are at EOF.  The item will be either a
// delimited string or a non-empty string between delimited
// strings. Tokens stop at (but include, if plain text) a newline.
// Action tokens on a line by themselves drop the white space on
// either side, up to and including the newline.
func (t *Template) nextItem() []byte {
	sawLeft := false;	// are we waiting for an opening delimiter?
	special := false;	// is this a {.foo} directive, which means trim white space?
	// Delete surrounding white space if this {.foo} is the only thing on the line.
	trim_white := t.p == 0 || t.buf[t.p-1] == '\n';
	only_white := true;	// we have seen only white space so far
	var i int;
	start := t.p;
Loop:
	for i = t.p; i < len(t.buf); i++ {
		switch {
		case t.buf[i] == '\n':
			t.linenum++;
			i++;
			break Loop;
		case white(t.buf[i]):
			// white space, do nothing
		case !sawLeft && equal(t.buf, i, t.ldelim):  // sawLeft checked because delims may be equal
			// anything interesting already on the line?
			if !only_white {
				break Loop;
			}
			// is it a directive or comment?
			j := i + len(t.ldelim);  // position after delimiter
			if j+1 < len(t.buf) && (t.buf[j] == '.' || t.buf[j] == '#') {
				special = true;
				if trim_white && only_white {
					start = i;
				}
			} else if i > t.p {  // have some text accumulated so stop before delimiter
				break Loop;
			}
			sawLeft = true;
			i = j - 1;
		case equal(t.buf, i, t.rdelim):
			if !sawLeft {
				t.parseError("unmatched closing delimiter")
			}
			sawLeft = false;
			i += len(t.rdelim);
			break Loop;
		default:
			only_white = false;
		}
	}
	if sawLeft {
		t.parseError("unmatched opening delimiter")
	}
	item := t.buf[start:i];
	if special && trim_white {
		// consume trailing white space
		for ; i < len(t.buf) && white(t.buf[i]); i++ {
			if t.buf[i] == '\n' {
				i++;
				break	// stop after newline
			}
		}
	}
	t.p = i;
	return item
}

// Turn a byte array into a white-space-split array of strings.
func words(buf []byte) []string {
	s := make([]string, 0, 5);
	p := 0; // position in buf
	// one word per loop
	for i := 0; ; i++ {
		// skip white space
		for ; p < len(buf) && white(buf[p]); p++ {
		}
		// grab word
		start := p;
		for ; p < len(buf) && !white(buf[p]); p++ {
		}
		if start == p {	// no text left
			break
		}
		if i == cap(s) {
			ns := make([]string, 2*cap(s));
			for j := range s {
				ns[j] = s[j]
			}
			s = ns;
		}
		s = s[0:i+1];
		s[i] = string(buf[start:p])
	}
	return s
}

// Analyze an item and return its token type and, if it's an action item, an array of
// its constituent words.
func (t *Template) analyze(item []byte) (tok int, w []string) {
	// item is known to be non-empty
	if !equal(item, 0, t.ldelim) {	// doesn't start with left delimiter
		tok = tokText;
		return
	}
	if !equal(item, len(item)-len(t.rdelim), t.rdelim) {	// doesn't end with right delimiter
		t.parseError("internal error: unmatched opening delimiter")	// lexing should prevent this
	}
	if len(item) <= len(t.ldelim)+len(t.rdelim) {	// no contents
		t.parseError("empty directive")
	}
	// Comment
	if item[len(t.ldelim)] == '#' {
		tok = tokComment;
		return
	}
	// Split into words
	w = words(item[len(t.ldelim): len(item)-len(t.rdelim)]);	// drop final delimiter
	if len(w) == 0 {
		t.parseError("empty directive")
	}
	if len(w) == 1 && w[0][0] != '.' {
		tok = tokVariable;
		return;
	}
	switch w[0] {
	case ".meta-left", ".meta-right", ".space", ".tab":
		tok = tokLiteral;
		return;
	case ".or":
		tok = tokOr;
		return;
	case ".end":
		tok = tokEnd;
		return;
	case ".section":
		if len(w) != 2 {
			t.parseError("incorrect fields for .section: %s", item)
		}
		tok = tokSection;
		return;
	case ".repeated":
		if len(w) != 3 || w[1] != "section" {
			t.parseError("incorrect fields for .repeated: %s", item)
		}
		tok = tokRepeated;
		return;
	case ".alternates":
		if len(w) != 2 || w[1] != "with" {
			t.parseError("incorrect fields for .alternates: %s", item)
		}
		tok = tokAlternates;
		return;
	}
	t.parseError("bad directive: %s", item);
	return
}

// -- Parsing

// Allocate a new variable-evaluation element.
func (t *Template) newVariable(name_formatter string) (v *variableElement) {
	name := name_formatter;
	formatter := "";
	bar := strings.Index(name_formatter, "|");
	if bar >= 0 {
		name = name_formatter[0:bar];
		formatter = name_formatter[bar+1:len(name_formatter)];
	}
	// Probably ok, so let's build it.
	v = &variableElement{t.linenum, name, formatter};

	// We could remember the function address here and avoid the lookup later,
	// but it's more dynamic to let the user change the map contents underfoot.
	// We do require the name to be present, though.

	// Is it in user-supplied map?
	if t.fmap != nil {
		if fn, ok := t.fmap[formatter]; ok {
			return
		}
	}
	// Is it in builtin map?
	if fn, ok := builtins[formatter]; ok {
		return
	}
	t.parseError("unknown formatter: %s", formatter);
	return
}

// Grab the next item.  If it's simple, just append it to the template.
// Otherwise return its details.
func (t *Template) parseSimple(item []byte) (done bool, tok int, w []string) {
	tok, w = t.analyze(item);
	done = true;	// assume for simplicity
	switch tok {
	case tokComment:
		return;
	case tokText:
		t.elems.Push(&textElement{item});
		return;
	case tokLiteral:
		switch w[0] {
		case ".meta-left":
			t.elems.Push(&literalElement{t.ldelim});
		case ".meta-right":
			t.elems.Push(&literalElement{t.rdelim});
		case ".space":
			t.elems.Push(&literalElement{space});
		case ".tab":
			t.elems.Push(&literalElement{tab});
		default:
			t.parseError("internal error: unknown literal: %s", w[0]);
		}
		return;
	case tokVariable:
		t.elems.Push(t.newVariable(w[0]));
		return;
	}
	return false, tok, w
}

// parseSection and parseRepeated are mutually recursive
func (t *Template) parseSection(words []string) *sectionElement

func (t *Template) parseRepeated(words []string) *repeatedElement {
	r := new(repeatedElement);
	t.elems.Push(r);
	r.linenum = t.linenum;
	r.field = words[2];
	// Scan section, collecting true and false (.or) blocks.
	r.start = t.elems.Len();
	r.or = -1;
	r.altstart = -1;
	r.altend = -1;
Loop:
	for {
		item := t.nextItem();
		if len(item) ==  0 {
			t.parseError("missing .end for .repeated section")
		}
		done, tok, w := t.parseSimple(item);
		if done {
			continue
		}
		switch tok {
		case tokEnd:
			break Loop;
		case tokOr:
			if r.or >= 0 {
				t.parseError("extra .or in .repeated section");
			}
			r.altend = t.elems.Len();
			r.or = t.elems.Len();
		case tokSection:
			t.parseSection(w);
		case tokRepeated:
			t.parseRepeated(w);
		case tokAlternates:
			if r.altstart >= 0 {
				t.parseError("extra .alternates in .repeated section");
			}
			if r.or >= 0 {
				t.parseError(".alternates inside .or block in .repeated section");
			}
			r.altstart = t.elems.Len();
		default:
			t.parseError("internal error: unknown repeated section item: %s", item);
		}
	}
	if r.altend < 0 {
		r.altend = t.elems.Len()
	}
	r.end = t.elems.Len();
	return r;
}

func (t *Template) parseSection(words []string) *sectionElement {
	s := new(sectionElement);
	t.elems.Push(s);
	s.linenum = t.linenum;
	s.field = words[1];
	// Scan section, collecting true and false (.or) blocks.
	s.start = t.elems.Len();
	s.or = -1;
Loop:
	for {
		item := t.nextItem();
		if len(item) ==  0 {
			t.parseError("missing .end for .section")
		}
		done, tok, w := t.parseSimple(item);
		if done {
			continue
		}
		switch tok {
		case tokEnd:
			break Loop;
		case tokOr:
			if s.or >= 0 {
				t.parseError("extra .or in .section");
			}
			s.or = t.elems.Len();
		case tokSection:
			t.parseSection(w);
		case tokRepeated:
			t.parseRepeated(w);
		case tokAlternates:
			t.parseError(".alternates not in .repeated");
		default:
			t.parseError("internal error: unknown section item: %s", item);
		}
	}
	s.end = t.elems.Len();
	return s;
}

func (t *Template) parse() {
	for {
		item := t.nextItem();
		if len(item) == 0 {
			break
		}
		done, tok, w := t.parseSimple(item);
		if done {
			continue
		}
		switch tok {
		case tokOr, tokEnd, tokAlternates:
			t.parseError("unexpected %s", w[0]);
		case tokSection:
			t.parseSection(w);
		case tokRepeated:
			t.parseRepeated(w);
		default:
			t.parseError("internal error: bad directive in parse: %s", item);
		}
	}
}

// -- Execution

// If the data for this template is a struct, find the named variable.
// The special name "@" (the "cursor") denotes the current data.
func (st *state) findVar(s string) reflect.Value {
	if s == "@" {
		return st.data
	}
	data := reflect.Indirect(st.data);
	typ, ok := data.Type().(reflect.StructType);
	if ok {
		for i := 0; i < typ.Len(); i++ {
			name, ftyp, tag, offset := typ.Field(i);
			if name == s {
				return data.(reflect.StructValue).Field(i)
			}
		}
	}
	return nil
}

// Is there no data to look at?
func empty(v reflect.Value, indirect_ok bool) bool {
	v = reflect.Indirect(v);
	if v == nil {
		return true
	}
	switch v.Type().Kind() {
	case reflect.StringKind:
		return v.(reflect.StringValue).Get() == "";
	case reflect.StructKind:
		return false;
	case reflect.ArrayKind:
		return v.(reflect.ArrayValue).Len() == 0;
	}
	return true;
}

// Look up a variable, up through the parent if necessary.
func (t *Template) varValue(v *variableElement, st *state) reflect.Value {
	field := st.findVar(v.name);
	if field == nil {
		if st.parent == nil {
			t.execError(st, t.linenum, "name not found: %s", v.name)
		}
		return t.varValue(v, st.parent);
	}
	return field;
}

// Evaluate a variable, looking up through the parent if necessary.
// If it has a formatter attached ({var|formatter}) run that too.
func (t *Template) writeVariable(v *variableElement, st *state) {
	formatter := v.formatter;
	val := t.varValue(v, st).Interface();
	// is it in user-supplied map?
	if t.fmap != nil {
		if fn, ok := t.fmap[v.formatter]; ok {
			fn(st.wr, val, v.formatter);
			return;
		}
	}
	// is it in builtin map?
	if fn, ok := builtins[v.formatter]; ok {
		fn(st.wr, val, v.formatter);
		return;
	}
	t.execError(st, v.linenum, "missing formatter %s for variable %s", v.formatter, v.name)
}

// execute{|Element|Section|Repeated} are mutually recursive
func (t *Template) executeSection(s *sectionElement, st *state)
func (t *Template) executeRepeated(r *repeatedElement, st *state)

// Execute element i.  Return next index to execute.
func (t *Template) executeElement(i int, st *state) int {
	switch elem := t.elems.At(i).(type) {
	case *textElement:
		st.wr.Write(elem.text);
		return i+1;
	case *literalElement:
		st.wr.Write(elem.text);
		return i+1;
	case *variableElement:
		t.writeVariable(elem, st);
		return i+1;
	case *sectionElement:
		t.executeSection(elem, st);
		return elem.end;
	case *repeatedElement:
		t.executeRepeated(elem, st);
		return elem.end;
	}
	e := t.elems.At(i);
	t.execError(st, 0, "internal error: bad directive in execute: %v %T\n", reflect.NewValue(e).Interface(), e);
	return 0
}

// Execute the template.
func (t *Template) execute(start, end int, st *state) {
	for i := start; i < end; {
		i = t.executeElement(i, st)
	}
}

// Execute a .section
func (t *Template) executeSection(s *sectionElement, st *state) {
	// Find driver data for this section.  It must be in the current struct.
	field := st.findVar(s.field);
	if field == nil {
		t.execError(st, s.linenum, ".section: cannot find field %s in %s", s.field, reflect.Indirect(st.data).Type());
	}
	st = st.clone(field);
	start, end := s.start, s.or;
	if !empty(field, true) {
		// Execute the normal block.
		if end < 0 {
			end = s.end
		}
	} else {
		// Execute the .or block.  If it's missing, do nothing.
		start, end = s.or, s.end;
		if start < 0 {
			return
		}
	}
	for i := start; i < end; {
		i = t.executeElement(i, st)
	}
}

// Execute a .repeated section
func (t *Template) executeRepeated(r *repeatedElement, st *state) {
	// Find driver data for this section.  It must be in the current struct.
	field := st.findVar(r.field);
	if field == nil {
		t.execError(st, r.linenum, ".repeated: cannot find field %s in %s", r.field, reflect.Indirect(st.data).Type());
	}
	field = reflect.Indirect(field);

	// Must be an array/slice
	if field != nil && field.Kind() != reflect.ArrayKind {
		t.execError(st, r.linenum, ".repeated: %s has bad type %s", r.field, field.Type());
	}
	if empty(field, true) {
		// Execute the .or block, once.  If it's missing, do nothing.
		start, end := r.or, r.end;
		if start >= 0 {
			newst := st.clone(field);
			for i := start; i < end; {
				i = t.executeElement(i, newst)
			}
		}
		return
	}
	// Execute the normal block.
	start, end := r.start, r.or;
	if end < 0 {
		end = r.end
	}
	if r.altstart >= 0 {
		end = r.altstart
	}
	if field != nil {
		array := field.(reflect.ArrayValue);
		for j := 0; j < array.Len(); j++ {
			newst := st.clone(array.Elem(j));
			for i := start; i < end; {
				i = t.executeElement(i, newst)
			}
			// If appropriate, do .alternates between elements
			if j < array.Len() - 1 && r.altstart >= 0 {
				for i := r.altstart; i < r.altend; i++ {
					i = t.executeElement(i, newst)
				}
			}
		}
	}
}

// A valid delimiter must contain no white space and be non-empty.
func validDelim(d []byte) bool {
	if len(d) == 0 {
		return false
	}
	for i, c := range d {
		if white(c) {
			return false
		}
	}
	return true;
}

// -- Public interface

// Parse initializes a Template by parsing its definition.  The string
// s contains the template text.  If any errors occur, Parse returns
// the error.
func (t *Template) Parse(s string) os.Error {
	if !validDelim(t.ldelim) || !validDelim(t.rdelim) {
		return ParseError{fmt.Sprintf("bad delimiter strings %q %q", t.ldelim, t.rdelim)}
	}
	t.buf = io.StringBytes(s);
	t.p = 0;
	t.linenum = 0;
	go func() {
		t.parse();
		t.errors <- nil;	// clean return;
	}();
	return <-t.errors;
}

// Execute applies a parsed template to the specified data object,
// generating output to wr.
func (t *Template) Execute(data interface{}, wr io.Write) os.Error {
	// Extract the driver data.
	val := reflect.NewValue(data);
	errors := make(chan os.Error);
	go func() {
		t.p = 0;
		t.execute(0, t.elems.Len(), &state{nil, val, wr, errors});
		errors <- nil;	// clean return;
	}();
	return <-errors;
}

// SetDelims sets the left and right delimiters for operations in the
// template.  They are validated during parsing.  They could be
// validated here but it's better to keep the routine simple.  The
// delimiters are very rarely invalid and Parse has the necessary
// error-handling interface already.
func (t *Template) SetDelims(left, right string) {
	t.ldelim = io.StringBytes(left);
	t.rdelim = io.StringBytes(right);
}

// Parse creates a Template with default parameters (such as {} for
// metacharacters).  The string s contains the template text while
// the formatter map fmap, which may be nil, defines auxiliary functions
// for formatting variables.  The template is returned. If any errors
// occur, err will be non-nil.
func Parse(s string, fmap FormatterMap) (t *Template, err os.Error) {
	t = New(fmap);
	err = t.Parse(s);
	return
}
