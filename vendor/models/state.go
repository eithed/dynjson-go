package models

import (
	"common"
	"fmt"
	"reflect"
)

type State struct {
	Definitions []Definition
	Root        *Record
}

func (state *State) TraverseStructure(record interface{}) *Record {
	switch record.(type) {
	// if the record is an array
	case []interface{}:
		ret := Record{
			TypeOf: "array"}

		// iterate through values that are stored within the array
		for _, value := range record.([]interface{}) {
			switch value.(type) {

			// value is an object
			case map[string]interface{}:
				record := state.TraverseStructure(value)

				ret.Records = append(ret.Records, *record)

			// value is an array
			case []interface{}:
				record := state.TraverseStructure(value)

				ret.Records = append(ret.Records, *record)

			// value is a primitive
			default:
				record := Record{
					TypeOf: reflect.TypeOf(value).Name(),
					Value:  common.Itos(value)}

				ret.Records = append(ret.Records, record)
			}
		}

		return &ret
	// if the record is an object
	case map[string]interface{}:
		ret := Record{
			TypeOf: "object"}

		// iterate through values that are stored within the object
		for name, value := range record.(map[string]interface{}) {
			switch value.(type) {

			// value is an object
			case map[string]interface{}:
				record := state.TraverseStructure(value)
				record.Name = name

				ret.Records = append(ret.Records, *record)

			// value is an array
			case []interface{}:
				record := state.TraverseStructure(value)
				record.Name = name

				ret.Records = append(ret.Records, *record)

			// value is a primitive
			default:
				record := Record{
					Name:   name,
					Value:  common.Itos(value),
					TypeOf: reflect.TypeOf(value).Name()}

				ret.Records = append(ret.Records, record)
			}
		}

		return &ret
	}

	return nil
}

// @todo - cater for mixed types within array
func (state *State) RecurseKeyName(record *Record) *string {
	if record.TypeOf != "array" {
		return nil
	}

	switch record.Records[0].TypeOf {
	case "array":
		value := fmt.Sprintf("[]%s", *state.RecurseKeyName(&record.Records[0]))

		return &value
	case "object":
		value := fmt.Sprintf("[]%s", state.Definition(&record.Records[0]).Name())

		return &value
	default:
		value := fmt.Sprintf("[]%s", record.Records[0].TypeOf)

		return &value
	}

	return nil
}

func (state *State) Definition(record *Record) *Definition {
	if record.TypeOf != "object" {
		return nil
	}

	definition := Definition{}

	for _, value := range record.Records {
		field := Field{
			Name:   value.Name,
			TypeOf: value.TypeOf}

		switch value.TypeOf {
		// array of somethings
		case "array":
			field.Translation = *state.RecurseKeyName(&value)
		// object of somethings
		case "object":
			field.Translation = fmt.Sprintf("%s", state.Definition(&value).Name())
		default:
			field.Translation = value.TypeOf
		}

		definition.Fields = append(definition.Fields, field)
	}

	existingDefinition := state.ExistingDefinition(definition)

	if existingDefinition != nil {
		return existingDefinition
	}

	return &definition
}

func (state *State) ExistingDefinition(definition Definition) *Definition {
	for _, def := range state.Definitions {
		if definition.Signature() == def.Signature() {
			return &definition
		}
	}

	return nil
}

func (state *State) TraverseDefinitions(record *Record, parent *Record) {
	switch record.TypeOf {
	case "object":
		definition := state.Definition(record)

		for _, value := range record.Records {
			if value.TypeOf == "array" || value.TypeOf == "object" {
				state.TraverseDefinitions(&value, record)
			}
		}

		existingDefinition := state.ExistingDefinition(*definition)

		if existingDefinition == nil {
			state.Definitions = append(state.Definitions, *definition)
		}
	case "array":
		for _, value := range record.Records {
			state.TraverseDefinitions(&value, record)
		}
	}
}

func (state *State) SetRoot(record interface{}) {
	state.Root = state.TraverseStructure(record)
}

func (state *State) SetDefinitions() {
	state.TraverseDefinitions(state.Root, nil)
}

func (state *State) ParseDefinitions() string {
	ret := ""
	for _, definition := range state.Definitions {
		ret += "\n" + definition.ToString()
	}

	return ret
}
