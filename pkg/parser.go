// Package pkg Copyright 2024 Lars Wilhelmsen <sral-backwards@sral.org>. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.
package pkg

import "fmt"

// Assuming Token and Lexer are already defined as in the previous lexer translation

// ParserError represents errors that can occur during parsing.
type ParserError struct {
	Message string
}

func (e ParserError) Error() string {
	return e.Message
}

// AuthorizationExpression is an interface for different expression types.
type AuthorizationExpression interface {
	Evaluate(authorizations map[string]bool) bool
}

// AndExpression represents an AND expression.
type AndExpression struct {
	Nodes []AuthorizationExpression
}

func (a AndExpression) Evaluate(authorizations map[string]bool) bool {
	for _, node := range a.Nodes {
		if !node.Evaluate(authorizations) {
			return false
		}
	}
	return true
}

// Implement Evaluate, ToJSONStr, and Normalize for AndExpression...

// OrExpression represents an OR expression.
type OrExpression struct {
	Nodes []AuthorizationExpression
}

func (o OrExpression) Evaluate(authorizations map[string]bool) bool {
	for _, node := range o.Nodes {
		if node.Evaluate(authorizations) {
			return true
		}
	}
	return false
}

// Implement Evaluate, ToJSONStr, and Normalize for OrExpression...

// AccessTokenExpression represents an access token.
type AccessTokenExpression struct {
	Token string
}

func (a AccessTokenExpression) Evaluate(authorizations map[string]bool) bool {
	return authorizations[a.Token]
}

// Implement Evaluate, ToJSONStr, and Normalize for AccessTokenExpression...

// Scope is used during parsing to build up expressions.
type Scope struct {
	Nodes    []AuthorizationExpression
	Labels   []AccessTokenExpression
	Operator Token
}

func newScope() *Scope {
	return &Scope{
		Nodes:    make([]AuthorizationExpression, 0),
		Labels:   make([]AccessTokenExpression, 0),
		Operator: None, // Assuming None is a defined Token value
	}
}

func (s *Scope) addNode(node AuthorizationExpression) {
	s.Nodes = append(s.Nodes, node)
}

func (s *Scope) addLabel(label string) {
	s.Labels = append(s.Labels, AccessTokenExpression{Token: label})
}

func (s *Scope) setOperator(operator Token) error {
	if s.Operator != None {
		return ParserError{Message: "unexpected operator"}
	}
	s.Operator = operator
	return nil
}

func (s *Scope) Build() (AuthorizationExpression, error) {
	if len(s.Labels) == 1 && len(s.Nodes) == 0 {
		return s.Labels[0], nil
	}

	if len(s.Nodes) == 1 && len(s.Labels) == 0 {
		return s.Nodes[0], nil
	}

	if s.Operator == None {
		return nil, ParserError{Message: "missing operator"}
	}
	// combine nodes and labels into one slice
	combined := make([]AuthorizationExpression, 0, len(s.Nodes)+len(s.Labels))
	for _, node := range s.Nodes {
		combined = append(combined, node)
	}
	for _, label := range s.Labels {
		combined = append(combined, label)
	}
	if s.Operator == And {
		return AndExpression{Nodes: combined}, nil
	}
	if s.Operator == Or {
		return OrExpression{Nodes: combined}, nil
	}
	return nil, ParserError{Message: fmt.Sprintf("unexpected operator: %v", s.Operator)}
}

// Parser is used to parse an expression and return an AuthorizationExpression tree.
type Parser struct {
	Lexer *Lexer
}

func NewParser(lexer *Lexer) *Parser {
	return &Parser{Lexer: lexer}
}

func (p *Parser) Parse() (AuthorizationExpression, error) {
	scopeStack := []*Scope{newScope()}

	for {
		tok, val, err := p.Lexer.nextToken()
		if err != nil {
			return nil, ParserError{Message: fmt.Sprintf("Lexer error: %v", err)}
		}

		if tok == None { // Assuming 0 represents the end of input
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
