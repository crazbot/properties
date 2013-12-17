// Copyright 2013 Frank Schroeder. All rights reserved. MIT licensed.

package properties

import (
	"fmt"
	// "log"
	"runtime"
)

type parser struct {
	lex *lexer
}

func newParser() *parser {
	return &parser{}
}

func (p *parser) Parse(input string) (props *Properties, err error) {
	// log.Printf("Parsing input '%s'", input)

	defer p.recover(&err)
	p.lex = lex(input)
	props = &Properties{m:make(map[string]string)}

	for {
		token := p.expectOneOf(itemKey, itemEOF)
		if token.typ == itemEOF {
			break
		}
		key := token.val
		token = p.expectOneOf(itemValue, itemEOF)
		if token.typ == itemEOF {
			props.Set(key, "")
			break
		}
		props.Set(key, token.val)
	}

	return props, nil
}

func (p *parser) errorf(format string, args ...interface{}) {
	format = fmt.Sprintf("properties: Line %d: %s", p.lex.lineNumber(), format)
	panic(fmt.Errorf(format, args...))
}

func (p *parser) expect(expected itemType) (token item) {
	token = p.lex.nextItem()
	if token.typ != expected {
		p.unexpected(token)
	}
	return token
}

func (p *parser) expectOneOf(expected1, expected2 itemType) (token item) {
	token = p.lex.nextItem()
	if token.typ != expected1 && token.typ != expected2 {
		p.unexpected(token)
	}
	return token
}

func (p *parser) unexpected(token item) {
	p.errorf(token.String())
}

// recover is the handler that turns panics into returns from the top level of Parse.
func (p *parser) recover(errp *error) {
	e := recover()
	if e != nil {
		if _, ok := e.(runtime.Error); ok {
			panic(e)
		}
		// if p != nil {
		// 	p.stopParse()
		// }
		*errp = e.(error)
	}
	return
}

// func (p *parser) stopParse() {

// }
