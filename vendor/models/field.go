package models

type Field struct {
	TypeOf string
}

func (field *Field) ToString() string {
	return field.TypeOf
}
