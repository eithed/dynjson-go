package dynjson

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"models"
	"plugin"
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
		"fmt"
	)

	%s
		
	func jsonDecode(input []byte) %s {
		var ret %s
		json.Unmarshal(input, &ret)
		return ret
	}`, definitions, definition, definition)
}

func ParseString(String string) {
	var i interface{}

	err := json.Unmarshal([]byte(String), &i)

	if err != nil {
		panic(err)
	}

	content := Parse(i)
	filename := fmt.Sprintf("./.temp/evaluated-%d.go", time.Now().UnixNano())

	ioutil.WriteFile(filename, []byte(content), 0644)

	p, err := plugin.Open(".temp/evaluated.go")

	if err != nil {
		panic(err)
	}

	fmt.Println(p)
}

func ParseFile(Name string) {
	Content, err := ioutil.ReadFile(Name)

	if err != nil {
		panic(err)
	}

	ParseString(string(Content))
}
