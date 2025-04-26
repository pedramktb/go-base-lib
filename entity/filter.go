package entity

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/pedramktb/go-base-lib/taggederror"
)

var (
	ErrInvalidValue = taggederror.ErrBadRequest.Wrap(taggederror.New(
		errors.New("invalid value"),
		"INVALID_VALUE",
	))
	ErrUnknownFieldInFilter = taggederror.ErrBadRequest.Wrap(taggederror.New(
		errors.New("unknown field in filter expression"),
		"UNKNOWN_FIELD_IN_FILTER",
	))
	ErrInvalidFilterExpression = errors.New("invalid filter expression")
)

type Expression[E Entity] interface {
	ToSQLQuery() squirrel.Sqlizer
}

type LogicalOperator string

const (
	LogicalAnd LogicalOperator = "and" // logical and
	LogicalOr  LogicalOperator = "or"  // logical or
)

type LogicalExpression[E Entity] struct {
	Left Expression[E]
	LogicalOperator
	Right Expression[E]
}

func (e *LogicalExpression[E]) ToSQLQuery() squirrel.Sqlizer {
	switch e.LogicalOperator {
	case LogicalAnd:
		return squirrel.And([]squirrel.Sqlizer{e.Left.ToSQLQuery(), e.Right.ToSQLQuery()})
	case LogicalOr:
		return squirrel.Or([]squirrel.Sqlizer{e.Left.ToSQLQuery(), e.Right.ToSQLQuery()})
	}
	return nil
}

type ConditionOperator string

const (
	ConditionRGX ConditionOperator = "rgx" // regular expression
	ConditionEQ  ConditionOperator = "eq"  // equality
	ConditionNE  ConditionOperator = "ne"  // not equal
	ConditionIN  ConditionOperator = "in"  // in
	ConditionNIN ConditionOperator = "nin" // not in
	ConditionGT  ConditionOperator = "gt"  // greater than
	ConditionLT  ConditionOperator = "lt"  // less than
	ConditionGTE ConditionOperator = "gte" // greater than or equal to
	ConditionLTE ConditionOperator = "lte" // less than or equal to
)

type ConditionExpression[E Entity] struct {
	Field    string
	Operator ConditionOperator
	Value    any
}

func (e *ConditionExpression[E]) ToSQLQuery() squirrel.Sqlizer {
	switch e.Operator {
	case ConditionRGX:
		return squirrel.Expr("REGEXP_LIKE(?, ?)", string(e.Field), e.Value)
	case ConditionEQ:
		return squirrel.Eq{string(e.Field): e.Value}
	case ConditionNE:
		return squirrel.NotEq{string(e.Field): e.Value}
	case ConditionGT:
		return squirrel.Gt{string(e.Field): e.Value}
	case ConditionLT:
		return squirrel.Lt{string(e.Field): e.Value}
	case ConditionGTE:
		return squirrel.GtOrEq{string(e.Field): e.Value}
	case ConditionLTE:
		return squirrel.LtOrEq{string(e.Field): e.Value}
	}
	return nil
}

func ExpressionFromDTO[E Entity](dto string) (Expression[E], error) {
	if dto == "" {
		return nil, nil
	}
	var exp map[string]any
	err := json.Unmarshal([]byte(dto), &exp)
	if err != nil {
		return nil, err
	}
	return expressionFromMap[E](exp)
}

func expressionFromMap[E Entity](m map[string]any) (Expression[E], error) {
	if len(m) == 0 {
		return nil, nil
	} else if len(m) != 1 {
		return nil, ErrInvalidFilterExpression
	}
	for k, v := range m {
		switch k {
		case "$and", "$or": // OP: [EXP1, EXP2]
			if exps, ok := v.([]map[string]any); ok {
				var op LogicalOperator
				switch k {
				case "$and":
					op = LogicalAnd
				case "$or":
					op = LogicalOr
				}
				return logicalExpressionFromMap[E](op, exps...)
			}
			return nil, ErrInvalidFilterExpression
		default: // FIELD: {OP: VALUE}
			if opVal, ok := v.(map[string]any); ok && len(opVal) == 1 {
				return conditionExpressionFromMap[E](k, opVal)
			}
			return nil, ErrInvalidFilterExpression
		}
	}
	return nil, nil
}

func logicalExpressionFromMap[E Entity](op LogicalOperator, exps ...map[string]any) (Expression[E], error) {
	if len(exps) == 0 {
		return nil, nil
	}
	base, err := expressionFromMap[E](exps[0])
	if err != nil {
		return nil, err
	}
	for i := 1; i < len(exps); i++ {
		exp, err := expressionFromMap[E](exps[i])
		if err != nil {
			return nil, err
		}
		base = &LogicalExpression[E]{
			LogicalOperator: op,
			Left:            base,
			Right:           exp,
		}
	}
	return base, nil
}

func conditionExpressionFromMap[E Entity](field string, opVal map[string]any) (Expression[E], error) {
	for k, v := range opVal {
		var op ConditionOperator
		switch k {
		case "$eq":
			op = ConditionEQ
		case "$nq":
			op = ConditionNE
		case "$in":
			op = ConditionIN
		case "$nin":
			op = ConditionNIN
		case "$gt":
			op = ConditionGT
		case "$lt":
			op = ConditionLT
		case "$gte":
			op = ConditionGTE
		case "$lte":
			op = ConditionLTE
		}
		switch k {
		case "$eq", "$nq", "$gt", "$lt", "$gte", "$lte":
			value, err := parseConditionExpressionValue[E](field, v)
			if err != nil {
				return nil, err
			}
			return &ConditionExpression[E]{
				Field:    field,
				Operator: op,
				Value:    value,
			}, nil
		case "$in", "$nin":
			if vs, ok := v.([]any); ok {
				values := make([]any, 0, len(vs))
				for _, v := range vs {
					value, err := parseConditionExpressionValue[E](field, v)
					if err != nil {
						return nil, err
					}
					values = append(values, value)
				}
				return &ConditionExpression[E]{
					Field:    field,
					Operator: op,
					Value:    values,
				}, nil
			}
		}
	}
	return nil, nil
}

func parseConditionExpressionValue[E Entity](field string, value any) (any, error) {
	fieldVal, ok := (*new(E)).Fields()[field]
	if !ok {
		return nil, ErrUnknownField.Wrap(fmt.Errorf("unknown field %q in filter expression", field))
	}
	if withUnmarshal, ok := fieldVal.(json.Unmarshaler); ok {
		json, _ := json.Marshal(value) // Can't fail, we just unmarshaled it above
		if err := withUnmarshal.UnmarshalJSON(json); err != nil {
			return nil, ErrInvalidValue.Wrap(fmt.Errorf("invalid value %q for field %q: %w", json, field, err))
		}
		return any(withUnmarshal), nil
	}
	return value, nil
}
