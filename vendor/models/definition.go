package models

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

type Definition struct {
	Fields map[string]Field
}

func (definition *Definition) Signature() string {
	var keys []string
	for key, field := range definition.Fields {
		keys = append(keys, key+":"+field.TypeOf)
	}

	sort.Strings(keys)

	return strings.Join(keys, ",")
}

func (definition *Definition) Name() string {
	hasher := sha256.New()
	hasher.Write([]byte(definition.Signature()))
	return strings.Title(hex.EncodeToString(hasher.Sum(nil)))
}

func (definition *Definition) Body() string {
	ret := ""

	for key, field := range definition.Fields {
		ret += fmt.Sprintf("\n%s %s", strings.Title(key), field.ToString())
	}

	return ret
}

func (definition *Definition) ToString() string {
	return fmt.Sprintf("type %s struct {%s\n}\n", definition.Name(), definition.Body())
}
