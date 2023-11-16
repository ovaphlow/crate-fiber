package events

import (
	"fmt"
	"log"
	"ovaphlow/crate/internal/utilities"
	"strings"
)

var columns = []string{"id", "relation_id", "reference_id", "tags", "detail", "time"}

func Filter(relationId int64, referenceId int64, tags []string, detail string, timeRange []string, skip int64, take int) ([]Event, error) {
	q := fmt.Sprintf(`
	select %s, cast(id as char) _id
	from crate.events
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
	if detail != "" {
		conditions = append(conditions, "json_contains(detail, ?)")
		params = append(params, detail)
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
	log.Println(q)
	log.Println(params)
	p := make([]interface{}, len(params))
	for i := range params {
		p[i] = &params[i]
	}
	log.Println(p)
	rows, err := utilities.MySQL.Query("select 1")
	// rows, err := utilities.MySQL.Query(q)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer rows.Close()
	var result []Event
	for rows.Next() {
		var row Event
		err = rows.Scan(&row.Id, &row.RelationId, &row.ReferenceId, &row.Tags, &row.Detail, &row.Time, &row.IdStr)
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}
		result = append(result, row)
	}
	return result, nil
}
