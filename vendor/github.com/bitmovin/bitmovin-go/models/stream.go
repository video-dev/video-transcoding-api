package models

import "github.com/bitmovin/bitmovin-go/bitmovintypes"

type Stream struct {
	ID                   *string                         `json:"id,omitempty"`
	Name                 *string                         `json:"name,omitempty"`
	Description          *string                         `json:"description,omitempty"`
	CustomData           map[string]interface{}          `json:"customData,omitempty"`
	InputStreams         []InputStream                   `json:"inputStreams,omitempty"`
	CodecConfigurationID *string                         `json:"codecConfigId,omitempty"`
	Outputs              []Output                        `json:"outputs,omitempty"`
	DecodingErrorMode    bitmovintypes.DecodingErrorMode `json:"decodingErrorMode,omitempty"`
	Conditions           *StreamCondition                `json:"conditions,omitempty"`
	Mode                 bitmovintypes.StreamMode        `json:"mode,omitempty"`
}

type StreamCondition struct {
	Attribute  bitmovintypes.ConditionAttribute `json:"attribute,omitempty"`
	Operator   string                           `json:"operator,omitempty"`
	Value      string                           `json:"value,omitempty"`
	Type       bitmovintypes.ConditionType      `json:"type"`
	Conditions []*StreamCondition               `json:"conditions,omitempty"`
}

// NewAttributeCondition creates a Condition that tests an attribute against a value given the operator
func NewAttributeCondition(attribute bitmovintypes.ConditionAttribute, operator, value string) *StreamCondition {
	return &StreamCondition{
		Attribute: attribute,
		Type:      bitmovintypes.ConditionTypeCondition,
		Operator:  operator,
		Value:     value,
	}
}

// NewAndConjunction creates a logical Conjunction (AND) of all the condition parameters
func NewAndConjunction(conditions ...*StreamCondition) *StreamCondition {
	return &StreamCondition{
		Type:       bitmovintypes.ConditionTypeAnd,
		Conditions: conditions,
	}
}

// NewOrDisjunction creates a logical Disjunction (OR) of all the condition parameters
func NewOrDisjunction(conditions ...*StreamCondition) *StreamCondition {
	return &StreamCondition{
		Type:       bitmovintypes.ConditionTypeOr,
		Conditions: conditions,
	}
}
