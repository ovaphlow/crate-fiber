package main

import (
	"errors"
	"strings"
)

func StringToMap(s string) (map[string]interface{}, error) {
	items := strings.Split(s, ",")
	if len(items)%2 != 0 {
		return nil, errors.New("输入的字符串不能转换为键值对")
	}

	result := make(map[string]interface{})
	for i := 0; i < len(items); i += 2 {
		result[items[i]] = items[i+1]
	}

	return result, nil
}
