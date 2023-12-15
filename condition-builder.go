package main

import (
	"fmt"
	"strings"
)

type ConditionBuilder struct {
	Conditions []string
	Params     []interface{}
}

func NewConditionBuilder(
	conditions []string,
	params []interface{},
) *ConditionBuilder {
	return &ConditionBuilder{
		Conditions: conditions,
		Params:     params,
	}
}

func (cb *ConditionBuilder) EqualBuilder(equal []string) {
	for i := 0; i < len(equal); i += 2 {
		cb.Conditions = append(cb.Conditions, fmt.Sprintf("%s = ?", equal[i]))
		cb.Params = append(cb.Params, equal[i+1])
	}
}

func (cb *ConditionBuilder) ObjectContainBuilder(objectContain []string) {
	for i := 0; i < len(objectContain); i += 3 {
		cb.Conditions = append(
			cb.Conditions,
			fmt.Sprintf(
				"json_contains(%s, json_object('%s', ?))",
				objectContain[i],
				objectContain[i+1],
			),
		)
		cb.Params = append(cb.Params, objectContain[i+2])
	}
}

func (cb *ConditionBuilder) ArrayContainBuilder(arrayContain []string) {
	for i := 0; i < len(arrayContain); i += 2 {
		cb.Conditions = append(
			cb.Conditions,
			fmt.Sprintf("json_contains(%s, json_array(?))", arrayContain[i]),
		)
		cb.Params = append(cb.Params, arrayContain[i+1])
	}
}

func (cb *ConditionBuilder) LikeBuilder(like []string) {
	for i := 0; i < len(like); i += 2 {
		cb.Conditions = append(
			cb.Conditions,
			fmt.Sprintf("position(? in %s)", like[i]),
		)
		cb.Params = append(cb.Params, like[i+1])
	}
}

func (cb *ConditionBuilder) ObjectLikeBuilder(objectLike []string) {
	for i := 0; i < len(objectLike); i += 3 {
		cb.Conditions = append(
			cb.Conditions,
			fmt.Sprintf(
				"position(? in %s->>'$.%s')",
				objectLike[i],
				objectLike[i+1],
			),
		)
		cb.Params = append(cb.Params, objectLike[i+2])
	}
}

func (cb *ConditionBuilder) InBuilder(in []string) {
	c := make([]string, len(in)-1)
	for i := range c {
		c[i] = "?"
	}
	cb.Conditions = append(
		cb.Conditions,
		fmt.Sprintf("%s in (%s)", in[0], strings.Join(c, ", ")),
	)
	params := make([]interface{}, len(in)-1)
	for i := range params {
		params[i] = in[i+1]
	}
	cb.Params = append(cb.Params, params...)
}

func (cb *ConditionBuilder) LesserBuilder(lesser []string) {
	for i := 0; i < len(lesser); i += 2 {
		cb.Conditions = append(
			cb.Conditions,
			fmt.Sprintf("%s <= ?", lesser[i]),
		)
		cb.Params = append(cb.Params, lesser[i+1])
	}
}

func (cb *ConditionBuilder) GreaterBuilder(greater []string) {
	for i := 0; i < len(greater); i += 2 {
		cb.Conditions = append(
			cb.Conditions,
			fmt.Sprintf("%s >= ?", greater[i]),
		)
		cb.Params = append(cb.Params, greater[i+1])
	}
}

func (cb *ConditionBuilder) Build() (string, []interface{}) {
	if len(cb.Conditions) == 0 {
		return "", []interface{}{}
	}
	return fmt.Sprintf(
		"where %s",
		strings.Join(cb.Conditions, " and "),
	), cb.Params
}
