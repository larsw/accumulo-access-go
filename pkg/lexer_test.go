package pkg

import "testing"

func TestLexer(t *testing.T) {
	lexer := newLexer("(a & b) | c")
	token, _, _ := lexer.nextToken()
	if token != OpenParen {
		t.Fatal("expected OpenParen")
	}
	token, val, _ := lexer.nextToken()
	if token != AccessToken {
		t.Fatal("expected AccessToken")
	}
	if val != "a" {
		t.Fatal("expected a")
	}
	token, _, _ = lexer.nextToken()
	if token != And {
		t.Fatal("expected And")
	}
	token, val, _ = lexer.nextToken()
	if token != AccessToken {
		t.Fatal("expected AccessToken")
	}
	if val != "b" {
		t.Fatal("expected b")
	}
	token, _, _ = lexer.nextToken()
	if token != CloseParen {
		t.Fatal("expected CloseParen")
	}
	token, _, _ = lexer.nextToken()
	if token != Or {
		t.Fatal("expected Or")
	}
	token, val, _ = lexer.nextToken()
	if token != AccessToken {
		t.Fatal("expected AccessToken")
	}
	if val != "c" {
		t.Fatal("expected c")
	}
	token, _, _ = lexer.nextToken()
	if token != None {
		t.Fatal("expected end of input")
	}
}
