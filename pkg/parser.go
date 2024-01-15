// Package pkg Copyright 2024 Lars Wilhelmsen <sral-backwards@sral.org>. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.
package pkg

import "fmt"

// Assuming value and Lexer are already defined as in the previous lexer translation

// ParserError represents errors that can occur during parsing.
type ParserError struct {
	Message string
}

func (e ParserError) Error() string {
	return e.Message
}

// AuthorizationExpression is an interface for different expression types.
type AuthorizationExpression interface {
	evaluate(authorizations map[string]bool) bool
}

// AndExpression represents an AND expression.
type AndExpression struct {
	nodes []AuthorizationExpression
}

func (a AndExpression) evaluate(authorizations map[string]bool) bool {
	for _, node := range a.nodes {
		if !node.evaluate(authorizations) {
			return false
		}
	}
	return true
}

// Implement evaluate, ToJSONStr, and Normalize for AndExpression...

// OrExpression represents an OR expression.
type OrExpression struct {
	nodes []AuthorizationExpression
}

func (o OrExpression) evaluate(authorizations map[string]bool) bool {
	for _, node := range o.nodes {
		if node.evaluate(authorizations) {
			return true
		}
	}
	return false
}

// AccessTokenExpression represents an access token.
type AccessTokenExpression struct {
	value string
}

func (a AccessTokenExpression) evaluate(authorizations map[string]bool) bool {
	return authorizations[a.value]
}

// Implement evaluate, ToJSONStr, and Normalize for AccessTokenExpression...

// Scope is used during parsing to build up expressions.
type Scope struct {
	nodes    []AuthorizationExpression
	labels   []AccessTokenExpression
	operator Token
}

func newScope() *Scope {
	return &Scope{
		nodes:    make([]AuthorizationExpression, 0),
		labels:   make([]AccessTokenExpression, 0),
		operator: None, // Assuming None is a defined value value
	}
}

func (s *Scope) addNode(node AuthorizationExpression) {
	s.nodes = append(s.nodes, node)
}

func (s *Scope) addLabel(label string) {
	s.labels = append(s.labels, AccessTokenExpression{value: label})
}

func (s *Scope) setOperator(operator Token) error {
	if s.operator == None {
		s.operator = operator
	} else if s.operator != operator {
		return ParserError{Message: "unexpected operator"}
	}
	return nil
}

func (s *Scope) Build() (AuthorizationExpression, error) {
	if len(s.labels) == 1 && len(s.nodes) == 0 {
		return s.labels[0], nil
	}

	if len(s.nodes) == 1 && len(s.labels) == 0 {
		return s.nodes[0], nil
	}

	if s.operator == None {
		return nil, ParserError{Message: "missing operator"}
	}
	// combine nodes and labels into one slice
	combined := make([]AuthorizationExpression, 0, len(s.nodes)+len(s.labels))
	for _, node := range s.nodes {
		combined = append(combined, node)
	}
	for _, label := range s.labels {
		combined = append(combined, label)
	}
	if s.operator == And {
		return AndExpression{nodes: combined}, nil
	}
	if s.operator == Or {
		return OrExpression{nodes: combined}, nil
	}
	return nil, ParserError{Message: fmt.Sprintf("unexpected operator: %v", s.operator)}
}

// Parser is used to parse an expression and return an AuthorizationExpression tree.
type Parser struct {
	Lexer *Lexer
}

func newParser(lexer *Lexer) *Parser {
	return &Parser{Lexer: lexer}
}

func (p *Parser) Parse() (AuthorizationExpression, error) {
	scopeStack := []*Scope{newScope()}

	for {
		tok, val, err := p.Lexer.nextToken()
		if err != nil {
			return nil, ParserError{Message: fmt.Sprintf("Lexer error: %v", err)}
		}

		if tok == None {
			break
		}

		currentScope := scopeStack[len(scopeStack)-1]

		switch tok {
		case AccessToken:
			currentScope.addLabel(val)
		case OpenParen:
			newScope := newScope()
			scopeStack = append(scopeStack, newScope)
		case And, Or:
			if err := currentScope.setOperator(tok); err != nil {
				return nil, err
			}
		case CloseParen:
			if len(scopeStack) == 1 {
				return nil, ParserError{Message: "unmatched closing parenthesis"}
			}
			finishedScope := scopeStack[len(scopeStack)-1]
			scopeStack = scopeStack[:len(scopeStack)-1]
			expression, err := finishedScope.Build()
			if err != nil {
				return nil, err
			}
			currentScope = scopeStack[len(scopeStack)-1]
			currentScope.addNode(expression)
		default:
			return nil, ParserError{Message: fmt.Sprintf("unexpected token: %v", tok)}
		}
	}

	if len(scopeStack) != 1 {
		return nil, ParserError{Message: "mismatched parentheses"}
	}

	return scopeStack[0].Build()
}
