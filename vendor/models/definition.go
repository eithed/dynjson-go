package models

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

type Definition struct {
	Fields []Field
}

func (definition *Definition) Signature() string {
	var keys []string
	for _, field := range definition.Fields {
		keys = append(keys, field.Name+":"+field.TypeOf)
	}

	sort.Strings(keys)

	return strings.Join(keys, ",")
}

func (definition *Definition) Name() string {
	hasher := sha256.New()
	hasher.Write([]byte(definition.Signature()))
	return "Struct_" + hex.EncodeToString(hasher.Sum(nil))
}

func (definition *Definition) Body() string {
	ret := ""

	for _, field := range definition.Fields {
		ret += fmt.Sprintf("%s\n", field.ToString())
	}

	return ret
}

func (definition *Definition) Accessors() string {
	ret := ""

	for _, field := range definition.Fields {
		ret += fmt.Sprintf(`func (str *%s) %s() %s {
			return str.%s
		}
		`, definition.Name(), strings.Title(field.Name), field.Translation, "Field_"+field.Name)
	}

	return ret
}

func (definition *Definition) ToString() string {
	return fmt.Sprintf(`type %s struct {
		%s}
	
	%s`, definition.Name(), definition.Body(), definition.Accessors())
}
