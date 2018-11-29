package dynjson

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"models"
	"time"
)

func Parse(json interface{}) string {
	state := models.State{}

	state.SetRoot(json)
	state.SetDefinitions()

	var definitions string

	for _, definition := range state.Definitions {
		definitions += definition.ToString()
	}

	var definition string

	if state.Root.TypeOf == "object" {
		definition = state.Definition(state.Root).Name()
	} else {
		definition = *state.RecurseKeyName(state.Root)
	}

	return fmt.Sprintf(`package main
	
	import (
		"encoding/json"
	)

	%s
	
	var Root %s
	func jsonDecode(input []byte) {
		json.Unmarshal(input, &Root)
	}`, definitions, definition)
}

func ParseString(str string) {
	var i interface{}

	err := json.Unmarshal([]byte(str), &i)

	if err != nil {
		panic(err)
	}

	content := Parse(i)
	filename := fmt.Sprintf("./evaluated-%d.go", time.Now().UnixNano())

	ioutil.WriteFile(filename, []byte(content), 0644)
}

func ParseFile(name string) {
	content, err := ioutil.ReadFile(name)

	if err != nil {
		panic(err)
	}

	ParseString(string(content))
}
