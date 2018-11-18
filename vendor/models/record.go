package models

type Record struct {
	Name   string
	Value  string
	TypeOf string

	Records []Record
}
