package interpreter

import "testing"

func TestInterpreter(t *testing.T) {
	a := NewVariable("is_admin")
	b := NewVariable("has_token")
	c := NewVariable("is_banned")

	rule := And(a, And(b, Not(c)))
	ctx := NewContext(map[string]bool{
		"is_admin":  true,
		"has_token": true,
		"is_banned": false,
	})

	if !rule.Interpret(ctx) {
		t.Fatal("expected rule to pass")
	}
}

func TestOrExpression(t *testing.T) {
	a := NewVariable("paid")
	b := NewVariable("trial")
	rule := Or(a, b)

	ctx := NewContext(map[string]bool{
		"paid":  false,
		"trial": true,
	})

	if !rule.Interpret(ctx) {
		t.Fatal("expected or rule to pass")
	}
}

func TestEvaluateRule(t *testing.T) {
	rule := Not(NewVariable("expired"))
	ctx := NewContext(map[string]bool{"expired": false})

	if EvaluateRule(rule, ctx) != "result=true" {
		t.Fatalf("unexpected result: %s", EvaluateRule(rule, ctx))
	}
}
