package models

import (
	"testing"

	"github.com/bitmovin/bitmovin-go/bitmovintypes"
)

func TestConditionBuilder(t *testing.T) {
	condition := NewAttributeCondition(bitmovintypes.ConditionAttributeFPS, "==", "25")
	testCondition(condition, t)
}

func testCondition(condition *StreamCondition, t *testing.T) {
	if condition.Type != bitmovintypes.ConditionTypeCondition {
		t.Errorf("Wanted ConditionType Condition got %v", condition.Type)
	}
	if condition.Attribute != bitmovintypes.ConditionAttributeFPS {
		t.Errorf("Wanted Attribute FPS got %v", condition.Attribute)
	}
	if condition.Operator != "==" {
		t.Errorf("Wanted Value == got %v", condition.Operator)
	}
	if condition.Value != "25" {
		t.Errorf("Wanted Value 25 got %v", condition.Value)
	}
}

func TestNewAndConjunction(t *testing.T) {
	cond := NewAndConjunction(
		NewAttributeCondition(bitmovintypes.ConditionAttributeFPS, "==", "25"),
	)

	if cond.Type != bitmovintypes.ConditionTypeAnd {
		t.Errorf("Wanted AndConjunction Type to be AND got %v", cond.Type)
	}

	if len(cond.Conditions) != 1 {
		t.Fatalf("Wanted 1 Condition, got %d", len(cond.Conditions))
	}
	condition := cond.Conditions[0]
	testCondition(condition, t)
}

func TestNewOrDisjunction(t *testing.T) {
	cond := NewOrDisjunction(
		NewAttributeCondition(bitmovintypes.ConditionAttributeFPS, "==", "25"),
	)
	if cond.Type != bitmovintypes.ConditionTypeOr {
		t.Errorf("Wanted AndConjunction Type to be OR got %v", cond.Type)
	}

	if len(cond.Conditions) != 1 {
		t.Fatalf("Wanted 1 Condition, got %d", len(cond.Conditions))
	}
	condition := cond.Conditions[0]
	testCondition(condition, t)
}

func buildNestedCondition() *StreamCondition {
	return NewOrDisjunction(
		NewAndConjunction(
			NewAttributeCondition(bitmovintypes.ConditionAttributeFPS, "==", "25"),
			NewAttributeCondition(bitmovintypes.ConditionAttributeBitrate, "==", "14000"),
		),
		NewAndConjunction(
			NewAttributeCondition(bitmovintypes.ConditionAttributeFPS, "==", "60"),
			NewAttributeCondition(bitmovintypes.ConditionAttributeBitrate, "==", "7000"),
		),
	)
}

func TestNestedConditions(t *testing.T) {
	cond := buildNestedCondition()
	if cond.Type != bitmovintypes.ConditionTypeOr {
		t.Errorf("Wanted ConditionType OR got %v", cond.Type)
	}
	if len(cond.Conditions) != 2 {
		t.Fatalf("Wanted 2 Conditions, got %d", len(cond.Conditions))
	}
	firstCond := cond.Conditions[0]
	if firstCond.Type != bitmovintypes.ConditionTypeAnd {
		t.Errorf("Expected ConditionType of first nested Condition to be AND got %v", firstCond.Type)
	}
	secondCond := cond.Conditions[1]
	if secondCond.Type != bitmovintypes.ConditionTypeAnd {
		t.Errorf("Expected ConditionType of second nested Condition to be AND got %v", secondCond.Type)
	}
}
