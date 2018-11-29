package models

import (
	"fmt"
)

type Field struct {
	Name        string
	TypeOf      string
	Translation string
}

func (field *Field) ToString() string {
	return fmt.Sprintf("%s %s `json:\"%s\"`", "Field_"+field.Name, field.Translation, field.Name)
}
