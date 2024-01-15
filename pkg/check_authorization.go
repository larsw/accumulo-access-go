// Package pkg Copyright 2024 Lars Wilhelmsen <sral-backwards@sral.org>. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in the LICENSE file.
package pkg

import "strings"

// CheckAuthorization checks if the given authorizations are allowed to perform the given expression.
// Arguments:
//
//	expression: The expression to check.
//	authorizations: A comma-separated list of authorizations.
//
// Returns:
//
//	True if the authorizations are allowed to perform the expression, false otherwise.
func CheckAuthorization(expression string, authorizations string) (bool, error) {
	authorizationMap := CommaSeparatedStringToMap(authorizations)
	return CheckAuthorizationByMap(expression, authorizationMap)
}

func CheckAuthorizationByMap(expression string, authorizations map[string]bool) (bool, error) {
	parser := newParser(newLexer(expression))
	ast, err := parser.Parse()
	if err != nil {
		return false, err
	}
	return ast.evaluate(authorizations), nil
}

func CommaSeparatedStringToMap(authorizations string) map[string]bool {
	// should also trim quotes from quoted strings
	authorizationMap := make(map[string]bool)
	for _, authorization := range strings.Split(authorizations, ",") {
		tmp := strings.TrimSpace(authorization)
		if tmp[0] == '"' {
			tmp = tmp[1:]
		}
		if tmp[len(tmp)-1] == '"' {
			tmp = tmp[:len(tmp)-1]
		}
		authorizationMap[tmp] = true
	}
	return authorizationMap
}

// PrepareAuthorizationCheck returns a function that can be used to check if the given authorizations are allowed to perform the given expression.
// Arguments:
//
//	authorizations: A comma-separated list of authorizations.
//
// Returns:
//
//	A function that can be used to check if the given authorizations are allowed to perform the given expression.
func PrepareAuthorizationCheck(authorizations string) func(string) (bool, error) {
	authorizationMap := CommaSeparatedStringToMap(authorizations)
	return func(expression string) (bool, error) {
		parser := newParser(newLexer(expression))
		ast, err := parser.Parse()
		if err != nil {
			return false, err
		}
		return ast.evaluate(authorizationMap), nil
	}
}
