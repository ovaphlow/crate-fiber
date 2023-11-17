package main

type Event struct {
	Id          int64  `json:"id"`
	RelationId  int64  `json:"relationId"`
	ReferenceId int64  `json:"referenceId"`
	Tags        string `json:"tags"`
	Detail      string `json:"detail"`
	Time        string `json:"time"`
	Id_         string `json:"_id"`
}
