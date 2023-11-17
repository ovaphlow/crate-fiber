package main

func SliceToMap(array []string) map[string]string {
	var mapResult = make(map[string]string)
	for i := 0; i < len(array); i += 2 {
		mapResult[array[i]] = array[i+1]
	}
	return mapResult
}
