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
	in []string,
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
	if len(in) >= 2 {
		conditionBuilder.InBuilder(in)
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

func EventsFilter(relationId int64, referenceId int64, tags []string, detail map[string]interface{}, timeRange []string, skip int64, take int) ([]Event, error) {
	q := fmt.Sprintf(`
	select %s, cast(id as char) _id
	from events
	`, strings.Join(columns, ","))
	var conditions []string
	var params []interface{}
	if relationId > 0 {
		conditions = append(conditions, "relation_id = ?")
		params = append(params, relationId)
	}
	if referenceId > 0 {
		conditions = append(conditions, "reference_id = ?")
		params = append(params, referenceId)
	}
	for _, tag := range tags {
		conditions = append(conditions, "json_contains(tags, json_array(?))")
		params = append(params, tag)
	}
	for k, v := range detail {
		conditions = append(conditions, "json_contains(detail, json_object(?, ?))")
		params = append(params, k, v)
	}
	if len(timeRange) == 2 {
		conditions = append(conditions, "time > ?", "time < ?")
		params = append(params, timeRange[0], timeRange[1])
	}
	if len(conditions) > 0 {
		q = fmt.Sprintf(`
		%s
		where %s
		`, q, strings.Join(conditions, " and "))
	}
	q = fmt.Sprintf(`
	%s
	order by id desc
	limit %d, %d
	`, q, skip, take)
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
