// Package pkg Copyright 2024 Lars Wilhelmsen <sral-backwards@sral.org>. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.
package pkg

import (
	"fmt"
	"unicode"
)

// Token represents the different types of tokens.
type Token int

// Define token constants.
const (
	AccessToken Token = iota
	OpenParen
	CloseParen
	And
	Or
	None
	Error
)

//func (t value) String() string {
//	switch t {
//	case AccessToken:
//		return "AccessToken"
//	case OpenParen:
//		return "("
//	case CloseParen:
//		return ")"
//	case And:
//		return "&"
//	case Or:
//		return "|"
//	case None:
//		return "None"
//	default:
//		return "Unknown"
//	}
//}

// Lexer represents a lexer for tokenizing strings.
type Lexer struct {
	input   string
	pos     int
	readPos int
	ch      byte
}

// newLexer creates a new Lexer instance.
func newLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// LexerError represents an error that occurred during lexing.
type LexerError struct {
	Char     byte
	Position int
}

func (e LexerError) Error() string {
	return fmt.Sprintf("unexpected character '%c' at position %d", e.Char, e.Position)
}

func (l *Lexer) readChar() {
	if l.readPos >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPos]
	}
	l.pos = l.readPos
	l.readPos++
}

func (l *Lexer) peekChar() byte {
	if l.readPos >= len(l.input) {
		return 0
	}
	return l.input[l.readPos]
}

func (l *Lexer) nextToken() (Token, string, error) {
	l.skipWhitespace()

	switch l.ch {
	case '(':
		l.readChar()
		return OpenParen, "", nil
	case ')':
		l.readChar()
		return CloseParen, "", nil
	case '&':
		l.readChar()
		return And, "", nil
	case '|':
		l.readChar()
		return Or, "", nil
	case 0:
		return None, "", nil
	case '"':
		strLiteral := l.readString()
		return AccessToken, strLiteral, nil
	default:
		if isLegalTokenLetter(l.ch) {
			val := l.readIdentifier()
			return AccessToken, val, nil
		} else {
			return Error, "", LexerError{Char: l.ch, Position: l.pos}
		}
	}
}

func (l *Lexer) readIdentifier() string {
	startPos := l.pos
	for isLegalTokenLetter(l.ch) {
		l.readChar()
	}
	return l.input[startPos:l.pos]
}

func (l *Lexer) readString() string {
	startPos := l.pos + 1 // Skip initial double quote
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break // End of string or end of input
		}

		// Handle escape sequences
		if l.ch == '\\' {
			l.readChar()
			if l.ch != '"' && l.ch != '\\' {
				// Handle invalid escape sequence
				return l.input[startPos : l.pos-1] // Return string up to invalid escape
			}
		}
	}
	str := l.input[startPos:l.pos]
	l.readChar() // Skip closing double quote
	return str
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func isLegalTokenLetter(ch byte) bool {

	return unicode.IsLetter(rune(ch)) ||
		unicode.IsDigit(rune(ch)) ||
		ch == '_' ||
		ch == '-' ||
		ch == '.' ||
		ch == ':'
}
