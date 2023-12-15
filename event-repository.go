package main

import (
	"fmt"
	"log"
	"strings"
)

var columns = []string{"id", "relation_id", "reference_id", "tags", "detail", "time"}

func EventDefaultFilter(
	skip int64,
	take int,
	equal []string,
	objectContain []string,
	arrayContain []string,
	like []string,
	objectLike []string,
	in []string,
	lesser []string,
	greater []string,
) ([]Event, error) {
	q := fmt.Sprintf(`select %s from events`, strings.Join(columns, ", "))
	var conditions []string
	var params []interface{}
	conditionBuilder := NewConditionBuilder(conditions, params)
	if len(equal) > 0 && len(equal)%2 == 0 {
		conditionBuilder.EqualBuilder(equal)
	}
	if len(objectContain) > 0 && len(objectContain)%3 == 0 {
		conditionBuilder.ObjectContainBuilder(objectContain)
	}
	if len(arrayContain) > 0 && len(arrayContain)%2 == 0 {
		conditionBuilder.ArrayContainBuilder(arrayContain)
	}
	if len(like) > 0 && len(like)%2 == 0 {
		conditionBuilder.LikeBuilder(like)
	}
	if len(objectLike) > 0 && len(objectLike)%3 == 0 {
		conditionBuilder.ObjectLikeBuilder(objectLike)
	}
	if len(in) >= 2 {
		conditionBuilder.InBuilder(in)
	}
	if len(lesser) > 0 && len(lesser)%2 == 0 {
		conditionBuilder.LesserBuilder(lesser)
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
	// slogger.Info(q)
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
	var result []Event
	for rows.Next() {
		var row Event
		err = rows.Scan(&row.Id, &row.RelationId, &row.ReferenceId, &row.Tags, &row.Detail, &row.Time)
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}
		result = append(result, row)
	}
	if len(result) == 0 {
		return []Event{}, nil
	}
	return result, nil
}
