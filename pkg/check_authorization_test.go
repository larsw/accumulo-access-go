package pkg

import (
	"fmt"
	"testing"
)

type testCase struct {
	expression     string
	authorizations string
	expected       bool
}

func TestCheckAuthorization(t *testing.T) {
	testCases := []testCase{
		{"label1", "label1", true},
		{"label1|label2", "label1", true},
		{"label1&label2", "label1", false},
		{"label1&label2", "label1,label2", true},
		{"label1&(label2 | label3)", "label1", false},
		{"label1&(label2 | label3)", "label1,label3", true},
		{"label1&(label2 | label3)", "label1,label2", true},
		{"(label2 | label3)", "label1", false},
		{"(label2 | label3)", "label2", true},
		{"(label2 & label3)", "label2", false},
		{"((label2 | label3))", "label2", true},
		{"((label2 & label3))", "label2", false},
		{"(((((label2 & label3)))))", "label2", false},
		{"(a & b) & (c & d)", "a,b,c,d", true},
		{"(a & b) & (c & d)", "a,b,c", false},
		{"(a & b) | (c & d)", "a,b,d", true},
		{"(a | b) & (c | d)", "a,d", true},
		{"\"a b c\"", "\"a b c\"", true},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("\"%v\" + \"%v\" -> %v", tc.expression, tc.authorizations, tc.expected), func(t *testing.T) {
			result, err := CheckAuthorization(tc.expression, tc.authorizations)
			if err != nil {
				t.Fatal(err)
			}
			if result != tc.expected {
				t.Fatalf("expected %v for %s with %s", tc.expected, tc.expression, tc.authorizations)
			}
		})
	}
}
