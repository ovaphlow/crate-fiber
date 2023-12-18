package main

import (
	"fmt"
	"log"
	"strings"
	"time"
)

var Schemas = map[string]map[string]interface{}{
	"event": {
		"table":   "events",
		"columns": []string{"id", "relation_id", "reference_id", "tags", "detail", "time"},
	},
}

func Retrieve(
	name string,
	skip int64,
	take int,
	option map[string]interface{},
) ([]map[string]interface{}, error) {
	schema, ok := Schemas[name]
	if !ok {
		return nil, fmt.Errorf("schema not found")
	}
	q := fmt.Sprintf(
		`select %s from %s`,
		// strings.Join(schema["columns"], ", "),
		strings.Join(schema["columns"].([]string), ", "),
		schema["table"],
	)
	var conditions []string
	var params []interface{}
	conditionBuilder := NewConditionBuilder(conditions, params)
	equal, ok := option["equal"].([]string)
	if !ok {
		equal = []string{}
	}
	if len(equal) > 0 && len(equal)%2 == 0 {
		conditionBuilder.EqualBuilder(equal)
	}
	objectContain, ok := option["objectContain"].([]string)
	if !ok {
		objectContain = []string{}
	}
	if len(objectContain) > 0 && len(objectContain)%3 == 0 {
		conditionBuilder.ObjectContainBuilder(objectContain)
	}
	arrayContain, ok := option["arrayContain"].([]string)
	if !ok {
		arrayContain = []string{}
	}
	if len(arrayContain) > 0 && len(arrayContain)%2 == 0 {
		conditionBuilder.ArrayContainBuilder(arrayContain)
	}
	like, ok := option["like"].([]string)
	if !ok {
		like = []string{}
	}
	if len(like) > 0 && len(like)%2 == 0 {
		conditionBuilder.LikeBuilder(like)
	}
	objectLike, ok := option["objectLike"].([]string)
	if !ok {
		objectLike = []string{}
	}
	if len(objectLike) > 0 && len(objectLike)%3 == 0 {
		conditionBuilder.ObjectLikeBuilder(objectLike)
	}
	in, ok := option["in"].([]string)
	if !ok {
		in = []string{}
	}
	if len(in) >= 2 {
		conditionBuilder.InBuilder(in)
	}
	lesser, ok := option["lesser"].([]string)
	if !ok {
		lesser = []string{}
	}
	if len(lesser) > 0 && len(lesser)%2 == 0 {
		conditionBuilder.LesserBuilder(lesser)
	}
	greater, ok := option["greater"].([]string)
	if !ok {
		greater = []string{}
	}
	if len(greater) > 0 && len(greater)%2 == 0 {
		conditionBuilder.GreaterBuilder(greater)
	}
	if len(conditionBuilder.Conditions) > 0 {
		var where string
		where, params = conditionBuilder.Build()
		q = fmt.Sprintf(`%s %s`, q, where)
	}
	q = fmt.Sprintf(`%s order by id desc limit %d, %d`, q, skip, take)
	statement, err := MySQL.Prepare(q)
	if err != nil {
		return nil, err
	}
	defer statement.Close()
	rows, err := statement.Query(params...)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer rows.Close()
	var result []map[string]interface{}
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		for i := range values {
			values[i] = new(interface{})
		}
		err = rows.Scan(values...)
		if err != nil {
			return nil, err
		}
		row := make(map[string]interface{})
		for i, column := range columns {
			val := *(values[i].(*interface{}))
			if b, ok := val.([]byte); ok {
				row[column] = string(b)
			} else if t, ok := val.(time.Time); ok {
				row[column] = t.Format("2006-01-02 15:04:05")
			} else {
				row[column] = val
			}
		}
		result = append(result, row)
	}
	if len(result) == 0 {
		return []map[string]interface{}{}, nil
	}
	return result, nil
}
